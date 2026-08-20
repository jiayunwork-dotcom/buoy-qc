package obs

import "math"

// GroupBy groups readings using a key function.
func GroupBy(readings []Reading, keyFn func(Reading) string) map[string][]Reading {
	m := make(map[string][]Reading)
	for _, r := range readings {
		k := keyFn(r)
		m[k] = append(m[k], r)
	}
	return m
}

// MaxField returns the maximum value of a field across readings.
func MaxField(readings []Reading, sel func(Reading) float64) float64 {
	if len(readings) == 0 {
		return 0
	}
	mx := sel(readings[0])
	for _, r := range readings[1:] {
		v := sel(r)
		if v > mx {
			mx = v
		}
	}
	return mx
}

// MinField returns the minimum value of a field across readings.
func MinField(readings []Reading, sel func(Reading) float64) float64 {
	if len(readings) == 0 {
		return 0
	}
	mn := sel(readings[0])
	for _, r := range readings[1:] {
		v := sel(r)
		if v < mn {
			mn = v
		}
	}
	return mn
}

// MeanField returns the mean value of a field across readings.
func MeanField(readings []Reading, sel func(Reading) float64) float64 {
	if len(readings) == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range readings {
		sum += sel(r)
	}
	return sum / float64(len(readings))
}

// StdField returns the population standard deviation of a field.
func StdField(readings []Reading, sel func(Reading) float64) float64 {
	n := len(readings)
	if n == 0 {
		return 0
	}
	mean := MeanField(readings, sel)
	var varSum float64
	for _, r := range readings {
		d := sel(r) - mean
		varSum += d * d
	}
	return math.Sqrt(varSum / float64(n))
}

// CountAbove returns the number of readings where field > threshold.
func CountAbove(readings []Reading, sel func(Reading) float64, threshold float64) int {
	n := 0
	for _, r := range readings {
		if sel(r) > threshold {
			n++
		}
	}
	return n
}

// CountBelow returns the number of readings where field < threshold.
func CountBelow(readings []Reading, sel func(Reading) float64, threshold float64) int {
	n := 0
	for _, r := range readings {
		if sel(r) < threshold {
			n++
		}
	}
	return n
}
