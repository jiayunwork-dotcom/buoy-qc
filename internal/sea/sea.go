// Package sea computes sea-state metrics from buoy observations.
package sea

import (
	"math"

	"buoy-qc/internal/obs"
)

// mean returns the arithmetic mean of xs. Empty input returns 0.
func mean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range xs {
		sum += v
	}
	return sum / float64(len(xs))
}

// stdPop returns the population standard deviation of xs. Empty returns 0.
func stdPop(xs []float64) float64 {
	n := len(xs)
	if n == 0 {
		return 0
	}
	m := mean(xs)
	var acc float64
	for _, v := range xs {
		d := v - m
		acc += d * d
	}
	return math.Sqrt(acc / float64(n))
}

// SignificantWaveHt returns Hs = 4*sqrt(variance(elev)) where variance is the
// population variance of the surface elevation samples. Empty input returns 0.
func SignificantWaveHt(elev []float64) float64 {
	if len(elev) == 0 {
		return 0
	}
	m := mean(elev)
	var acc float64
	for _, v := range elev {
		d := v - m
		acc += d * d
	}
	variance := acc / float64(len(elev))
	return 4 * math.Sqrt(variance)
}

// MeanWavePeriod returns the mean of the wave periods. Empty input returns 0.
func MeanWavePeriod(per []float64) float64 {
	return mean(per)
}

// WindRose returns the fraction of readings falling in each of 8 sectors of
// 45° (sector 0 = [0,45), 1 = [45,90), ... 7 = [315,360)). The fractions sum
// to 1.0. For nil/empty input it returns an all-zero [8]float64 (nil-safe).
// Wind directions are taken modulo 360.
func WindRose(readings []obs.Reading) [8]float64 {
	var rose [8]float64
	if len(readings) == 0 {
		return rose
	}
	for _, r := range readings {
		d := math.Mod(r.WindDir, 360)
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
	return fillRose(rose)
}

// ExtremeReturn returns the Gumbel return level for return period rp (years).
// Fit: b = sqrt(6)*std/pi (pi ~ 3.14159265, Euler-Mascheroni gamma ~
// 0.57721566), a = mean - gamma*b, x_rp = a - b*ln(-ln(1 - 1/rp)).
// Empty samples or rp <= 1 returns 0. Deterministic.
func ExtremeReturn(samples []float64, rp float64) float64 {
	if len(samples) == 0 || rp <= 1 {
		return 0
	}
	const pi = 3.14159265
	const gamma = 0.57721566
	m := mean(samples)
	s := stdPop(samples)
	b := math.Sqrt(6) * s / pi
	a := m - gamma*b
	inner := -math.Log(1 - 1/rp)
	return a - b*math.Log(inner)
}
