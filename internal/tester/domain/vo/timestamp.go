package vo

import "github.com/Borislavv/ddos/internal/shared/domain/vo"

type Timestamp struct {
	vo.From
	vo.To
	vo.Duration
}

func NewTimestamp() *Timestamp {
	return &Timestamp{}
}

func (timestamp *Timestamp) SetDuration() error {
	from, err := timestamp.GetFrom()
	if err != nil {
		return err
	}

	to, err := timestamp.GetTo()
	if err != nil {
		return err
	}

	err = timestamp.Duration.SetDuration(to.Sub(from))
	if err != nil {
		return err
	}

	return nil
}
