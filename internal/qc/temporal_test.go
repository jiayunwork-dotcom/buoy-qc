package qc

import "testing"

func TestTemporalConsistency(t *testing.T) {
	series := []float64{10, 11, 15, 16, 10}
	flags := TemporalConsistency(series, 3.0)
	if !flags[2] { // jump from 11 to 15 = 4 > 3
		t.Error("expected flag at index 2")
	}
	if flags[1] { // 11-10=1 <= 3
		t.Error("index 1 should not be flagged")
	}
	if !flags[4] { // 16 to 10 = 6 > 3
		t.Error("expected flag at index 4")
	}
}

func TestPersistenceCheck(t *testing.T) {
	series := []float64{1, 1, 1, 1, 2, 3, 3, 3, 3, 3}
	flags := PersistenceCheck(series, 0.01, 3)
	// first 4 values form a run of 4 > maxRun=3
	if !flags[0] || !flags[3] {
		t.Error("first run should be flagged")
	}
	// last 5 values form a run of 5 > 3
	if !flags[5] || !flags[9] {
		t.Error("second run should be flagged")
	}
}

func TestAccelerationCheck(t *testing.T) {
	// smooth: 1,2,3,4,5 → accel = 0 everywhere
	smooth := []float64{1, 2, 3, 4, 5}
	flags := AccelerationCheck(smooth, 0.5)
	for i, f := range flags {
		if f {
			t.Errorf("smooth series should not flag index %d", i)
		}
	}
	// sharp: sudden change at index 2
	sharp := []float64{1, 2, 10, 3, 4}
	flags2 := AccelerationCheck(sharp, 5.0)
	if !flags2[2] { // |2 - 20 + 3| = 15 > 5
		t.Error("expected flag at index 2 for sharp change")
	}
}

func TestRepeatValueCheck(t *testing.T) {
	series := []float64{1, 2, 2, 3, 3, 3}
	flags := RepeatValueCheck(series)
	if flags[0] {
		t.Error("first element should not be flagged")
	}
	if !flags[2] || !flags[4] || !flags[5] {
		t.Error("repeated values should be flagged")
	}
}
