package service

import (
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
}

func NewProvider(
	wg *sync.WaitGroup,
	displayer IDisplayer,
	settings *model.Settings,
	tasksCh chan *model.Task,
) *Provider {
	return &Provider{
		mu:        &sync.Mutex{},
		wg:        wg,
		wgInt:     &sync.WaitGroup{},
		displayer: displayer,
		settings:  settings,
		tasksCh:   tasksCh,
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

			for {
				select {
				case <-p.stopCh:
					return
				default:
					p.tasksCh <- model.NewTask(model.NewRequest(&http.Request{}))
					time.Sleep(time.Millisecond * 100)
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
