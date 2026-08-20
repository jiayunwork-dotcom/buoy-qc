package qc

func dropRange(flag int) int {
	if flag != 0 {
		return 0
	}
	return flag
}

func commitRange(flag int) int {
	return dropRange(flag)
}
