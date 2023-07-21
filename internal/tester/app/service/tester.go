package service

import (
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"sync"
)

type Tester struct {
	displayer IDisplayer
	consumer  IConsumer
	provider  IProvider
	meter     IMeter

	settings *model.Settings

	wg      *sync.WaitGroup
	tasksCh chan *model.Task
}

func NewTester(
	displayer IDisplayer,
	consumer IConsumer,
	provider IProvider,
	meter IMeter,
	settings *model.Settings,
	wg *sync.WaitGroup,
	tasksCh chan *model.Task,
) *Tester {
	return &Tester{
		wg:        wg,
		displayer: displayer,
		meter:     meter,
		settings:  settings,
		tasksCh:   tasksCh,
		consumer:  consumer,
		provider:  provider,
	}
}

func (t *Tester) Start() {
	go func() {
		defer t.wg.Done()
		t.wg.Add(1)

		t.meter.Start()
		t.displayer.Start()
		t.provider.Provide()
		t.consumer.Consume()
	}()
}

func (t *Tester) Stop() {
	go func() {
		defer t.wg.Done()
		t.wg.Add(1)

		t.meter.Stop()
		t.provider.Stop()
		t.consumer.Stop()
		t.displayer.Stop()
		t.meter.Summary()
	}()

	t.wg.Wait()
}
