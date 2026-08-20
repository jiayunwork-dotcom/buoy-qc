package qc

import "math"

// TemporalConsistency checks whether consecutive readings are temporally
// consistent. A reading at index i is flagged if |x[i]-x[i-1]| > maxRate * dt,
// where dt is always 1 (unit time step). Returns a bool slice of same length.
func TemporalConsistency(series []float64, maxRate float64) []bool {
	n := len(series)
	out := make([]bool, n)
	for i := 1; i < n; i++ {
		if math.Abs(series[i]-series[i-1]) > maxRate {
			out[i] = true
		}
	}
	return out
}

// PersistenceCheck flags runs of identical values longer than maxRun.
// A value is "identical" if |x[i]-x[i-1]| < eps.
func PersistenceCheck(series []float64, eps float64, maxRun int) []bool {
	n := len(series)
	out := make([]bool, n)
	if n < 2 || maxRun < 1 {
		return out
	}
	runStart := 0
	for i := 1; i < n; i++ {
		if math.Abs(series[i]-series[i-1]) < eps {
			if i-runStart+1 > maxRun {
				for j := runStart; j <= i; j++ {
					out[j] = true
				}
			}
		} else {
			runStart = i
		}
	}
	return out
}

// AccelerationCheck flags index i where the second derivative approximation
// |x[i-1] - 2*x[i] + x[i+1]| exceeds threshold. Requires len >= 3.
func AccelerationCheck(series []float64, threshold float64) []bool {
	n := len(series)
	out := make([]bool, n)
	if n < 3 {
		return out
	}
	for i := 1; i < n-1; i++ {
		accel := math.Abs(series[i-1] - 2*series[i] + series[i+1])
		if accel > threshold {
			out[i] = true
		}
	}
	return out
}

// RepeatValueCheck flags indices where value equals the previous value exactly.
func RepeatValueCheck(series []float64) []bool {
	n := len(series)
	out := make([]bool, n)
	for i := 1; i < n; i++ {
		if series[i] == series[i-1] {
			out[i] = true
		}
	}
	return out
}
