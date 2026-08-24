package climate

import "math"

// ThermalStats holds ocean thermal climate statistics.
type ThermalStats struct {
	MeanSST    float64
	MinSST     float64
	MaxSST     float64
	Range      float64
	StdSST     float64
	Trend      float64 // °C per unit time
}

// ComputeThermalStats calculates SST climate statistics.
func ComputeThermalStats(sst []float64) ThermalStats {
	n := len(sst)
	if n == 0 {
		return ThermalStats{}
	}
	ts := ThermalStats{MinSST: sst[0], MaxSST: sst[0]}
	sum := 0.0
	for _, v := range sst {
		sum += v
		if v < ts.MinSST {
			ts.MinSST = v
		}
		if v > ts.MaxSST {
			ts.MaxSST = v
		}
	}
	ts.MeanSST = sum / float64(n)
	ts.Range = ts.MaxSST - ts.MinSST
	var varSum float64
	for _, v := range sst {
		d := v - ts.MeanSST
		varSum += d * d
	}
	ts.StdSST = math.Sqrt(varSum / float64(n))
	// Simple linear trend
	if n >= 2 {
		var sumX, sumXY, sumXX float64
		for i, v := range sst {
			x := float64(i)
			sumX += x
			sumXY += x * v
			sumXX += x * x
		}
		nf := float64(n)
		denom := nf*sumXX - sumX*sumX
		if denom != 0 {
			ts.Trend = (nf*sumXY - sumX*sum) / denom
		}
	}
	return ts
}

// MixedLayerDepth estimates mixed layer depth from temperature profile.
// Uses a threshold criterion: depth where T drops more than deltaT from surface.
func MixedLayerDepth(temps []float64, depths []float64, deltaT float64) float64 {
	if len(temps) < 2 || len(temps) != len(depths) || deltaT <= 0 {
		return 0
	}
	surfaceT := temps[0]
	for i := 1; i < len(temps); i++ {
		if surfaceT-temps[i] > deltaT {
			// Linear interpolation between i-1 and i
			if i > 0 {
				frac := (deltaT - (surfaceT - temps[i-1])) / (temps[i-1] - temps[i])
				return depths[i-1] + frac*(depths[i]-depths[i-1])
			}
			return depths[i]
		}
	}
	return depths[len(depths)-1]
}

// DegreeDay computes heating degree-days above a base temperature.
func DegreeDay(temps []float64, baseTemp float64) float64 {
	sum := 0.0
	for _, t := range temps {
		if t > baseTemp {
			sum += t - baseTemp
		}
	}
	return sum
}

// CoolingDegreeDay computes cooling degree-days below a base temperature.
func CoolingDegreeDay(temps []float64, baseTemp float64) float64 {
	sum := 0.0
	for _, t := range temps {
		if t < baseTemp {
			sum += baseTemp - t
		}
	}
	return sum
}

// ThermalAnomaly computes anomaly = value - climatological mean.
func ThermalAnomaly(value, climateMean float64) float64 {
	return value - climateMean
}

// AnomalySeries computes anomaly for each value against a single climate mean.
func AnomalySeries(data []float64, climateMean float64) []float64 {
	out := make([]float64, len(data))
	for i, v := range data {
		out[i] = v - climateMean
	}
	return out
}
