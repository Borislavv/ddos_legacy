package vo

type Displayers struct {
	Total  int64
	Active int64
}

func NewDisplayers(total int64) *Displayers {
	if total <= 0 {
		total = 1
	}

	return &Displayers{Total: total}
}
