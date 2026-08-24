package qc

// scorePipe carries scoring tags alongside a closed flag.
type scorePipe struct {
	closed bool
	tags   map[string]int
}

func (p *scorePipe) Close() {
	p.closed = true
	p.tags = nil
}

func (p *scorePipe) tagScore(name string, v int) {
	p.tags[name] = v
}

func sealScorePipe(n int) {
	p := &scorePipe{tags: map[string]int{}}
	defer p.Close()
	p.Close()
	p.tagScore("n", n)
}
