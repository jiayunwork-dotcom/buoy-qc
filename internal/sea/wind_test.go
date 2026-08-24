package sea

import (
	"math"
	"testing"

	"buoy-qc/internal/obs"
)

func TestComputeWindStats(t *testing.T) {
	readings := []obs.Reading{
		{WindSpd: 5, WindDir: 0},
		{WindSpd: 10, WindDir: 0},
		{WindSpd: 15, WindDir: 0},
	}
	ws := ComputeWindStats(readings)
	if math.Abs(ws.MeanSpeed-10) > 0.01 {
		t.Errorf("mean speed expected 10, got %f", ws.MeanSpeed)
	}
	if ws.MaxSpeed != 15 {
		t.Errorf("max speed expected 15, got %f", ws.MaxSpeed)
	}
	if ws.Steadiness < 0.99 {
		t.Errorf("same direction steadiness should be ~1, got %f", ws.Steadiness)
	}
}

func TestBeaufortScale(t *testing.T) {
	tests := []struct {
		spd float64
		bf  int
	}{
		{0, 0}, {2, 2}, {6, 4}, {12, 6}, {35, 12},
	}
	for _, tt := range tests {
		got := BeaufortScale(tt.spd)
		if got != tt.bf {
			t.Errorf("Beaufort(%.0f) = %d, want %d", tt.spd, got, tt.bf)
		}
	}
}

func TestWindChill(t *testing.T) {
	// Standard: airTemp=0, wind=5 m/s ≈ 18 km/h
	wc := WindChill(0, 5)
	if wc > 0 {
		t.Errorf("wind chill at 0°C 5m/s should be negative, got %f", wc)
	}
}

func TestWindChill_WarmDay(t *testing.T) {
	// Above 10°C, should return airTemp
	wc := WindChill(20, 10)
	if wc != 20 {
		t.Errorf("warm day wind chill should be airTemp, got %f", wc)
	}
}

func TestGustFactor(t *testing.T) {
	gf := GustFactor(10, 15)
	if math.Abs(gf-1.5) > 0.001 {
		t.Errorf("gust factor expected 1.5, got %f", gf)
	}
}

func TestWindPowerDensity(t *testing.T) {
	p := WindPowerDensity(10)
	// 0.5 * 1.225 * 1000 = 612.5
	expected := 0.5 * 1.225 * 1000
	if math.Abs(p-expected) > 0.1 {
		t.Errorf("power density expected %f, got %f", expected, p)
	}
}
