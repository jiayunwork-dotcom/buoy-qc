package sea

import (
	"math"
	"testing"
)

func TestHsFromSpectrum(t *testing.T) {
	spec := JONSWAPSpectrum(0.1, 2.0, 3.3, 50)
	hs := HsFromSpectrum(spec)
	// Should be approximately 2.0m (input Hs)
	if hs < 1.0 || hs > 4.0 {
		t.Errorf("Hs from JONSWAP expected ~2.0, got %f", hs)
	}
}

func TestPeakPeriodFromSpectrum(t *testing.T) {
	spec := JONSWAPSpectrum(0.1, 2.0, 3.3, 100)
	tp := PeakPeriodFromSpectrum(spec)
	// peak frequency 0.1 Hz → period 10s
	if math.Abs(tp-10.0) > 2.0 {
		t.Errorf("Tp expected ~10s, got %f", tp)
	}
}

func TestMeanPeriodFromSpectrum(t *testing.T) {
	spec := JONSWAPSpectrum(0.1, 2.0, 3.3, 100)
	tm := MeanPeriodFromSpectrum(spec)
	if tm <= 0 {
		t.Errorf("Tm02 should be > 0, got %f", tm)
	}
}

func TestWaveSteepness(t *testing.T) {
	st := WaveSteepness(2.0, 10.0)
	if st <= 0 || st > 0.2 {
		t.Errorf("steepness out of range: %f", st)
	}
}

func TestBreakingCriterion(t *testing.T) {
	// Very steep wave
	if !BreakingCriterion(5.0, 3.0) {
		t.Error("very steep wave should break")
	}
	// Gentle wave
	if BreakingCriterion(0.5, 10.0) {
		t.Error("gentle wave should not break")
	}
}

func TestSpectralMoment_Empty(t *testing.T) {
	spec := WaveSpectrum{}
	if SpectralMoment(spec, 0) != 0 {
		t.Error("empty spectrum moment should be 0")
	}
}
