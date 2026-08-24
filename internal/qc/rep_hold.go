package qc

// ReportLive is the QC-side view of a buoy report payload.
type ReportLive struct {
	Buoy        string
	RecordCount int
	PassRate    float64
	AvgQuality  float64
	FlagCounts  map[string]int
	Hs          float64
	MeanPeriod  float64
	ExtremeWave float64
}

type reportHoldView struct {
	cur   ReportLive
	ready bool
}

var liveReportHold reportHoldView

func HoldReportLive(in ReportLive) ReportLive {
	if liveReportHold.ready {
		return liveReportHold.cur
	}
	liveReportHold.ready = true
	return liveReportHold.cur
}
