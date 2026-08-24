package qc

import "buoy-qc/internal/obs"

// Score represents an overall QC quality score for a reading.
type Score struct {
	Index      int
	RangeFlag  bool
	SpikeFlag  bool
	FlatFlag   bool
	GradFlag   bool
	TempFlag   bool
	Total      int // count of flags set
	Quality    float64 // 0.0 (worst) to 1.0 (best)
}

// ScoreReadings runs all QC checks on a set of readings and returns per-reading scores.
// Parameters: spikeK for spike detection, flatEps for flat check, gradMax for gradient,
// tempRate for temporal consistency.
func ScoreReadings(readings []obs.Reading, spikeK, flatEps, gradMax, tempRate float64) []Score {
	n := len(readings)
	if n == 0 {
		return nil
	}

	waveht := make([]float64, n)
	for i, r := range readings {
		waveht[i] = r.WaveHt
	}

	spikeFlags := SpikeDetect(waveht, spikeK)
	flatFlags := FlatCheck(waveht, flatEps)
	gradFlags := GradientCheck(waveht, gradMax)
	tempFlags := TemporalConsistency(waveht, tempRate)

	scores := make([]Score, n)
	for i, r := range readings {
		s := Score{Index: i}
		if RangeCheck(r) == 1 {
			s.RangeFlag = true
			s.Total++
		}
		if spikeFlags != nil && i < len(spikeFlags) && spikeFlags[i] {
			s.SpikeFlag = true
			s.Total++
		}
		if i < len(flatFlags) && flatFlags[i] {
			s.FlatFlag = true
			s.Total++
		}
		if i < len(gradFlags) && gradFlags[i] {
			s.GradFlag = true
			s.Total++
		}
		if i < len(tempFlags) && tempFlags[i] {
			s.TempFlag = true
			s.Total++
		}
		// Quality: 5 possible flags, each reduces by 0.2
		s.Quality = 1.0 - float64(s.Total)*0.2
		if s.Quality < 0 {
			s.Quality = 0
		}
		scores[i] = s
	}
	return scores
}

// AverageQuality returns the mean quality score across all readings.
func AverageQuality(scores []Score) float64 {
	if len(scores) == 0 {
		return 0
	}
	sum := 0.0
	for _, s := range scores {
		sum += s.Quality
	}
	return sum / float64(len(scores))
}

// FlagCounts returns the total count of each flag type across all scores.
func FlagCounts(scores []Score) map[string]int {
	counts := map[string]int{
		"range":    0,
		"spike":    0,
		"flat":     0,
		"gradient": 0,
		"temporal": 0,
	}
	for _, s := range scores {
		if s.RangeFlag {
			counts["range"]++
		}
		if s.SpikeFlag {
			counts["spike"]++
		}
		if s.FlatFlag {
			counts["flat"]++
		}
		if s.GradFlag {
			counts["gradient"]++
		}
		if s.TempFlag {
			counts["temporal"]++
		}
	}
	return counts
}

// PassRate returns the fraction of readings with Quality == 1.0.
func PassRate(scores []Score) float64 {
	if len(scores) == 0 {
		return 0
	}
	pass := 0
	for _, s := range scores {
		if s.Total == 0 {
			pass++
		}
	}
	return float64(pass) / float64(len(scores))
}
