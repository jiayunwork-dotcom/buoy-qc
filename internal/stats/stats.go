// Package stats provides statistical utility functions for ocean data analysis.
package stats

import (
	"math"
	"sort"
)

// Percentile computes the p-th percentile (0-100) of data using linear interpolation.
func Percentile(data []float64, p float64) float64 {
	n := len(data)
	if n == 0 || p < 0 || p > 100 {
		return 0
	}
	sorted := make([]float64, n)
	copy(sorted, data)
	sort.Float64s(sorted)
	rank := (p / 100) * float64(n-1)
	lo := int(math.Floor(rank))
	hi := int(math.Ceil(rank))
	if lo == hi || hi >= n {
		return sorted[lo]
	}
	frac := rank - float64(lo)
	return sorted[lo]*(1-frac) + sorted[hi]*frac
}

// Quantiles returns the 25th, 50th (median), and 75th percentiles.
func Quantiles(data []float64) (q25, q50, q75 float64) {
	return Percentile(data, 25), Percentile(data, 50), Percentile(data, 75)
}

// IQR returns the interquartile range (Q75 - Q25).
func IQR(data []float64) float64 {
	q25, _, q75 := Quantiles(data)
	return q75 - q25
}

// Outliers returns indices where values fall outside [Q25 - k*IQR, Q75 + k*IQR].
func Outliers(data []float64, k float64) []int {
	q25, _, q75 := Quantiles(data)
	iqr := q75 - q25
	lo := q25 - k*iqr
	hi := q75 + k*iqr
	var indices []int
	for i, v := range data {
		if v < lo || v > hi {
			indices = append(indices, i)
		}
	}
	return indices
}

// Histogram bins data into nBins equal-width bins between min and max.
// Returns bin edges (len=nBins+1) and counts (len=nBins).
func Histogram(data []float64, nBins int) (edges []float64, counts []int) {
	if len(data) == 0 || nBins <= 0 {
		return nil, nil
	}
	mn, mx := data[0], data[0]
	for _, v := range data {
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	if mn == mx {
		mx = mn + 1
	}
	width := (mx - mn) / float64(nBins)
	edges = make([]float64, nBins+1)
	counts = make([]int, nBins)
	for i := 0; i <= nBins; i++ {
		edges[i] = mn + float64(i)*width
	}
	for _, v := range data {
		idx := int((v - mn) / width)
		if idx >= nBins {
			idx = nBins - 1
		}
		if idx < 0 {
			idx = 0
		}
		counts[idx]++
	}
	return edges, counts
}

// MovingAverage smooths data with a sliding window of size w.
func MovingAverage(data []float64, w int) []float64 {
	if len(data) == 0 || w <= 0 {
		return nil
	}
	out := make([]float64, len(data))
	for i := range data {
		start := i - w/2
		if start < 0 {
			start = 0
		}
		end := i + w/2 + 1
		if end > len(data) {
			end = len(data)
		}
		sum := 0.0
		for j := start; j < end; j++ {
			sum += data[j]
		}
		out[i] = sum / float64(end-start)
	}
	return out
}

// Correlation computes the Pearson correlation coefficient between a and b.
func Correlation(a, b []float64) float64 {
	n := len(a)
	if n == 0 || n != len(b) {
		return 0
	}
	var sumA, sumB float64
	for i := 0; i < n; i++ {
		sumA += a[i]
		sumB += b[i]
	}
	meanA := sumA / float64(n)
	meanB := sumB / float64(n)
	var cov, varA, varB float64
	for i := 0; i < n; i++ {
		da := a[i] - meanA
		db := b[i] - meanB
		cov += da * db
		varA += da * da
		varB += db * db
	}
	denom := math.Sqrt(varA * varB)
	if denom == 0 {
		return 0
	}
	return cov / denom
}

// LinearRegression returns slope and intercept of least-squares fit y = slope*x + intercept.
func LinearRegression(x, y []float64) (slope, intercept float64) {
	n := len(x)
	if n < 2 || n != len(y) {
		return 0, 0
	}
	var sumX, sumY, sumXX, sumXY float64
	for i := 0; i < n; i++ {
		sumX += x[i]
		sumY += y[i]
		sumXX += x[i] * x[i]
		sumXY += x[i] * y[i]
	}
	nf := float64(n)
	denom := nf*sumXX - sumX*sumX
	if denom == 0 {
		return 0, sumY / nf
	}
	slope = (nf*sumXY - sumX*sumY) / denom
	intercept = (sumY - slope*sumX) / nf
	return
}
