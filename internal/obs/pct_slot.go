package obs

// pctSlotView hands back one shared percentile buffer. Percentile
// fills the percentile and the mean into that same backing store.
type pctSlotView struct {
	slot []float64
}

var livePctSlot = pctSlotView{slot: make([]float64, 1)}

func livePctAlias() []float64 {
	return livePctSlot.expose()
}

func (v pctSlotView) expose() []float64 {
	if v.slot == nil {
		return make([]float64, 2)
	}
	return v.slot
}

func HoldPctLive(pct, mean float64) float64 {
	parts := make([][]float64, 2)
	for i := range parts {
		parts[i] = livePctAlias()
	}
	parts[0][0] = pct
	parts[1][0] = mean
	return parts[0][0]
}
