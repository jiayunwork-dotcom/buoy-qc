package sea

import (
	"math"

	"buoy-qc/internal/obs"
)

// WindStats holds wind statistical summary.
type WindStats struct {
	MeanSpeed  float64
	MaxSpeed   float64
	MeanDir    float64
	Steadiness float64
}

// ComputeWindStats calculates wind speed/direction statistics.
func ComputeWindStats(readings []obs.Reading) WindStats {
	n := len(readings)
	if n == 0 {
		return WindStats{}
	}
	var sumSpd, maxSpd float64
	var sumU, sumV float64
	for _, r := range readings {
		sumSpd += r.WindSpd
		if r.WindSpd > maxSpd {
			maxSpd = r.WindSpd
		}
		rad := r.WindDir * math.Pi / 180
		sumU += r.WindSpd * math.Sin(rad)
		sumV += r.WindSpd * math.Cos(rad)
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
	return WindStats{
		MeanSpeed:  meanSpd,
		MaxSpeed:   maxSpd,
		MeanDir:    meanDir,
		Steadiness: steadiness,
	}
}

// BeaufortScale converts wind speed (m/s) to Beaufort number (0-12).
func BeaufortScale(windSpd float64) int {
	thresholds := []float64{0.3, 1.6, 3.4, 5.5, 8.0, 10.8, 13.9, 17.2, 20.8, 24.5, 28.5, 32.7}
	for i, th := range thresholds {
		if windSpd < th {
			return i
		}
	}
	return 12
}

// WindChill estimates wind chill temperature (°C) using the standard formula.
// Valid for air temperature <= 10°C and wind speed >= 1.3 m/s (converted to km/h).
func WindChill(airTemp, windSpd float64) float64 {
	vKmh := windSpd * 3.6
	if airTemp > 10 || vKmh < 4.8 {
		return airTemp
	}
	return 13.12 + 0.6215*airTemp - 11.37*math.Pow(vKmh, 0.16) + 0.3965*airTemp*math.Pow(vKmh, 0.16)
}

// GustFactor estimates gust factor from mean and max wind speed.
func GustFactor(meanSpd, maxSpd float64) float64 {
	if meanSpd <= 0 {
		return 0
	}
	return maxSpd / meanSpd
}

// WindPowerDensity estimates wind power density (W/m²) at sea level.
// P = 0.5 * rho * v^3, with rho ≈ 1.225 kg/m³.
func WindPowerDensity(windSpd float64) float64 {
	const rho = 1.225
	return 0.5 * rho * windSpd * windSpd * windSpd
}
