package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Stat struct {
	From     time.Time
	To       time.Time
	Duration time.Duration

	ReqsNum        int64
	SuccessReqs    int64
	FailedReqs     int64
	UnfinishedReqs int64

	ReqsDur        time.Duration
	SuccessReqsDur time.Duration
	FailedReqsDur  time.Duration

	AvgReqsDur        time.Duration
	AvgSuccessReqsDur time.Duration
	AvgFailedReqsDur  time.Duration
}

type ReqStat struct {
	Duration     time.Duration
	IsFailed     bool
	IsUnfinished bool
}

func provide() (chan struct{}, context.Context, context.CancelFunc, *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	dataCh := make(chan struct{}, 14)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for id := 1; id <= 2; id++ {
			wg.Add(1)
			go func(id int) {
				defer func() {
					wg.Done()
					log.Printf("p #%d is stopped\n", id)
				}()
				log.Printf("p #%d is started\n", id)

				for {
					select {
					case <-ctx.Done():
						return
					case dataCh <- struct{}{}:
					}
				}
			}(id)
		}
	}()
	return dataCh, ctx, cancel, wg
}

func stopProvide(dataCh chan struct{}, cancel context.CancelFunc, pwg *sync.WaitGroup) {
	cancel()
	pwg.Wait()
	close(dataCh)
	log.Println("awaiting closing of providers...")
}

func consume(dataCh chan struct{}, reqsStatCh chan<- *ReqStat, ctx context.Context) *sync.WaitGroup {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for id := 0; id < 5; id++ {
			wg.Add(1)
			go func(id int) {
				defer func() {
					wg.Done()
					log.Printf("consumer #%d is stopped", id)
				}()
				log.Printf("consumer #%d is started", id)
				for _ = range dataCh {
					func() {
						timestamp := time.Now()

						rand.Seed(timestamp.UnixNano())

						req, err := http.NewRequestWithContext(
							ctx,   //https://seo-wp1xv3n-swoole-debug.lux.kube.xbet.lan/
							"GET", // https://seo-wp1xv3n-6558-swoole.lux.kube.xbet.lan
							fmt.Sprintf( // https://seo-master-timings.lux.kube.xbet.lan
								"http://localhost:8080" +
									"/api/v1/pagedata?group_id=455&ref_id=1&url=https://1xlite-%d.top/registration&geo=cy&language=ru&timestamp=%d",
								rand.Int63(),
								timestamp.UnixMicro(),
							),
							nil,
						)
						if err != nil {
							log.Println(err)
						}

						reqStat := &ReqStat{}
						defer func() {
							reqStat.Duration = time.Since(timestamp)
							reqsStatCh <- reqStat
						}()

						resp, err := (&http.Client{}).Do(req)
						if err != nil {
							if strings.Contains(err.Error(), "context canceled") {
								reqStat.IsUnfinished = true
							} else {
								reqStat.IsFailed = true
								log.Println(err)
							}
							return
						}
						defer resp.Body.Close()

						if resp.StatusCode != 200 {
							bytes, errR := ioutil.ReadAll(resp.Body)
							if errR != nil {
								log.Println(err)
								reqStat.IsFailed = true
								return
							}

							log.Printf(
								"[dur: %s][resp: %s]\nResponse: %s",
								time.Since(timestamp),
								resp.Status,
								string(bytes),
							)
							reqStat.IsFailed = true
							return
						}

						log.Printf(
							"[dur: %s][resp: %s]\n",
							time.Since(timestamp),
							resp.Status,
						)
					}()
				}
			}(id)
		}
	}()

	return wg
}

func stopConsume(cwg *sync.WaitGroup, reqsStatCh chan *ReqStat) {
	log.Println("stopping consumers...")
	cwg.Wait()
	log.Println("stopping requests handling...")
	close(reqsStatCh)
}

func handleReqs() (chan *ReqStat, chan *Stat) {
	reqsStatCh := make(chan *ReqStat, 14)
	statCh := make(chan *Stat)

	stat := &Stat{
		From: time.Now(),
	}

	go func() {
		defer log.Println("stopped requests handling")

		for reqStat := range reqsStatCh {
			if reqStat.IsUnfinished {
				stat.UnfinishedReqs++
				continue
			}

			stat.ReqsNum++
			stat.ReqsDur += reqStat.Duration

			if reqStat.IsFailed {
				stat.FailedReqs++
				stat.FailedReqsDur += reqStat.Duration
			} else {
				stat.SuccessReqs++
				stat.SuccessReqsDur += reqStat.Duration
			}
		}

		stat.To = time.Now()
		stat.Duration = stat.To.Sub(stat.From)

		if stat.ReqsNum != 0 {
			stat.AvgReqsDur = time.Duration(stat.ReqsDur.Nanoseconds() / stat.ReqsNum)
		}

		if stat.SuccessReqs != 0 {
			stat.AvgSuccessReqsDur = time.Duration(stat.SuccessReqsDur.Nanoseconds() / stat.SuccessReqs)
		}

		if stat.FailedReqs != 0 {
			stat.AvgFailedReqsDur = time.Duration(stat.FailedReqsDur.Nanoseconds() / stat.FailedReqs)
		}

		statCh <- stat
	}()

	return reqsStatCh, statCh
}

func printStat(stat *Stat) {
	fmt.Printf(
		"Results:\n"+
			// Counters
			"\tCounters:\n"+
			"\t\tTotal: %d\n"+
			"\t\tSuccess: %d\n"+
			"\t\tFaileds: %d\n"+
			"\t\tUnfinished: %d\n"+
			// Durations
			"\tDurations:\n"+
			"\t\tTest: %s\n"+
			"\t\tRequests:\n"+
			"\t\t\tTotal: %s\n"+
			"\t\t\tSuccess: %s\n"+
			"\t\t\tFaileds: %s\n"+
			// Durations[AVG]
			"\tDurations[AVG]:\n"+
			"\t\tTotal: %s\n"+
			"\t\tSuccess: %s\n"+
			"\t\tFaileds: %s\n"+
			// Test
			"\tTest:\n"+
			"\t\tStart: %s\n"+
			"\t\tStop: %s\n",
		// Counters
		stat.ReqsNum,
		stat.SuccessReqs,
		stat.FailedReqs,
		stat.UnfinishedReqs,
		// Durations
		stat.Duration,
		stat.ReqsDur,
		stat.SuccessReqsDur,
		stat.FailedReqsDur,
		// Durations[AVG]
		stat.AvgReqsDur,
		stat.AvgSuccessReqsDur,
		stat.AvgFailedReqsDur,
		// Test
		stat.From,
		stat.To,
	)
}

func main() {
	osSigsCh := make(chan os.Signal, 1)
	defer close(osSigsCh)
	signal.Notify(osSigsCh, os.Interrupt)
	wg := &sync.WaitGroup{}
	start := time.Now()

	reqsStatCh, statCh := handleReqs()
	dCh, ctx, cancel, pwg := provide()
	cwg := consume(dCh, reqsStatCh, ctx)

	defer func() {
		printStat(<-statCh)
	}()
	defer func() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stopProvide(dCh, cancel, pwg)
			stopConsume(cwg, reqsStatCh)
		}()
		wg.Wait()
	}()
	for {
		select {
		case <-osSigsCh:
			return
		default:
			if time.Since(start) > time.Second*600 {
				return
			}
			runtime.Gosched()
		}
	}
}
