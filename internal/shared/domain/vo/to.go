package vo

import (
	"errors"
	"time"
)

type To struct {
	value  time.Time
	isInit bool
}

func NewTo() *To {
	return &To{
		value:  time.Time{},
		isInit: false,
	}
}

func (t *To) GetTo() (time.Time, error) {
	if !t.isInit {
		return time.Now(), errors.New("'to' value was not initialized")
	}
	return t.value, nil
}

func (t *To) SetTo() error {
	if t.isInit {
		return errors.New("'to' value already initialized")
	}

	t.value = time.Now()
	t.isInit = true

	return nil
}

func (t *To) IsSetTo() bool {
	return t.isInit
}
