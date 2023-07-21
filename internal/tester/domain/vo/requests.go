package vo

type Requests struct {
	Total int64
	IsSet bool
}

func NewRequests(total int64) *Requests {
	reqs := &Requests{}

	if total > 0 {
		reqs.Total = total
		reqs.IsSet = true
	}

	return reqs
}
