package qc

import (
	"testing"

	"buoy-qc/internal/obs"
)

func goodReading() obs.Reading {
	return obs.Reading{
		Buoy: "B1", WindSpd: 10, WindDir: 180, AirTemp: 20,
		Pressure: 1013, WaveHt: 2.0, WavePer: 8, SST: 18, Salinity: 35,
		CurrentSpd: 0.5, CurrentDir: 90,
	}
}

func TestScoreReadings_AllPass(t *testing.T) {
	readings := make([]obs.Reading, 6)
	for i := range readings {
		readings[i] = goodReading()
		readings[i].WaveHt = 2.0 + float64(i)*0.3 // vary to avoid flat/spike
	}
	scores := ScoreReadings(readings, 3.0, 0.05, 5.0, 5.0)
	for i, s := range scores {
		if s.Total != 0 {
			t.Errorf("reading %d should pass all checks, total=%d", i, s.Total)
		}
	}
}

func TestScoreReadings_RangeFlag(t *testing.T) {
	readings := []obs.Reading{goodReading(), goodReading(), goodReading(),
		goodReading(), goodReading(), goodReading()}
	readings[2].Pressure = 800 // out of range
	scores := ScoreReadings(readings, 2.0, 0.05, 5.0, 5.0)
	if !scores[2].RangeFlag {
		t.Error("expected range flag at index 2")
	}
}

func TestAverageQuality(t *testing.T) {
	scores := []Score{
		{Quality: 1.0},
		{Quality: 0.8},
		{Quality: 0.6},
	}
	avg := AverageQuality(scores)
	expected := 0.8
	if avg < expected-0.01 || avg > expected+0.01 {
		t.Errorf("average quality expected ~0.8, got %f", avg)
	}
}

func TestFlagCounts(t *testing.T) {
	scores := []Score{
		{RangeFlag: true, Total: 1},
		{SpikeFlag: true, GradFlag: true, Total: 2},
	}
	counts := FlagCounts(scores)
	if counts["range"] != 1 {
		t.Errorf("range count: %d", counts["range"])
	}
	if counts["spike"] != 1 {
		t.Errorf("spike count: %d", counts["spike"])
	}
}

func TestPassRate(t *testing.T) {
	scores := []Score{
		{Total: 0}, {Total: 0}, {Total: 1}, {Total: 0},
	}
	rate := PassRate(scores)
	if rate < 0.74 || rate > 0.76 {
		t.Errorf("pass rate expected 0.75, got %f", rate)
	}
}
