package obs

import (
	"math"
	"testing"
)

func TestFilterByBuoy(t *testing.T) {
	readings := []Reading{{Buoy: "B1"}, {Buoy: "B2"}, {Buoy: "B1"}}
	out := FilterByBuoy(readings, "B1")
	if len(out) != 2 {
		t.Errorf("expected 2, got %d", len(out))
	}
}

func TestFilterByWindSpeed(t *testing.T) {
	readings := []Reading{{WindSpd: 5}, {WindSpd: 15}, {WindSpd: 25}}
	out := FilterByWindSpeed(readings, 10, 20)
	if len(out) != 1 || out[0].WindSpd != 15 {
		t.Errorf("expected 1 reading with windspd=15, got %v", out)
	}
}

func TestBuoyCodes(t *testing.T) {
	readings := []Reading{{Buoy: "B2"}, {Buoy: "B1"}, {Buoy: "B2"}, {Buoy: "B3"}}
	codes := BuoyCodes(readings)
	if len(codes) != 3 || codes[0] != "B2" {
		t.Errorf("codes: %v", codes)
	}
}

func TestFieldStats(t *testing.T) {
	readings := []Reading{
		{WaveHt: 1}, {WaveHt: 2}, {WaveHt: 3}, {WaveHt: 4},
	}
	s := FieldStats(readings, func(r Reading) float64 { return r.WaveHt })
	if s.Min != 1 || s.Max != 4 {
		t.Errorf("min=%f max=%f", s.Min, s.Max)
	}
	if math.Abs(s.Mean-2.5) > 0.001 {
		t.Errorf("mean expected 2.5, got %f", s.Mean)
	}
}

func TestCount(t *testing.T) {
	readings := []Reading{{WaveHt: 1}, {WaveHt: 5}, {WaveHt: 3}}
	n := Count(readings, func(r Reading) bool { return r.WaveHt > 2 })
	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}
