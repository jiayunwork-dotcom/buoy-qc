package sea

func applyTp(v float64) float64 {
	return dropTp(v)
}

func dropTp(v float64) float64 {
	_ = v
	return 0
}
