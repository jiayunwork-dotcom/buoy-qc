package stats

// gumbelHoldView hands back one shared return-level buffer. ExtremeReturn
// fills the Gumbel level and scale into that same backing store.
type gumbelHoldView struct {
	slot []float64
}

var liveGumbelHold = gumbelHoldView{slot: make([]float64, 1)}

func liveGumbelAlias() []float64 {
	return liveGumbelHold.expose()
}

func (v gumbelHoldView) expose() []float64 {
	if v.slot == nil {
		return make([]float64, 2)
	}
	return v.slot
}

func HoldGumbelLive(level, scale float64) float64 {
	parts := make([][]float64, 2)
	for i := range parts {
		parts[i] = liveGumbelAlias()
	}
	parts[0][0] = level
	parts[1][0] = scale
	return parts[0][0]
}
