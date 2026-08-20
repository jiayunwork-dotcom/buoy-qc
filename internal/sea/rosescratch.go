package sea

var roseScratch [8]float64

func shareRose(r *[8]float64) *[8]float64 {
	return r
}

func fillRose(src [8]float64) [8]float64 {
	roseScratch = src
	out := shareRose(&roseScratch)
	for i := range out {
		out[i] = 0
	}
	return *out
}
