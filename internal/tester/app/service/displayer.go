package service

import (
	"fmt"
	"github.com/Borislavv/ddos/internal/shared/infrastructure/channel"
	"sync"
)

type Displayer struct {
	wg     *sync.WaitGroup
	wgInt  *sync.WaitGroup
	infoCh *channel.SyncCh
}

func NewDisplayer(wg *sync.WaitGroup, buffer int64) *Displayer {
	return &Displayer{
		wg:     wg,
		wgInt:  &sync.WaitGroup{},
		infoCh: channel.NewSyncCh(make(chan interface{}, buffer)),
	}
}

func (d *Displayer) Display(pattern string, args ...interface{}) {
	d.infoCh.Write(fmt.Sprintf(pattern, args...))
}

func (d *Displayer) DisplayError(err error) {
	d.infoCh.Write(err.Error())
}

func (d *Displayer) Start() {
	go func() {
		defer d.wg.Done()
		defer d.wgInt.Done()
		d.wg.Add(1)
		d.wgInt.Add(1)

		for msg := range d.infoCh.Get() {
			fmt.Println(msg)
		}
	}()
}

func (d *Displayer) Stop() {
	go func() {
		defer d.wg.Done()
		defer d.wgInt.Done()
		d.wg.Add(1)
		d.wgInt.Add(1)
		d.infoCh.Close()
	}()

	d.wgInt.Wait()
}
