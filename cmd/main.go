package main

import (
	"github.com/Borislavv/ddos/internal/tester/app/service"
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(7)

	start := time.Now()

	settings := model.NewSettings(15, 5, 2, time.Minute*5)

	tasksCh := make(chan *model.Task, settings.Workers.Total)
	defer close(tasksCh)
	errorsCh := make(chan error)
	stopProvidersCh := make(chan struct{})
	defer close(stopProvidersCh)
	stopConsumersCh := make(chan struct{})
	defer close(stopConsumersCh)
	osSigsCh := make(chan os.Signal, 1)
	defer close(osSigsCh)
	signal.Notify(osSigsCh, os.Interrupt)

	wg := &sync.WaitGroup{}
	displayer := service.NewDisplayer(wg, 1000)
	provider := service.NewProvider(wg, displayer, settings, tasksCh)
	consumer := service.NewConsumer(wg, displayer, settings, tasksCh)
	tester := service.NewTester(displayer, consumer, provider, settings, wg, tasksCh, errorsCh)
	tester.Start()

	for {
		select {
		case <-osSigsCh:
			displayer.Display("user init. interrupting...")
			tester.Stop()
			return
		default:
			if time.Since(start) > time.Second*14 {
				log.Println("stop providing due to timeout...")
				tester.Stop()
				return
			}
			log.Println("waiting")
			time.Sleep(time.Millisecond * 500)
		}
	}
}
