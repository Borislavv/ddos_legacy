package vo

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
