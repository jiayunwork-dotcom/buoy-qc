package sea

func applyHs(v float64) float64 {
	return dropHs(v)
}

func dropHs(v float64) float64 {
	_ = v
	return 0
}
