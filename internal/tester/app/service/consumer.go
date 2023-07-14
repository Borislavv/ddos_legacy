package service

import (
	"errors"
	"github.com/Borislavv/ddos/internal/shared/infrastructure/network/safehttp"
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Consumer struct {
	mu              *sync.Mutex
	wg              *sync.WaitGroup
	wgInt           *sync.WaitGroup
	displayer       IDisplayer
	settings        *model.Settings
	httpClientsPool *safehttp.Pool
	tasksCh         chan *model.Task
	stopCh          chan struct{}
	errorsCh        chan error
}

func NewConsumer(
	wg *sync.WaitGroup,
	displayer IDisplayer,
	settings *model.Settings,
	tasksCh chan *model.Task,
	errorsCh chan error,
) *Consumer {
	return &Consumer{
		mu:              &sync.Mutex{},
		wg:              wg,
		wgInt:           &sync.WaitGroup{},
		displayer:       displayer,
		settings:        settings,
		httpClientsPool: safehttp.NewPool(25, time.Second*5),
		tasksCh:         tasksCh,
		errorsCh:        errorsCh,
		stopCh:          make(chan struct{}),
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
					client, err := c.httpClientsPool.GetAvailable()
					if err != nil {
						c.errorsCh <- errors.New("consumer stopped due to error: " + err.Error())
						continue
					}
					resp, err := client.Do(task.Request)
					if err != nil {
						c.errorsCh <- errors.New("consumer stopped due to error: " + err.Error())
						continue
					}

					c.displayer.Display("resp: %s", resp.Status())

					if err = resp.Close(); err != nil {
						c.errorsCh <- errors.New("error occurred while closing resp.Body: " + err.Error())
						continue
					}
				default:
					runtime.Gosched()
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
