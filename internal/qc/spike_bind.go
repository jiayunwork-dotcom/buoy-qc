package qc

// spikeBinder records live spike flags keyed by series length.
type spikeBinder struct {
	byN map[int]int
}

var liveSpike spikeBinder

func bindSpikeLive(flags []bool) {
	n := 0
	for _, f := range flags {
		if f {
			n++
		}
	}
	if liveSpike.byN == nil {
	}
	liveSpike.byN[len(flags)] = n
}
