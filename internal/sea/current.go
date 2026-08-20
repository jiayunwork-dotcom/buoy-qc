package sea

import (
	"math"

	"buoy-qc/internal/obs"
)

// CurrentStats holds statistical summary of current observations.
type CurrentStats struct {
	MeanSpeed  float64
	MaxSpeed   float64
	MeanDir    float64
	Steadiness float64 // ratio of vector mean to scalar mean speed
}

// ComputeCurrentStats calculates current statistics from readings.
func ComputeCurrentStats(readings []obs.Reading) CurrentStats {
	n := len(readings)
	if n == 0 {
		return CurrentStats{}
	}
	var sumSpd, maxSpd float64
	var sumU, sumV float64
	for _, r := range readings {
		sumSpd += r.CurrentSpd
		if r.CurrentSpd > maxSpd {
			maxSpd = r.CurrentSpd
		}
		rad := r.CurrentDir * math.Pi / 180
		sumU += r.CurrentSpd * math.Sin(rad)
		sumV += r.CurrentSpd * math.Cos(rad)
	}
	nf := float64(n)
	meanSpd := sumSpd / nf
	meanU := sumU / nf
	meanV := sumV / nf
	vectorMean := math.Sqrt(meanU*meanU + meanV*meanV)
	steadiness := 0.0
	if meanSpd > 0 {
		steadiness = vectorMean / meanSpd
	}
	meanDir := math.Atan2(meanU, meanV) * 180 / math.Pi
	if meanDir < 0 {
		meanDir += 360
	}
	return CurrentStats{
		MeanSpeed:  meanSpd,
		MaxSpeed:   maxSpd,
		MeanDir:    meanDir,
		Steadiness: steadiness,
	}
}

// CurrentRose returns fraction of readings in 8 directional sectors for current.
func CurrentRose(readings []obs.Reading) [8]float64 {
	var rose [8]float64
	if len(readings) == 0 {
		return rose
	}
	for _, r := range readings {
		d := math.Mod(r.CurrentDir, 360)
		if d < 0 {
			d += 360
		}
		idx := int(d / 45)
		if idx > 7 {
			idx = 7
		}
		rose[idx]++
	}
	total := 0.0
	for _, c := range rose {
		total += c
	}
	if total > 0 {
		for i := range rose {
			rose[i] /= total
		}
	}
	return rose
}

// TidalHarmonic represents a single tidal constituent.
type TidalHarmonic struct {
	Name      string
	Frequency float64 // cycles per hour
	Amplitude float64
	Phase     float64 // radians
}

// FitTidalHarmonics fits simple sinusoidal harmonics to a time series.
// Uses least-squares fit for each constituent independently.
// dt is the time step in hours.
func FitTidalHarmonics(series []float64, dt float64, constituents []float64) []TidalHarmonic {
	n := len(series)
	if n < 4 || dt <= 0 {
		return nil
	}
	var harmonics []TidalHarmonic
	for _, freq := range constituents {
		var sumCos, sumSin, sumCC, sumSS, sumCS float64
		m := mean(series)
		for i := 0; i < n; i++ {
			t := float64(i) * dt
			phase := 2 * math.Pi * freq * t
			c := math.Cos(phase)
			s := math.Sin(phase)
			val := series[i] - m
			sumCos += val * c
			sumSin += val * s
			sumCC += c * c
			sumSS += s * s
			sumCS += c * s
		}
		det := sumCC*sumSS - sumCS*sumCS
		if math.Abs(det) < 1e-10 {
			continue
		}
		a := (sumSS*sumCos - sumCS*sumSin) / det
		b := (sumCC*sumSin - sumCS*sumCos) / det
		amp := math.Sqrt(a*a + b*b)
		ph := math.Atan2(b, a)
		harmonics = append(harmonics, TidalHarmonic{
			Frequency: freq,
			Amplitude: amp,
			Phase:     ph,
		})
	}
	return harmonics
}
