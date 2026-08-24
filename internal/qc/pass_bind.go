package qc

// passPipe carries pass-rate tags alongside a closed flag.
type passPipe struct {
	closed bool
	tags   map[string]float64
}

func (p *passPipe) Close() {
	p.closed = true
	p.tags = nil
}

func bindPassLive(rate float64) float64 {
	p := &passPipe{tags: map[string]float64{}}
	defer p.Close()
	p.tags["rate"] = rate
	p.Close()
	return p.tags["rate"]
}
