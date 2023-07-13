package service

import (
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"log"
	"sync"
	"sync/atomic"
)

type Consumer struct {
	wg       *sync.WaitGroup
	settings *model.Settings
	stopCh   chan struct{}
	tasksCh  chan *model.Task
}

func NewConsumer(
	settings *model.Settings,
	tasksCh chan *model.Task,
	stopCh chan struct{},
) *Consumer {
	return &Consumer{
		wg:       &sync.WaitGroup{},
		settings: settings,
		tasksCh:  tasksCh,
		stopCh:   stopCh,
	}
}

func (c *Consumer) Consume() {
	log.Printf("starting #%d consumers...\n", c.settings.Workers.Total)

	for i := int64(1); i <= c.settings.Workers.Total; i++ {
		go func(i int64) {
			defer func() {
				atomic.AddInt64(&c.settings.Workers.Active, -1)
				c.wg.Done()
				log.Printf("\t - #%d consumer stopped\n", i)
			}()
			atomic.AddInt64(&c.settings.Workers.Active, 1)
			c.wg.Add(1)

			log.Printf("\t - #%d consumer started\n", i)

			for {
				select {
				case <-c.stopCh:
					return
				case data := <-c.tasksCh:
					log.Printf("data received: %+v\n", data)
				}
			}
		}(i)
	}
}

func (c *Consumer) Stop() {
	go func() {
		defer c.wg.Done()
		c.wg.Add(1)

		log.Printf("stopping #%d consumers...\n", c.settings.Workers.Active)
		for i := int64(0); i < c.settings.Workers.Active; i++ {
			c.stopCh <- struct{}{}
		}
	}()

	c.wg.Wait()

	if c.settings.Workers.Active > 0 {
		log.Fatalf(
			"'consumer' workers does not stopped properly, number of active workers #%d",
			c.settings.Workers.Active,
		)
	}
}
