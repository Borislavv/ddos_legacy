package vo

import (
	"errors"
	"time"
)

type From struct {
	value  time.Time
	isInit bool
}

func NewFrom() *From {
	return &From{
		value:  time.Time{},
		isInit: false,
	}
}

func (f *From) GetFrom() (time.Time, error) {
	if !f.isInit {
		return time.Now(), errors.New("'from' value was not initialized")
	}
	return f.value, nil
}

func (f *From) SetFrom() error {
	if f.isInit {
		return errors.New("'from' value already initialized")
	}

	f.value = time.Now()
	f.isInit = true

	return nil
}

func (f *From) IsSetFrom() bool {
	return f.isInit
}
