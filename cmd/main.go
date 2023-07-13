package main

import (
	"github.com/Borislavv/ddos/internal/tester/app/service"
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"log"
	"os"
	"os/signal"
	"runtime"
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

	tester := service.NewTester(settings, tasksCh, errorsCh, stopProvidersCh, stopConsumersCh)
	tester.Start()

	for {
		select {
		case <-osSigsCh:
			log.Println("user init. interrupting...")
			tester.Stop()
			time.Sleep(time.Second * 3)
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
