package model

import (
	sharedvo "github.com/Borislavv/ddos/internal/shared/domain/vo"
	"github.com/Borislavv/ddos/internal/tester/domain/vo"
	"time"
)

type Settings struct {
	Providers *vo.Providers
	Workers   *vo.Workers
	Requests  *vo.Requests
	Duration  *sharedvo.Duration
}

func NewSettings(requests int64, workers int64, providers int64, duration time.Duration) *Settings {
	return &Settings{
		Providers: vo.NewProviders(providers),
		Workers:   vo.NewWorkers(workers),
		Requests:  vo.NewRequests(requests),
		Duration:  sharedvo.NewDuration(duration),
	}
}