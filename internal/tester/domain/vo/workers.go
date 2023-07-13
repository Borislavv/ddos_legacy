package vo

import "runtime"

type Workers struct {
	Total  int64
	Active int64
}

func NewWorkers(total int64) *Workers {
	if total <= 0 {
		total = int64(runtime.NumCPU())
	}

	return &Workers{Active: 0, Total: total}
}
