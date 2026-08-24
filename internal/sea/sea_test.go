package sea

import (
	"math"
	"testing"

	"buoy-qc/internal/obs"
)

// TestSignificantWaveHt verifies Hs = 4*sqrt(variance) and the empty case.
func TestSignificantWaveHt(t *testing.T) {
	if SignificantWaveHt(nil) != 0 {
		t.Fatal("expected 0 for empty input")
	}
	// variance of {-1, 1} is 1 -> Hs should be 4.0
	elev := []float64{-1, 1}
	if got := SignificantWaveHt(elev); math.Abs(got-4.0) > 1e-9 {
		t.Fatalf("expected 4.0, got %v", got)
	}
}

// TestMeanWavePeriod verifies the mean and the empty case.
func TestMeanWavePeriod(t *testing.T) {
	if MeanWavePeriod(nil) != 0 {
		t.Fatal("expected 0 for empty input")
	}
	if got := MeanWavePeriod([]float64{4, 6}); math.Abs(got-5.0) > 1e-9 {
		t.Fatalf("expected 5.0, got %v", got)
	}
}

// TestWindRose verifies sector fractions sum to 1.0, correct mapping, and the
// nil-safe all-zero result.
func TestWindRose(t *testing.T) {
	r := WindRose(nil)
	for i, v := range r {
		if v != 0 {
			t.Fatalf("expected 0 at sector %d for empty, got %v", i, v)
		}
	}
	readings := []obs.Reading{
		{WindDir: 10}, {WindDir: 20}, {WindDir: 50}, {WindDir: 200},
	}
	r = WindRose(readings)
	sum := 0.0
	for _, v := range r {
		sum += v
	}
	if math.Abs(sum-1.0) > 1e-9 {
		t.Fatalf("expected sectors to sum to 1.0, got %v", sum)
	}
	if math.Abs(r[0]-0.5) > 1e-9 {
		t.Fatalf("expected sector0 = 0.5, got %v", r[0])
	}
	if math.Abs(r[1]-0.25) > 1e-9 {
		t.Fatalf("expected sector1 = 0.25, got %v", r[1])
	}
	if math.Abs(r[4]-0.25) > 1e-9 {
		t.Fatalf("expected sector4 = 0.25, got %v", r[4])
	}
}

// TestExtremeReturn verifies the empty/rp<=1 guards and that a larger return
// period yields a larger level.
func TestExtremeReturn(t *testing.T) {
	if ExtremeReturn(nil, 10) != 0 {
		t.Fatal("expected 0 for empty input")
	}
	if ExtremeReturn([]float64{1, 2, 3}, 1) != 0 {
		t.Fatal("expected 0 for rp<=1")
	}
	s := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	a := ExtremeReturn(s, 10)
	b := ExtremeReturn(s, 100)
	if !(b > a) {
		t.Fatalf("expected larger level for larger rp (a=%v b=%v)", a, b)
	}
}
