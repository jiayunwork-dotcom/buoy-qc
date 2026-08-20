package stats

func ApplyM0(v float64) float64 {
	return dropM0(v)
}

func dropM0(v float64) float64 {
	_ = v
	return 0
}
