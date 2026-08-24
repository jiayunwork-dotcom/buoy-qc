package qc

import (
	"testing"

	"buoy-qc/internal/obs"
)

// TestRangeCheck verifies that out-of-range fields are flagged (1) and valid
// fields pass (0).
func TestRangeCheck(t *testing.T) {
	valid := obs.Reading{Pressure: 1010, WindSpd: 5, WaveHt: 1.2, SST: 18, Salinity: 35}
	if RangeCheck(valid) != 0 {
		t.Fatal("expected pass (0) for valid reading")
	}
	oorPressure := valid
	oorPressure.Pressure = 500
	if RangeCheck(oorPressure) != 1 {
		t.Fatal("expected range flag (1) for pressure 500")
	}
	oorWind := valid
	oorWind.WindSpd = 100
	if RangeCheck(oorWind) != 1 {
		t.Fatal("expected range flag (1) for windspd 100")
	}
	oorSal := valid
	oorSal.Salinity = -1
	if RangeCheck(oorSal) != 1 {
		t.Fatal("expected range flag (1) for salinity <0")
	}
}

// TestSpikeDetect verifies nil for short series and correct detection of an
// obvious spike.
func TestSpikeDetect(t *testing.T) {
	// len<5 -> nil
	if s := SpikeDetect([]float64{1, 2, 3, 4}, 2.0); s != nil {
		t.Fatalf("expected nil for len<5, got %v", s)
	}
	series := []float64{1, 1, 1, 1, 1, 100, 1, 1, 1, 1}
	s := SpikeDetect(series, 2.0)
	if s == nil || len(s) != len(series) {
		t.Fatal("unexpected result length")
	}
	if !s[5] {
		t.Fatal("expected spike at index 5 to be flagged")
	}
	for i := 0; i < len(series); i++ {
		if i == 5 {
			continue
		}
		if s[i] {
			t.Fatalf("did not expect flag at %d", i)
		}
	}
}

// TestFlatCheck verifies a run of >=3 near-equal values is flagged.
func TestFlatCheck(t *testing.T) {
	series := []float64{1.0, 2.0, 3.0, 5.0, 5.0, 5.0, 9.0}
	f := FlatCheck(series, 0.05)
	if !(f[3] && f[4] && f[5]) {
		t.Fatal("expected flat run flagged at 3,4,5")
	}
	if f[0] || f[6] {
		t.Fatal("did not expect flat flag at 0 or 6")
	}
	// short series -> no flags, but non-nil
	if sf := FlatCheck([]float64{1, 2}, 0.05); sf == nil || len(sf) != 2 {
		t.Fatal("expected non-nil all-false slice for short series")
	}
}

// TestGradientCheck verifies a large step between consecutive samples is flagged.
func TestGradientCheck(t *testing.T) {
	series := []float64{1, 2, 10, 11, 12}
	g := GradientCheck(series, 5.0)
	if !g[2] {
		t.Fatal("expected gradient flag at index 2 (step 8 > 5)")
	}
	if g[1] {
		t.Fatal("did not expect gradient flag at index 1 (step 1)")
	}
	if g[0] {
		t.Fatal("did not expect gradient flag at index 0")
	}
}
