package climate

import (
	"math"
	"testing"
)

func TestComputeThermalStats(t *testing.T) {
	sst := []float64{18, 19, 20, 21, 22}
	ts := ComputeThermalStats(sst)
	if ts.MinSST != 18 || ts.MaxSST != 22 {
		t.Errorf("min=%f max=%f", ts.MinSST, ts.MaxSST)
	}
	if math.Abs(ts.MeanSST-20) > 0.01 {
		t.Errorf("mean expected 20, got %f", ts.MeanSST)
	}
	if ts.Range != 4 {
		t.Errorf("range expected 4, got %f", ts.Range)
	}
	if ts.Trend <= 0 {
		t.Error("increasing SST should have positive trend")
	}
}

func TestMixedLayerDepth(t *testing.T) {
	temps := []float64{25, 25, 24.5, 22, 18}
	depths := []float64{0, 10, 20, 30, 50}
	mld := MixedLayerDepth(temps, depths, 1.0)
	// Drop of 1°C happens between 20m (24.5) and 30m (22)
	if mld < 20 || mld > 30 {
		t.Errorf("MLD expected between 20-30, got %f", mld)
	}
}

func TestDegreeDay(t *testing.T) {
	temps := []float64{20, 22, 25, 18, 30}
	dd := DegreeDay(temps, 20)
	// (0 + 2 + 5 + 0 + 10) = 17
	if math.Abs(dd-17) > 0.01 {
		t.Errorf("degree-days expected 17, got %f", dd)
	}
}

func TestCoolingDegreeDay(t *testing.T) {
	temps := []float64{15, 18, 20, 12}
	cdd := CoolingDegreeDay(temps, 18)
	// (3 + 0 + 0 + 6) = 9
	if math.Abs(cdd-9) > 0.01 {
		t.Errorf("cooling degree-days expected 9, got %f", cdd)
	}
}

func TestAnomalySeries(t *testing.T) {
	data := []float64{20, 22, 18}
	anom := AnomalySeries(data, 20)
	if anom[0] != 0 || anom[1] != 2 || anom[2] != -2 {
		t.Errorf("anomaly: %v", anom)
	}
}

func TestComputeThermalStats_Empty(t *testing.T) {
	ts := ComputeThermalStats(nil)
	if ts.MeanSST != 0 {
		t.Error("empty should return 0")
	}
}
