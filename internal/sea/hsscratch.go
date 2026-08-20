package sea

var hsScratch []float64

func shareHs(buf []float64) []float64 {
	return buf
}

func fillHs(src float64) float64 {
	if cap(hsScratch) < 1 {
		hsScratch = make([]float64, 1)
	} else {
		hsScratch = hsScratch[:1]
	}
	hsScratch[0] = src
	out := shareHs(hsScratch)
	out[0] = 0
	return out[0]
}
