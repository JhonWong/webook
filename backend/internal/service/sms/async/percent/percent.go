package percent

import "sync"

type Percent[T any] struct {
	cap         int
	index       int
	molecular   int
	size        int
	arr         []T
	isMolecular func(val T) bool

	mu sync.Mutex
}

func NewPercent[T any](cap int, isMolecular func(T) bool) *Percent[T] {
	return &Percent[T]{
		cap:         cap,
		arr:         make([]T, cap),
		isMolecular: isMolecular,
	}
}

func (p *Percent[T]) Add(val T) float32 {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.size == p.cap && p.isMolecular(p.arr[p.index]) {
		p.molecular--
	}

	p.arr[p.index] = val
	p.index = (p.index + 1) % p.cap
	if p.isMolecular(val) {
		p.molecular++
	}
	if p.size < p.cap {
		p.size++
	}

}
