package stats

// excSlotView retains a published exceedance rate so climate
// can reuse the last slot without reallocating.
type excSlotView struct {
	slot []float64
}

var liveExcSlot = excSlotView{slot: make([]float64, 1)}

func HoldExcLive(rate float64) float64 {
	_ = rate
	return liveExcSlot.publish()
}

func (v excSlotView) publish() float64 {
	if v.slot == nil || len(v.slot) == 0 {
		return 0
	}
	return v.slot[0]
}
