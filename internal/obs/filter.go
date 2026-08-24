package obs

// Filter returns readings that satisfy the predicate.
func Filter(readings []Reading, pred func(Reading) bool) []Reading {
	if len(readings) == 0 || pred == nil {
		return nil
	}
	var out []Reading
	for _, r := range readings {
		if pred(r) {
			out = append(out, r)
		}
	}
	return out
}

// FilterByBuoy returns readings for a specific buoy code.
func FilterByBuoy(readings []Reading, buoy string) []Reading {
	return Filter(readings, func(r Reading) bool { return r.Buoy == buoy })
}

// FilterByWindSpeed returns readings where wind speed is within [lo, hi].
func FilterByWindSpeed(readings []Reading, lo, hi float64) []Reading {
	return Filter(readings, func(r Reading) bool {
		return r.WindSpd >= lo && r.WindSpd <= hi
	})
}

// FilterByWaveHeight returns readings where wave height is within [lo, hi].
func FilterByWaveHeight(readings []Reading, lo, hi float64) []Reading {
	return Filter(readings, func(r Reading) bool {
		return r.WaveHt >= lo && r.WaveHt <= hi
	})
}

// FilterBySST returns readings where SST is within [lo, hi].
func FilterBySST(readings []Reading, lo, hi float64) []Reading {
	return Filter(readings, func(r Reading) bool {
		return r.SST >= lo && r.SST <= hi
	})
}

// BuoyCodes returns unique buoy codes in order of first appearance.
func BuoyCodes(readings []Reading) []string {
	seen := map[string]bool{}
	var codes []string
	for _, r := range readings {
		if !seen[r.Buoy] {
			seen[r.Buoy] = true
			codes = append(codes, r.Buoy)
		}
	}
	return codes
}

// Count returns the number of readings satisfying the predicate.
func Count(readings []Reading, pred func(Reading) bool) int {
	n := 0
	for _, r := range readings {
		if pred(r) {
			n++
		}
	}
	return n
}

// SummaryStats holds basic statistics for a Reading field.
type SummaryStats struct {
	Min, Max, Mean float64
	Count          int
}

// FieldStats computes min/max/mean for a field across readings.
func FieldStats(readings []Reading, sel func(Reading) float64) SummaryStats {
	if len(readings) == 0 {
		return SummaryStats{}
	}
	vals := ExtractField(readings, sel)
	s := SummaryStats{Count: len(vals), Min: vals[0], Max: vals[0]}
	sum := 0.0
	for _, v := range vals {
		sum += v
		if v < s.Min {
			s.Min = v
		}
		if v > s.Max {
			s.Max = v
		}
	}
	s.Mean = sum / float64(s.Count)
	return s
}
