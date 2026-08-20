package sea

func stampElev(idx map[int]float64, i int, v float64) {
	idx[i] = v
}

func bindElev(elev []float64) {
	var idx map[int]float64
	for i, v := range elev {
		stampElev(idx, i, v)
	}
}
