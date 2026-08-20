package qc

import (
	"math"

	"buoy-qc/internal/obs"
)

// SpatialConsistency checks whether a reading is consistent with neighbors.
// A reading is flagged if its value deviates more than k * std(neighbors) from
// the mean of its neighbors. fieldSel selects which field to check.
func SpatialConsistency(target obs.Reading, neighbors []obs.Reading, fieldSel func(obs.Reading) float64, k float64) bool {
	if len(neighbors) == 0 {
		return false
	}
	vals := make([]float64, len(neighbors))
	sum := 0.0
	for i, n := range neighbors {
		vals[i] = fieldSel(n)
		sum += vals[i]
	}
	meanVal := sum / float64(len(vals))
	var varSum float64
	for _, v := range vals {
		d := v - meanVal
		varSum += d * d
	}
	std := math.Sqrt(varSum / float64(len(vals)))
	targetVal := fieldSel(target)
	if std == 0 {
		return targetVal != meanVal
	}
	return math.Abs(targetVal-meanVal) > k*std
}

// BuoyDistance computes approximate distance (km) between two buoy positions
// using the Haversine formula. lat/lon in degrees.
func BuoyDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // Earth radius km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	lat1R := lat1 * math.Pi / 180
	lat2R := lat2 * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1R)*math.Cos(lat2R)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// NearbyBuoys returns buoy codes within maxDist km of (lat, lon).
// positions maps buoy code to [lat, lon].
func NearbyBuoys(lat, lon float64, positions map[string][2]float64, maxDist float64) []string {
	var near []string
	for code, pos := range positions {
		d := BuoyDistance(lat, lon, pos[0], pos[1])
		if d <= maxDist {
			near = append(near, code)
		}
	}
	return near
}
