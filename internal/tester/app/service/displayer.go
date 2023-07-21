package service

import (
	"fmt"
	"github.com/Borislavv/ddos/internal/shared/infrastructure/channel"
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"sync"
	"sync/atomic"
)

type Displayer struct {
	wg     *sync.WaitGroup
	wgInt  *sync.WaitGroup
	cfg    *model.Settings
	infoCh *channel.SyncCh
	stopCh chan struct{}
}

func NewDisplayer(wg *sync.WaitGroup, buffer int64, cfg *model.Settings) *Displayer {
	return &Displayer{
		wg:     wg,
		wgInt:  &sync.WaitGroup{},
		cfg:    cfg,
		infoCh: channel.NewSyncCh(make(chan interface{}, buffer)),
		stopCh: make(chan struct{}),
	}
}

func (d *Displayer) Display(pattern string, args ...interface{}) {
	d.infoCh.Write(fmt.Sprintf(pattern, args...))
}

func (d *Displayer) DisplayError(err error) {
	d.infoCh.Write(err.Error())
}

func (d *Displayer) Start() {
	for i := int64(1); i <= d.cfg.Displayers.Total; i++ {
		go func(i int64) {
			defer func() {
				d.wg.Done()
				d.wgInt.Done()
				atomic.AddInt64(&d.cfg.Displayers.Active, -1)
			}()
			d.wg.Add(1)
			d.wgInt.Add(1)
			atomic.AddInt64(&d.cfg.Displayers.Active, 1)

			for {
				select {
				case <-d.stopCh:
					fmt.Printf("[displayer #%d] displayer stopped\n", i)
					return
				case msg := <-d.infoCh.Get():
					fmt.Printf("%s\n", msg)
				}
			}
		}(i)
	}
}

func (d *Displayer) Stop() {
	go func() {
		defer d.wg.Done()
		defer d.wgInt.Done()
		d.wg.Add(1)
		d.wgInt.Add(1)

		for i := int64(1); i <= d.cfg.Displayers.Total; i++ {
			d.stopCh <- struct{}{}
		}

		d.infoCh.Close()
	}()

	d.wgInt.Wait()
}
