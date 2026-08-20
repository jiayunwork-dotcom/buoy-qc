package qc

import (
	"math"
	"testing"

	"buoy-qc/internal/obs"
)

func TestSpatialConsistency_Normal(t *testing.T) {
	target := obs.Reading{SST: 20}
	neighbors := []obs.Reading{{SST: 19.5}, {SST: 20.2}, {SST: 19.8}}
	flagged := SpatialConsistency(target, neighbors, func(r obs.Reading) float64 { return r.SST }, 2.0)
	if flagged {
		t.Error("normal value should not be flagged")
	}
}

func TestSpatialConsistency_Outlier(t *testing.T) {
	target := obs.Reading{SST: 30}
	neighbors := []obs.Reading{{SST: 20}, {SST: 20.5}, {SST: 19.5}}
	flagged := SpatialConsistency(target, neighbors, func(r obs.Reading) float64 { return r.SST }, 2.0)
	if !flagged {
		t.Error("outlier should be flagged")
	}
}

func TestBuoyDistance_SamePoint(t *testing.T) {
	d := BuoyDistance(30, 120, 30, 120)
	if d != 0 {
		t.Errorf("same point distance should be 0, got %f", d)
	}
}

func TestBuoyDistance_Known(t *testing.T) {
	// Shanghai (31.2, 121.5) to Tokyo (35.7, 139.7) ≈ 1760 km
	d := BuoyDistance(31.2, 121.5, 35.7, 139.7)
	if math.Abs(d-1760) > 100 {
		t.Errorf("distance Shanghai-Tokyo expected ~1760km, got %f", d)
	}
}

func TestNearbyBuoys(t *testing.T) {
	positions := map[string][2]float64{
		"B1": {30, 120},
		"B2": {30.01, 120.01},
		"B3": {35, 130},
	}
	near := NearbyBuoys(30, 120, positions, 10)
	if len(near) < 1 {
		t.Error("B1 and B2 should be nearby")
	}
}
