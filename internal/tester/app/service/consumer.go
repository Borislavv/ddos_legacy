package service

import (
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"sync"
	"sync/atomic"
)

type Consumer struct {
	mu        *sync.Mutex
	wg        *sync.WaitGroup
	wgInt     *sync.WaitGroup
	displayer IDisplayer
	settings  *model.Settings
	stopCh    chan struct{}
	tasksCh   chan *model.Task
}

func NewConsumer(
	wg *sync.WaitGroup,
	displayer IDisplayer,
	settings *model.Settings,
	tasksCh chan *model.Task,
) *Consumer {
	return &Consumer{
		mu:        &sync.Mutex{},
		wg:        wg,
		wgInt:     &sync.WaitGroup{},
		displayer: displayer,
		settings:  settings,
		tasksCh:   tasksCh,
		stopCh:    make(chan struct{}),
	}
}

func (c *Consumer) Consume() {
	c.displayer.Display("starting #%d consumers...", c.settings.Workers.Total)

	for i := int64(1); i <= c.settings.Workers.Total; i++ {
		go func(i int64) {
			defer func() {
				atomic.AddInt64(&c.settings.Workers.Active, -1)
				c.wg.Done()
				c.wgInt.Done()
				c.displayer.Display("\t - #%d consumer stopped", i)
			}()
			atomic.AddInt64(&c.settings.Workers.Active, 1)
			c.wg.Add(1)
			c.wgInt.Add(1)
			c.displayer.Display("\t - #%d consumer started", i)

			for {
				select {
				case <-c.stopCh:
					return
				case task := <-c.tasksCh:
					c.displayer.Display("worker #%d, received task: %+v", i, task)
				}
			}
		}(i)
	}
}

func (c *Consumer) Stop() {
	go func() {
		defer c.wg.Done()
		defer c.wgInt.Done()
		c.wg.Add(1)
		c.wgInt.Add(1)

		c.mu.Lock()
		activeWorkers := c.settings.Workers.Active
		c.mu.Unlock()

		c.displayer.Display("stopping #%d consumers...", activeWorkers)
		for i := int64(0); i < activeWorkers; i++ {
			c.stopCh <- struct{}{}
		}
	}()

	c.wgInt.Wait()
}
