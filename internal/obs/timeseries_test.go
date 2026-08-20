package obs

import (
	"math"
	"testing"
)

func TestResample_Double(t *testing.T) {
	data := []float64{0, 2, 4}
	out := Resample(data, 5)
	if len(out) != 5 {
		t.Fatalf("expected 5 samples, got %d", len(out))
	}
	if math.Abs(out[0]) > 0.001 || math.Abs(out[4]-4) > 0.001 {
		t.Errorf("endpoints: [%f, %f]", out[0], out[4])
	}
	if math.Abs(out[2]-2) > 0.001 {
		t.Errorf("midpoint expected 2, got %f", out[2])
	}
}

func TestInterpolate(t *testing.T) {
	data := []float64{1, -999, -999, 4, 5}
	out := Interpolate(data, -999)
	if math.Abs(out[1]-2) > 0.001 {
		t.Errorf("out[1] expected 2, got %f", out[1])
	}
	if math.Abs(out[2]-3) > 0.001 {
		t.Errorf("out[2] expected 3, got %f", out[2])
	}
}

func TestDownsample(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6}
	out := Downsample(data, 2)
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
	if math.Abs(out[0]-1.5) > 0.001 {
		t.Errorf("out[0] expected 1.5, got %f", out[0])
	}
}

func TestDiff(t *testing.T) {
	data := []float64{1, 3, 6, 10}
	d := Diff(data)
	expected := []float64{2, 3, 4}
	for i, v := range expected {
		if math.Abs(d[i]-v) > 0.001 {
			t.Errorf("diff[%d] expected %f, got %f", i, v, d[i])
		}
	}
}

func TestCumulativeSum(t *testing.T) {
	data := []float64{1, 2, 3, 4}
	cs := CumulativeSum(data)
	if cs[3] != 10 {
		t.Errorf("cumsum[3] expected 10, got %f", cs[3])
	}
}

func TestExtractField(t *testing.T) {
	readings := []Reading{
		{WaveHt: 1.5}, {WaveHt: 2.0}, {WaveHt: 2.5},
	}
	wh := ExtractField(readings, func(r Reading) float64 { return r.WaveHt })
	if len(wh) != 3 || wh[1] != 2.0 {
		t.Errorf("extract failed: %v", wh)
	}
}
