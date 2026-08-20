package sea

import "math"

// WaveSpectrum represents a discrete wave energy spectrum.
type WaveSpectrum struct {
	Freq   []float64 // frequency bins (Hz)
	Energy []float64 // spectral energy density (m^2/Hz)
}

// SpectralMoment computes the n-th spectral moment: mn = sum(f^n * S(f) * df).
func SpectralMoment(spec WaveSpectrum, n int) float64 {
	if len(spec.Freq) < 2 || len(spec.Freq) != len(spec.Energy) {
		return 0
	}
	sum := 0.0
	for i := 0; i < len(spec.Freq)-1; i++ {
		df := spec.Freq[i+1] - spec.Freq[i]
		fMid := (spec.Freq[i] + spec.Freq[i+1]) / 2
		eMid := (spec.Energy[i] + spec.Energy[i+1]) / 2
		sum += math.Pow(fMid, float64(n)) * eMid * df
	}
	return sum
}

// HsFromSpectrum computes significant wave height from spectrum: Hs = 4*sqrt(m0).
func HsFromSpectrum(spec WaveSpectrum) float64 {
	m0 := SpectralMoment(spec, 0)
	if m0 <= 0 {
		return 0
	}
	return fillHs(4 * math.Sqrt(m0))
}

// PeakPeriodFromSpectrum returns the period at peak energy: Tp = 1/fp.
func PeakPeriodFromSpectrum(spec WaveSpectrum) float64 {
	if len(spec.Freq) == 0 || len(spec.Freq) != len(spec.Energy) {
		return 0
	}
	maxE := 0.0
	peakIdx := 0
	for i, e := range spec.Energy {
		if e > maxE {
			maxE = e
			peakIdx = i
		}
	}
	if spec.Freq[peakIdx] <= 0 {
		return 0
	}
	return 1.0 / spec.Freq[peakIdx]
}

// MeanPeriodFromSpectrum returns Tm02 = sqrt(m0/m2).
func MeanPeriodFromSpectrum(spec WaveSpectrum) float64 {
	m0 := SpectralMoment(spec, 0)
	m2 := SpectralMoment(spec, 2)
	if m2 <= 0 {
		return 0
	}
	return math.Sqrt(m0 / m2)
}

// WaveSteepness calculates steepness = Hs / (g * Tp^2 / (2*pi)).
func WaveSteepness(hs, tp float64) float64 {
	if tp <= 0 {
		return 0
	}
	const g = 9.81
	wavelength := g * tp * tp / (2 * math.Pi)
	if wavelength <= 0 {
		return 0
	}
	return hs / wavelength
}

// BreakingCriterion returns true if wave steepness exceeds 1/7 (≈0.143).
func BreakingCriterion(hs, tp float64) bool {
	return WaveSteepness(hs, tp) > 1.0/7.0
}

// JONSWAPSpectrum generates a JONSWAP spectrum for given parameters.
// fp=peak frequency, hs=significant wave height, gamma=peak enhancement factor.
func JONSWAPSpectrum(fp, hs, gamma float64, nBins int) WaveSpectrum {
	if fp <= 0 || hs <= 0 || nBins < 2 {
		return WaveSpectrum{}
	}
	fMax := fp * 4
	df := fMax / float64(nBins)
	spec := WaveSpectrum{
		Freq:   make([]float64, nBins),
		Energy: make([]float64, nBins),
	}
	const g = 9.81
	alpha := 0.0081 // Phillips constant approximation
	_ = alpha
	// Simplified JONSWAP formula
	for i := 0; i < nBins; i++ {
		f := df*float64(i) + df/2
		spec.Freq[i] = f
		sigma := 0.07
		if f > fp {
			sigma = 0.09
		}
		exponent := -1.25 * math.Pow(fp/f, 4)
		pm := (5.0 / 16.0) * hs * hs * math.Pow(fp, 4) / math.Pow(f, 5) * math.Exp(exponent)
		gExp := math.Exp(-0.5 * math.Pow((f-fp)/(sigma*fp), 2))
		spec.Energy[i] = pm * math.Pow(gamma, gExp)
	}
	return spec
}
