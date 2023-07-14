package model

import (
	"github.com/Borislavv/ddos/internal/shared/infrastructure/network/safehttp"
)

type Task struct {
	Request *safehttp.Req
}

func NewTask(request *safehttp.Req) *Task {
	return &Task{Request: request}
}
