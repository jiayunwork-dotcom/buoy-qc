package obs

import "math"

// ExtractField extracts a float64 field from readings using a selector function.
func ExtractField(readings []Reading, sel func(Reading) float64) []float64 {
	out := make([]float64, len(readings))
	for i, r := range readings {
		out[i] = sel(r)
	}
	return out
}

// Resample resamples a series to a new length using linear interpolation.
func Resample(data []float64, newLen int) []float64 {
	n := len(data)
	if n == 0 || newLen <= 0 {
		return nil
	}
	if newLen == 1 {
		return []float64{data[0]}
	}
	out := make([]float64, newLen)
	for i := 0; i < newLen; i++ {
		pos := float64(i) * float64(n-1) / float64(newLen-1)
		lo := int(math.Floor(pos))
		hi := lo + 1
		if hi >= n {
			hi = n - 1
		}
		frac := pos - float64(lo)
		out[i] = data[lo]*(1-frac) + data[hi]*frac
	}
	return out
}

// Interpolate fills NaN values (represented by the sentinel) using linear interpolation.
// sentinel is the value considered as missing.
func Interpolate(data []float64, sentinel float64) []float64 {
	n := len(data)
	if n == 0 {
		return nil
	}
	out := make([]float64, n)
	copy(out, data)

	for i := 0; i < n; i++ {
		if out[i] == sentinel {
			// find previous valid
			prev := -1
			for j := i - 1; j >= 0; j-- {
				if out[j] != sentinel {
					prev = j
					break
				}
			}
			// find next valid
			next := -1
			for j := i + 1; j < n; j++ {
				if out[j] != sentinel {
					next = j
					break
				}
			}
			if prev >= 0 && next >= 0 {
				frac := float64(i-prev) / float64(next-prev)
				out[i] = out[prev] + frac*(out[next]-out[prev])
			} else if prev >= 0 {
				out[i] = out[prev]
			} else if next >= 0 {
				out[i] = out[next]
			}
		}
	}
	return out
}

// Downsample reduces resolution by averaging blocks of size factor.
func Downsample(data []float64, factor int) []float64 {
	if len(data) == 0 || factor <= 0 {
		return nil
	}
	n := len(data) / factor
	if n == 0 {
		n = 1
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		start := i * factor
		end := start + factor
		if end > len(data) {
			end = len(data)
		}
		sum := 0.0
		for j := start; j < end; j++ {
			sum += data[j]
		}
		out[i] = sum / float64(end-start)
	}
	return out
}

// Diff computes first-order differences: out[i] = data[i+1] - data[i].
func Diff(data []float64) []float64 {
	if len(data) < 2 {
		return nil
	}
	out := make([]float64, len(data)-1)
	for i := 0; i < len(data)-1; i++ {
		out[i] = data[i+1] - data[i]
	}
	return out
}

// CumulativeSum computes the running cumulative sum.
func CumulativeSum(data []float64) []float64 {
	if len(data) == 0 {
		return nil
	}
	out := make([]float64, len(data))
	out[0] = data[0]
	for i := 1; i < len(data); i++ {
		out[i] = out[i-1] + data[i]
	}
	return out
}
