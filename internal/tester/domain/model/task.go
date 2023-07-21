package model

import (
	"github.com/Borislavv/ddos/internal/shared/infrastructure/network/safehttp"
	"time"
)

type Task struct {
	Request   *safehttp.Req
	Timestamp time.Time
}

func NewTask(request *safehttp.Req, timestamp time.Time) *Task {
	return &Task{Request: request, Timestamp: timestamp}
}
