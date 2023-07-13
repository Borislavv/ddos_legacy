package vo

import "sync/atomic"

type Providers struct {
	Total  int64
	Active int64
}

func NewProviders(total int64) *Providers {
	if total <= 0 {
		total = 1
	}

	return &Providers{Total: total}
}

func (p *Providers) Activate() {
	atomic.AddInt64(&p.Active, 1)
}

func (p *Providers) Deactivate() {
	atomic.AddInt64(&p.Active, -1)
}
