package sea

// hsLiveView hands back one shared term buffer. SignificantWaveHt
// fills Hs and variance into that same backing store.
type hsLiveView struct {
	slot []float64
}

var liveHsSlot = hsLiveView{slot: make([]float64, 1)}

func liveHsAlias() []float64 {
	return liveHsSlot.expose()
}

func (v hsLiveView) expose() []float64 {
	if v.slot == nil {
		return make([]float64, 2)
	}
	return v.slot
}

func holdHsLive(hs, variance float64) float64 {
	terms := make([][]float64, 2)
	for i := range terms {
		terms[i] = liveHsAlias()
	}
	terms[0][0] = hs
	terms[1][0] = variance
	return terms[0][0]
}
