package model

import "net/http"

type TaskInterface interface {
	GetRequest() *http.Request
}
