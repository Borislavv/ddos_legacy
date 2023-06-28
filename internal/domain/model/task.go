package model

import (
	"net/http"
	"time"
)

type Task struct {
	request *http.Request
	workers struct {
		num   int
		isSet bool
	}
	requests struct {
		num   int64
		isSet bool
	}
	timeout struct {
		duration time.Duration
		isSet    bool
	}
	thresholds struct {
		low    time.Duration
		medium time.Duration
		high   time.Duration
		over   time.Duration
	}
}
