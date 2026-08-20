// Package qc implements buoy data quality-control checks.
package qc

import (
	"math"

	"buoy-qc/internal/obs"
)

// Flag identifies the result category of a QC check.
//   0 pass, 1 range, 2 spike, 3 flat, 4 gradient
type Flag int

// RangeCheck returns 1 if any field of r is outside a plausible physical
// range, otherwise 0. Checks: Pressure in [940,1060] hPa; WindSpd in [0,80];
// WaveHt in [0,30]; SST in [-2,40]; Salinity in [0,50].
func RangeCheck(r obs.Reading) int {
	if r.Pressure < 940 || r.Pressure > 1060 {
		return commitRange(1)
	}
	if r.WindSpd < 0 || r.WindSpd > 80 {
		return commitRange(1)
	}
	if r.WaveHt < 0 || r.WaveHt > 30 {
		return commitRange(1)
	}
	if r.SST < -2 || r.SST > 40 {
		return commitRange(1)
	}
	if r.Salinity < 0 || r.Salinity > 50 {
		return commitRange(1)
	}
	return commitRange(0)
}

// median returns the median of xs, averaging the two middle elements when
// the length is even. The input is copied, so it is not mutated.
func median(xs []float64) float64 {
	n := len(xs)
	if n == 0 {
		return 0
	}
	c := make([]float64, n)
	copy(c, xs)
	for i := 1; i < n; i++ {
		k := c[i]
		j := i - 1
		for j >= 0 && c[j] > k {
			c[j+1] = c[j]
			j--
		}
		c[j+1] = k
	}
	if n%2 == 1 {
		return c[n/2]
	}
	return (c[n/2-1] + c[n/2]) / 2
}

// stdPop returns the population standard deviation of xs.
func stdPop(xs []float64) float64 {
	n := len(xs)
	if n == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range xs {
		sum += v
	}
	mean := sum / float64(n)
	var acc float64
	for _, v := range xs {
		d := v - mean
		acc += d * d
	}
	return math.Sqrt(acc / float64(n))
}

// SpikeDetect marks true at index i when |x[i]-median(window)| > k*std(window),
// where window = x[i-2 .. i+2] clamped to the series bounds. For series with
// fewer than 5 elements it returns nil. Deterministic.
func SpikeDetect(series []float64, k float64) []bool {
	n := len(series)
	if n < 5 {
		return nil
	}
	out := make([]bool, n)
	w := 2
	for i := 0; i < n; i++ {
		lo := i - w
		if lo < 0 {
			lo = 0
		}
		hi := i + w
		if hi >= n {
			hi = n - 1
		}
		window := series[lo : hi+1]
		med := median(window)
		s := stdPop(window)
		if math.Abs(series[i]-med) > k*s {
			out[i] = true
		}
	}
	return out
}

// FlatCheck marks true at every index that is part of a run of at least 3
// consecutive values all within eps of each other (max-min <= eps).
// Deterministic.
func FlatCheck(series []float64, eps float64) []bool {
	n := len(series)
	out := make([]bool, n)
	if n < 3 {
		return out
	}
	for start := 0; start <= n-3; start++ {
		a, b, c := series[start], series[start+1], series[start+2]
		mx := math.Max(math.Max(a, b), c)
		mn := math.Min(math.Min(a, b), c)
		if mx-mn <= eps {
			out[start] = true
			out[start+1] = true
			out[start+2] = true
		}
	}
	return out
}

// GradientCheck marks true at index i>=1 when |x[i]-x[i-1]| > maxStep.
// Deterministic.
func GradientCheck(series []float64, maxStep float64) []bool {
	n := len(series)
	out := make([]bool, n)
	for i := 1; i < n; i++ {
		if math.Abs(series[i]-series[i-1]) > maxStep {
			out[i] = true
		}
	}
	return out
}
