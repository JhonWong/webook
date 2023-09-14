package percent

import "sync"

type Percent[T any] struct {
	size        int
	index       int
	molecular   int
	isFilled    bool
	arr         []T
	isMolecular func(val T) bool

	mu sync.Mutex
}

func NewPercent[T any](size int, isMolecular func(T) bool) *Percent[T] {
	return &Percent[T]{
		size:        size,
		arr:         make([]T, size),
		isMolecular: isMolecular,
	}
}

func (p *Percent[T]) Add(val T) float64 {
	p.mu.Lock()

	if p.isFilled && p.isMolecular(p.arr[p.index]) {
		p.molecular--
	}

	p.arr[p.index] = val
	p.index = (p.index + 1) % p.size
	if p.isMolecular(val) {
		p.molecular++
	}
	if p.index == 0 && !p.isFilled {
		p.isFilled = true
	}
	defer p.mu.Unlock()

	return p.Percent()
}

func (p *Percent[T]) Percent() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	var res float64
	if p.isFilled {
		res = float64(p.molecular) / float64(p.size)
	} else {
		res = float64(p.molecular) / float64(p.index)
	}
	return res
}
