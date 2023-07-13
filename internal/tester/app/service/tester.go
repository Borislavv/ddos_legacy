package service

import (
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"log"
	"sync"
)

type Tester struct {
	wg       *sync.WaitGroup
	settings *model.Settings
	tasksCh  chan *model.Task
	errorsCh chan error
	consumer *Consumer
	provider *Provider
}

func NewTester(
	settings *model.Settings,
	tasksCh chan *model.Task,
	errorsCh chan error,
	stopProvidersCh chan struct{},
	stopConsumersCh chan struct{},
) *Tester {
	return &Tester{
		wg:       &sync.WaitGroup{},
		settings: settings,
		tasksCh:  tasksCh,
		errorsCh: errorsCh,
		consumer: NewConsumer(settings, tasksCh, stopConsumersCh),
		provider: NewProvider(settings, tasksCh, stopProvidersCh),
	}
}

func (t *Tester) Start() {
	go func() {
		defer t.wg.Done()
		t.wg.Add(1)

		t.provider.Provide()
		t.consumer.Consume()
	}()
}

func (t *Tester) Stop() {
	go func() {
		defer t.wg.Done()
		t.wg.Add(1)

		log.Println("stopping providers and consumers...")
		t.provider.Stop()
		t.consumer.Stop()
	}()

	t.wg.Wait()
}
