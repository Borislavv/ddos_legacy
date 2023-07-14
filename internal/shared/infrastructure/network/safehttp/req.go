package safehttp

import (
	"github.com/Borislavv/ddos/internal/tester/domain/vo"
	"net/http"
)

type Req struct {
	Request   *http.Request
	IsFailed  bool
	Timestamp vo.Timestamp
}

func NewReq(request *http.Request) *Req {
	return &Req{Request: request}
}

func (r *Req) MarkFailed() {
	r.IsFailed = true
}
