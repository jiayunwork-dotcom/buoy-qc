package sea

import (
	"math"
	"testing"

	"buoy-qc/internal/obs"
)

func TestComputeCurrentStats(t *testing.T) {
	readings := []obs.Reading{
		{CurrentSpd: 1.0, CurrentDir: 0},
		{CurrentSpd: 1.0, CurrentDir: 0},
		{CurrentSpd: 2.0, CurrentDir: 0},
	}
	stats := ComputeCurrentStats(readings)
	expectedMean := 4.0 / 3.0
	if math.Abs(stats.MeanSpeed-expectedMean) > 0.01 {
		t.Errorf("mean speed expected %f, got %f", expectedMean, stats.MeanSpeed)
	}
	if stats.MaxSpeed != 2.0 {
		t.Errorf("max speed expected 2.0, got %f", stats.MaxSpeed)
	}
	// All same direction → steadiness should be ~1
	if stats.Steadiness < 0.99 {
		t.Errorf("steadiness expected ~1.0, got %f", stats.Steadiness)
	}
}

func TestComputeCurrentStats_Empty(t *testing.T) {
	stats := ComputeCurrentStats(nil)
	if stats.MeanSpeed != 0 {
		t.Error("empty should return 0")
	}
}

func TestCurrentRose(t *testing.T) {
	readings := []obs.Reading{
		{CurrentDir: 10},  // sector 0
		{CurrentDir: 50},  // sector 1
		{CurrentDir: 10},  // sector 0
		{CurrentDir: 200}, // sector 4
	}
	rose := CurrentRose(readings)
	if math.Abs(rose[0]-0.5) > 0.01 {
		t.Errorf("sector 0 expected 0.5, got %f", rose[0])
	}
}

func TestFitTidalHarmonics(t *testing.T) {
	// Generate a simple sinusoidal signal: 1.5*cos(2*pi*0.08*t)
	n := 100
	dt := 1.0 // 1 hour
	series := make([]float64, n)
	freq := 0.08
	for i := range series {
		series[i] = 1.5 * math.Cos(2*math.Pi*freq*float64(i)*dt)
	}
	harmonics := FitTidalHarmonics(series, dt, []float64{freq})
	if len(harmonics) != 1 {
		t.Fatalf("expected 1 harmonic, got %d", len(harmonics))
	}
	if math.Abs(harmonics[0].Amplitude-1.5) > 0.2 {
		t.Errorf("amplitude expected ~1.5, got %f", harmonics[0].Amplitude)
	}
}
