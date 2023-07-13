package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(7)

	wg := &sync.WaitGroup{}

	stopCh := make(chan struct{})
	defer close(stopCh)
	doneCh := make(chan struct{}, 1)
	defer close(doneCh)
	dataCh := make(chan struct{}, 5)
	defer close(dataCh)

	provide(wg, stopCh, dataCh)
	consume(wg, doneCh, dataCh)

	osCh := make(chan os.Signal, 1)
	signal.Notify(osCh, os.Interrupt)

	i := 0
loop:
	for {
		select {
		case <-osCh:
			log.Println("interrupted by user")
			stop(wg, stopCh, doneCh)
			break loop
		default:
			log.Println("awaiting in main thread....")
			time.Sleep(time.Millisecond * 10)
			i++
		}
	}

	log.Printf("iterations #%d", i)
}

func provide(wg *sync.WaitGroup, stopCh chan struct{}, dataCh chan struct{}) {
	go func() {
		defer wg.Done()
		wg.Add(1)

		for {
			select {
			case <-stopCh:
				log.Println("provider stopped by stopCh")
				return
			default:
				dataCh <- struct{}{}
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()
}

func consume(wg *sync.WaitGroup, doneCh chan struct{}, dataCh chan struct{}) {
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			wg.Add(1)

			for {
				select {
				case <-doneCh:
					log.Println("consumer stopped by doneCh")
					return
				case data := <-dataCh:
					log.Printf("data received: %+v\n", data)
				}
			}
		}()
	}
}

func stop(wg *sync.WaitGroup, stopCh chan struct{}, doneCh chan struct{}) {
	go func() {
		wg.Add(1)
		defer func() {
			wg.Done()
		}()

		stopCh <- struct{}{}
		for i := 0; i < 5; i++ {
			doneCh <- struct{}{}
		}
	}()

	wg.Wait()
}
