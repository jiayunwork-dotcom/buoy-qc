// Package cli implements the buoy-qc command-line interface with subcommands.
package cli

import (
	"flag"
	"fmt"
	"io"
	"sort"

	"buoy-qc/internal/export"
	"buoy-qc/internal/obs"
	"buoy-qc/internal/qc"
	"buoy-qc/internal/sea"
)

// Run is the CLI entry point.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}
	cmd := args[0]
	sub := args[1:]
	switch cmd {
	case "analyze":
		return cmdAnalyze(sub, stdout, stderr)
	case "report":
		return cmdReport(sub, stdout, stderr)
	case "help", "--help", "-h":
		printUsage(stdout)
		return 0
	default:
		// Legacy mode: treat as flags (backward compat)
		return cmdAnalyze(args, stdout, stderr)
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: buoy-qc <command> [options]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "commands:")
	fmt.Fprintln(w, "  analyze  -readings <path> [-buoy <code>]   QC + sea-state analysis")
	fmt.Fprintln(w, "  report   -readings <path> [-json]          Generate report")
	fmt.Fprintln(w, "  help                                       Show this message")
}

func cmdAnalyze(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	fs.SetOutput(stderr)
	readingsPath := fs.String("readings", "", "path to readings CSV")
	buoyFilter := fs.String("buoy", "", "optional buoy code filter")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *readingsPath == "" {
		fmt.Fprintln(stderr, "error: -readings is required")
		return 2
	}
	rs, err := obs.ParseReadings(*readingsPath)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	if *buoyFilter != "" {
		filtered := make([]obs.Reading, 0)
		for _, r := range rs {
			if r.Buoy == *buoyFilter {
				filtered = append(filtered, r)
			}
		}
		rs = filtered
	}
	if len(rs) == 0 {
		fmt.Fprintln(stderr, "error: no readings")
		return 1
	}
	byBuoy := obs.ByBuoy(rs)
	codes := sortedKeys(byBuoy)
	for _, code := range codes {
		list := byBuoy[code]
		waveht := obs.ExtractField(list, func(r obs.Reading) float64 { return r.WaveHt })
		waveper := obs.ExtractField(list, func(r obs.Reading) float64 { return r.WavePer })
		scores := qc.ScoreReadings(list, 2.0, 0.05, 2.0, 3.0)
		hs := sea.SignificantWaveHt(waveht)
		mwp := sea.MeanWavePeriod(waveper)
		ext := sea.ExtremeReturn(waveht, 50)
		counts := qc.FlagCounts(scores)
		fmt.Fprintf(stdout, "Buoy %s: n=%d pass=%.1f%% Hs=%.3f meanT=%.3f ext50=%.3f flags(range=%d spike=%d flat=%d grad=%d temp=%d)\n",
			code, len(list), qc.PassRate(scores)*100, hs, mwp, ext,
			counts["range"], counts["spike"], counts["flat"], counts["gradient"], counts["temporal"])
	}
	rose := sea.WindRose(rs)
	fmt.Fprintln(stdout, "WindRose:")
	for i, f := range rose {
		fmt.Fprintf(stdout, "  sector %d [%d,%d): %.4f\n", i, i*45, (i+1)*45, f)
	}
	return 0
}

func cmdReport(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	fs.SetOutput(stderr)
	readingsPath := fs.String("readings", "", "path to readings CSV")
	jsonFlag := fs.Bool("json", false, "JSON output")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *readingsPath == "" {
		fmt.Fprintln(stderr, "error: -readings is required")
		return 2
	}
	rs, err := obs.ParseReadings(*readingsPath)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	rpt := export.NewFullReport("Buoy QC Report")
	byBuoy := obs.ByBuoy(rs)
	for _, code := range sortedKeys(byBuoy) {
		list := byBuoy[code]
		waveht := obs.ExtractField(list, func(r obs.Reading) float64 { return r.WaveHt })
		waveper := obs.ExtractField(list, func(r obs.Reading) float64 { return r.WavePer })
		scores := qc.ScoreReadings(list, 2.0, 0.05, 2.0, 3.0)
		hs := sea.SignificantWaveHt(waveht)
		mwp := sea.MeanWavePeriod(waveper)
		ext := sea.ExtremeReturn(waveht, 50)
		br := export.MakeBuoyReport(code, scores, hs, mwp, ext)
		rpt.AddBuoy(br)
	}
	if *jsonFlag {
		if err := rpt.WriteJSON(stdout); err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
	} else {
		rpt.WriteText(stdout)
	}
	return 0
}

func sortedKeys(m map[string][]obs.Reading) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}


