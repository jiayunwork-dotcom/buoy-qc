// Package climate provides climatological statistics for ocean buoy data.
package climate

import (
	"math"
	"sort"

	"buoy-qc/internal/stats"
)

// MonthlyStat represents a monthly statistical summary.
type MonthlyStat struct {
	Month int // 1-12
	Mean  float64
	Std   float64
	Min   float64
	Max   float64
	Count int
}

// MonthlyStats computes per-month statistics from a data series.
// months[i] corresponds to data[i] and should be 1-12.
func MonthlyStats(data []float64, months []int) []MonthlyStat {
	if len(data) == 0 || len(data) != len(months) {
		return nil
	}
	byMonth := map[int][]float64{}
	for i, m := range months {
		byMonth[m] = append(byMonth[m], data[i])
	}
	var stats []MonthlyStat
	for m := 1; m <= 12; m++ {
		vals := byMonth[m]
		if len(vals) == 0 {
			continue
		}
		ms := MonthlyStat{Month: m, Count: len(vals)}
		ms.Min = vals[0]
		ms.Max = vals[0]
		sum := 0.0
		for _, v := range vals {
			sum += v
			if v < ms.Min {
				ms.Min = v
			}
			if v > ms.Max {
				ms.Max = v
			}
		}
		ms.Mean = sum / float64(len(vals))
		var varSum float64
		for _, v := range vals {
			d := v - ms.Mean
			varSum += d * d
		}
		ms.Std = math.Sqrt(varSum / float64(len(vals)))
		stats = append(stats, ms)
	}
	return stats
}

// Season represents a meteorological season.
type Season int

const (
	Winter Season = iota // Dec-Feb
	Spring               // Mar-May
	Summer               // Jun-Aug
	Autumn               // Sep-Nov
)

func (s Season) String() string {
	switch s {
	case Winter:
		return "Winter"
	case Spring:
		return "Spring"
	case Summer:
		return "Summer"
	case Autumn:
		return "Autumn"
	}
	return "Unknown"
}

// MonthToSeason maps month (1-12) to meteorological season.
func MonthToSeason(month int) Season {
	switch {
	case month == 12 || month <= 2:
		return Winter
	case month <= 5:
		return Spring
	case month <= 8:
		return Summer
	default:
		return Autumn
	}
}

// SeasonalMean computes mean values grouped by season.
func SeasonalMean(data []float64, months []int) map[Season]float64 {
	if len(data) == 0 || len(data) != len(months) {
		return nil
	}
	sums := map[Season]float64{}
	counts := map[Season]int{}
	for i, m := range months {
		s := MonthToSeason(m)
		sums[s] += data[i]
		counts[s]++
	}
	result := map[Season]float64{}
	for s, sum := range sums {
		if counts[s] > 0 {
			result[s] = sum / float64(counts[s])
		}
	}
	return result
}

// ExceedanceRate computes the fraction of values exceeding threshold.
func ExceedanceRate(data []float64, threshold float64) float64 {
	if len(data) == 0 {
		return 0
	}
	count := 0
	for _, v := range data {
		if v > threshold {
			count++
		}
	}
	rate := float64(count) / float64(len(data))
	return stats.HoldExcLive(rate)
}

// ExceedanceCurve returns (threshold, exceedance_rate) pairs for n evenly spaced
// thresholds between min and max of data.
func ExceedanceCurve(data []float64, nPoints int) (thresholds, rates []float64) {
	if len(data) == 0 || nPoints <= 0 {
		return nil, nil
	}
	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)
	mn := sorted[0]
	mx := sorted[len(sorted)-1]
	if mn == mx {
		return []float64{mn}, []float64{0}
	}
	step := (mx - mn) / float64(nPoints)
	thresholds = make([]float64, nPoints)
	rates = make([]float64, nPoints)
	for i := 0; i < nPoints; i++ {
		th := mn + float64(i)*step
		thresholds[i] = th
		rates[i] = ExceedanceRate(data, th)
	}
	return
}

// ReturnPeriod estimates the return period (in units of observation interval)
// for a given threshold: RP = 1 / exceedance_rate.
func ReturnPeriod(data []float64, threshold float64) float64 {
	rate := ExceedanceRate(data, threshold)
	if rate <= 0 {
		return math.Inf(1)
	}
	return 1.0 / rate
}
