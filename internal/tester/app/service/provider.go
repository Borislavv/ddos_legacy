package service

import (
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Provider struct {
	wg       *sync.WaitGroup
	settings *model.Settings
	tasksCh  chan *model.Task
	stopCh   chan struct{}
}

func NewProvider(
	settings *model.Settings,
	tasksCh chan *model.Task,
	stopCh chan struct{},
) *Provider {
	return &Provider{
		wg:       &sync.WaitGroup{},
		settings: settings,
		tasksCh:  tasksCh,
		stopCh:   stopCh,
	}
}

func (p *Provider) Provide() {
	log.Printf("starting #%d providers...\n", p.settings.Providers.Total)

	for i := int64(1); i <= p.settings.Providers.Total; i++ {
		go func(i int64) {
			defer func() {
				atomic.AddInt64(&p.settings.Providers.Active, -1)
				p.wg.Done()
				log.Printf("\t - #%d provider stopped\n", i)
			}()
			atomic.AddInt64(&p.settings.Providers.Active, 1)
			p.wg.Add(1)
			log.Printf("\t - #%d provider started\n", i)

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
		defer p.wg.Done()

		log.Printf("stopping #%d providers...\n", p.settings.Providers.Active)
		for i := int64(1); i <= p.settings.Providers.Total; i++ {
			p.stopCh <- struct{}{}
		}
	}()

	p.wg.Wait()
}
