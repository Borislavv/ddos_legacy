package vo

import (
	"errors"
	"log"
	"time"
)

type Duration struct {
	value  time.Duration
	isInit bool
}

func NewDuration(duration time.Duration) *Duration {
	dur := &Duration{}
	if duration > 0 {
		dur.value = duration
		dur.isInit = true
	}

	return dur
}

func (t *Duration) GetDuration() time.Duration {
	if !t.isInit {
		log.Fatalln("'duration' value was not initialized yet, probably 'from' was not set up")
	}
	return t.value
}

func (t *Duration) SetDuration(duration time.Duration) error {
	if t.isInit {
		return errors.New("'duration' value already initialized")
	}

	t.value = duration
	t.isInit = true

	return nil
}

func (t *Duration) IsSetDuration() bool {
	return t.isInit
}
