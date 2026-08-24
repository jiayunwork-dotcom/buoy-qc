package export

import (
	"fmt"
	"io"
	"strings"

	"buoy-qc/internal/qc"
)

// WriteCSV writes QC scores to a CSV writer.
func WriteCSV(w io.Writer, scores []qc.Score) {
	fmt.Fprintln(w, "index,range,spike,flat,gradient,temporal,total,quality")
	for _, s := range scores {
		fmt.Fprintf(w, "%d,%s,%s,%s,%s,%s,%d,%.3f\n",
			s.Index,
			boolStr(s.RangeFlag),
			boolStr(s.SpikeFlag),
			boolStr(s.FlatFlag),
			boolStr(s.GradFlag),
			boolStr(s.TempFlag),
			s.Total,
			s.Quality,
		)
	}
}

// WriteCSVString returns the CSV as a string.
func WriteCSVString(scores []qc.Score) string {
	var b strings.Builder
	WriteCSV(&b, scores)
	return b.String()
}

// WriteSummaryCSV writes per-buoy summary as CSV.
func WriteSummaryCSV(w io.Writer, reports []BuoyReport) {
	fmt.Fprintln(w, "buoy,n,pass_rate,avg_quality,hs,mean_period,extreme_50y")
	for _, r := range reports {
		fmt.Fprintf(w, "%s,%d,%.4f,%.4f,%.3f,%.3f,%.3f\n",
			r.Buoy, r.RecordCount, r.PassRate, r.AvgQuality, r.Hs, r.MeanPeriod, r.ExtremeWave)
	}
}

// WriteSummaryCSVString returns per-buoy summary CSV as string.
func WriteSummaryCSVString(reports []BuoyReport) string {
	var b strings.Builder
	WriteSummaryCSV(&b, reports)
	return b.String()
}

func boolStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
