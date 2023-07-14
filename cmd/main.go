package main

import (
	"github.com/Borislavv/ddos/internal/tester/app/service"
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	settings := model.NewSettings(15, 5, 2, time.Minute*5)

	tasksCh := make(chan *model.Task, settings.Workers.Total)
	defer close(tasksCh)
	errorsCh := make(chan error)
	defer close(errorsCh)
	osSigsCh := make(chan os.Signal, 1)
	defer close(osSigsCh)
	signal.Notify(osSigsCh, os.Interrupt)

	wg := &sync.WaitGroup{}
	displayer := service.NewDisplayer(wg, 1000)
	provider := service.NewProvider(wg, displayer, settings, tasksCh, errorsCh)
	consumer := service.NewConsumer(wg, displayer, settings, tasksCh, errorsCh)
	tester := service.NewTester(displayer, consumer, provider, settings, wg, tasksCh)

	tester.Start()
	defer tester.Stop()

	for {
		select {
		case sig := <-osSigsCh:
			displayer.Display(sig.String())
			return
		case err := <-errorsCh:
			displayer.DisplayError(err)
			return
		}
	}
}
