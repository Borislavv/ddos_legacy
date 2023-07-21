package service

import (
	"errors"
	"fmt"
	"github.com/Borislavv/ddos/internal/shared/infrastructure/network/safehttp"
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Provider struct {
	mu        *sync.Mutex
	wg        *sync.WaitGroup
	wgInt     *sync.WaitGroup
	displayer IDisplayer
	settings  *model.Settings
	tasksCh   chan *model.Task
	stopCh    chan struct{}
	errorsCh  chan error
}

func NewProvider(
	wg *sync.WaitGroup,
	displayer IDisplayer,
	settings *model.Settings,
	tasksCh chan *model.Task,
	errorsCh chan error,
) *Provider {
	return &Provider{
		mu:        &sync.Mutex{},
		wg:        wg,
		wgInt:     &sync.WaitGroup{},
		displayer: displayer,
		settings:  settings,
		tasksCh:   tasksCh,
		errorsCh:  errorsCh,
		stopCh:    make(chan struct{}),
	}
}

func (p *Provider) Provide() {
	p.displayer.Display("starting #%d providers...", p.settings.Providers.Total)

	for i := int64(1); i <= p.settings.Providers.Total; i++ {
		go func(i int64) {
			defer func() {
				atomic.AddInt64(&p.settings.Providers.Active, -1)
				p.wg.Done()
				p.wgInt.Done()
				p.displayer.Display("\t - #%d provider stopped", i)
			}()
			atomic.AddInt64(&p.settings.Providers.Active, 1)
			p.wg.Add(1)
			p.wgInt.Add(1)
			p.displayer.Display("\t - #%d provider started", i)

			n := i * 1000000
			for {
				select {
				case <-p.stopCh:
					return
				default:
					timestamp := time.Now()

					request, err := http.NewRequest(
						"GET", // https://seo-wp1xv3n-6558-swoole.lux.kube.xbet.lan
						fmt.Sprintf(
							"https://seo-master-timings.lux.kube.xbet.lan/api/v1/pagedata?group_id=455&ref_id=1&url=https://1xlite-522762%d.top/registration&geo=cy&language=ru",
							n,
						),
						nil,
					)
					if err != nil {
						p.errorsCh <- errors.New("unable to create request: " + err.Error())
						continue
					}
					p.tasksCh <- model.NewTask(safehttp.NewReq(request), timestamp)
					n++
				}
			}
		}(i)
	}
}

func (p *Provider) Stop() {
	go func() {
		p.wg.Add(1)
		p.wgInt.Add(1)
		defer p.wg.Done()
		defer p.wgInt.Done()

		p.mu.Lock()
		activeProviders := p.settings.Providers.Active
		p.mu.Unlock()

		p.displayer.Display("stopping #%d providers...", activeProviders)
		for i := int64(1); i <= activeProviders; i++ {
			p.stopCh <- struct{}{}
		}
	}()

	p.wgInt.Wait()
}
