package service

import (
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"sync"
)

type Tester struct {
	displayer IDisplayer
	consumer  IConsumer
	provider  IProvider

	settings *model.Settings

	wg       *sync.WaitGroup
	tasksCh  chan *model.Task
	errorsCh chan error
}

func NewTester(
	displayer IDisplayer,
	consumer IConsumer,
	provider IProvider,
	settings *model.Settings,
	wg *sync.WaitGroup,
	tasksCh chan *model.Task,
	errorsCh chan error,
) *Tester {
	return &Tester{
		wg:        wg,
		displayer: displayer,
		settings:  settings,
		tasksCh:   tasksCh,
		errorsCh:  errorsCh,
		consumer:  consumer,
		provider:  provider,
	}
}

func (t *Tester) Start() {
	go func() {
		defer t.wg.Done()
		t.wg.Add(1)

		t.displayer.Start()
		t.provider.Provide()
		t.consumer.Consume()
	}()
}

func (t *Tester) Stop() {
	go func() {
		defer t.wg.Done()
		t.wg.Add(1)

		t.provider.Stop()
		t.consumer.Stop()
		t.displayer.Stop()
	}()

	t.wg.Wait()
}
