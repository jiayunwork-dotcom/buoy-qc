// Package export generates QC analysis reports in text and JSON formats.
package export

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"buoy-qc/internal/qc"
)

// BuoyReport is the report for a single buoy.
type BuoyReport struct {
	Buoy         string         `json:"buoy"`
	RecordCount  int            `json:"record_count"`
	PassRate     float64        `json:"pass_rate"`
	AvgQuality   float64        `json:"avg_quality"`
	FlagCounts   map[string]int `json:"flag_counts"`
	Hs           float64        `json:"hs"`
	MeanPeriod   float64        `json:"mean_period"`
	ExtremeWave  float64        `json:"extreme_wave_50y"`
}

// FullReport is the complete analysis output.
type FullReport struct {
	Title  string       `json:"title"`
	Buoys  []BuoyReport `json:"buoys"`
}

// NewFullReport creates an empty report with title.
func NewFullReport(title string) *FullReport {
	return &FullReport{Title: title}
}

// AddBuoy appends a buoy report.
func (r *FullReport) AddBuoy(br BuoyReport) {
	r.Buoys = append(r.Buoys, br)
}

// WriteText writes a human-readable text report.
func (r *FullReport) WriteText(w io.Writer) {
	fmt.Fprintf(w, "=== %s ===\n\n", r.Title)
	for _, b := range r.Buoys {
		fmt.Fprintf(w, "Buoy: %s (n=%d)\n", b.Buoy, b.RecordCount)
		fmt.Fprintf(w, "  Pass Rate:      %.1f%%\n", b.PassRate*100)
		fmt.Fprintf(w, "  Avg Quality:    %.3f\n", b.AvgQuality)
		fmt.Fprintf(w, "  Hs:             %.3f m\n", b.Hs)
		fmt.Fprintf(w, "  Mean Period:    %.3f s\n", b.MeanPeriod)
		fmt.Fprintf(w, "  50y Extreme:    %.3f m\n", b.ExtremeWave)
		fmt.Fprintf(w, "  Flags: %s\n", formatFlags(b.FlagCounts))
		fmt.Fprintln(w)
	}
}

// WriteJSON writes a JSON report.
func (r *FullReport) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// String returns a text report as string.
func (r *FullReport) String() string {
	var b strings.Builder
	r.WriteText(&b)
	return b.String()
}

// MakeBuoyReport creates a BuoyReport from QC scores and sea-state metrics.
func MakeBuoyReport(buoy string, scores []qc.Score, hs, meanPer, extreme float64) BuoyReport {
	return BuoyReport{
		Buoy:        buoy,
		RecordCount: len(scores),
		PassRate:    qc.PassRate(scores),
		AvgQuality:  qc.AverageQuality(scores),
		FlagCounts:  qc.FlagCounts(scores),
		Hs:          hs,
		MeanPeriod:  meanPer,
		ExtremeWave: extreme,
	}
}

func formatFlags(counts map[string]int) string {
	parts := make([]string, 0, len(counts))
	for k, v := range counts {
		if v > 0 {
			parts = append(parts, fmt.Sprintf("%s=%d", k, v))
		}
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, " ")
}
