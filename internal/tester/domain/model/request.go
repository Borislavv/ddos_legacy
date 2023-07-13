package model

import (
	"github.com/Borislavv/ddos/internal/tester/domain/vo"
	"net/http"
)

type Request struct {
	Request   *http.Request
	Timestamp vo.Timestamp
}

func NewRequest(request *http.Request) *Request {
	return &Request{Request: request}
}
