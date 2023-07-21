package service

import (
	"errors"
	"fmt"
	"github.com/Borislavv/ddos/internal/shared/infrastructure/network/safehttp"
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"github.com/Borislavv/ddos/internal/tester/infrastructure/helper"
	"net/url"
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
	meter           IMeter
	settings        *model.Settings
	httpClientsPool *safehttp.Pool
	tasksCh         chan *model.Task
	stopCh          chan struct{}
	errorsCh        chan error
	counter         int64
}

func NewConsumer(
	wg *sync.WaitGroup,
	displayer IDisplayer,
	meter IMeter,
	settings *model.Settings,
	tasksCh chan *model.Task,
	errorsCh chan error,
) *Consumer {
	return &Consumer{
		mu:              &sync.Mutex{},
		wg:              wg,
		wgInt:           &sync.WaitGroup{},
		displayer:       displayer,
		meter:           meter,
		settings:        settings,
		httpClientsPool: safehttp.NewPool(25, time.Second*1),
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
						c.errorsCh <- errors.New("consumer skipped iteration due to error: " + err.Error())
						continue
					}

					timestamp := time.Now()

					URL, err := url.Parse(fmt.Sprintf(task.Request.Request.URL.String()+"&timestamp=%d", timestamp.UnixMicro()))
					if err != nil {
						c.errorsCh <- errors.New("consumer skipped iteration due to error: " + err.Error())
						continue
					}
					task.Request.Request.URL = URL
					resp, err := client.Do(task.Request)
					if err != nil {
						c.errorsCh <- errors.New("consumer skipped iteration due to error: " + err.Error())
						continue
					}
					_, err = resp.Body()
					if err != nil {
						c.errorsCh <- errors.New("consumer skipped iteration due to error: " + err.Error())
						continue
					}
					c.meter.CommitReq(resp.Req())

					c.displayer.Display(
						fmt.Sprintf(
							"[#%d] [resp: %s] [dur: %s] [client->server: %s] [p: %s ] [req->resp: %s] [server->client: %s]",
							atomic.AddInt64(&c.counter, 1),
							resp.Status(),
							time.Since(timestamp),
							resp.Origin().Header.Get("Server-Timing-Client-To-Server-Duration"),
							resp.Origin().Header.Get("Server-Timing"),
							resp.Origin().Header.Get("Server-Timing-Request-To-Response-Duration"),
							time.Since(helper.ParsePhpMicroTime(resp.Origin().Header.Get("Server-Timing-Response-Timestamp"))),
							//helper.ParseMillisecondsDur(resp.Origin().Header.Get("Server-Timing-Client-To-Server-Duration")),
							//helper.ParsePDur(resp.Origin().Header.Get("Server-Timing")),
							//helper.ParseMillisecondsDur(resp.Origin().Header.Get("Server-Timing-Request-To-Response-Duration")),
							//time.Since(helper.ParsePhpMicroTime(resp.Origin().Header.Get("Server-Timing-Response-Timestamp"))),
						),
					)

					if err = resp.Close(); err != nil {
						c.errorsCh <- errors.New("consumer skipped iteration due to error occurred while closing resp.Body: " + err.Error())
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
