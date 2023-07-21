package service

import (
	"fmt"
	"github.com/Borislavv/ddos/internal/shared/infrastructure/network/safehttp"
	"github.com/Borislavv/ddos/internal/tester/domain/model"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Meter struct {
	mu           *sync.Mutex
	cfg          *model.Settings
	measurements *model.Measurements
}

func NewMeter(cfg *model.Settings) *Meter {
	return &Meter{
		mu:           &sync.Mutex{},
		cfg:          cfg,
		measurements: &model.Measurements{},
	}
}

func (m *Meter) Start() {
	// todo must be solved
	if err := m.measurements.Timestamp.SetFrom(); err != nil {
		log.Fatalln(err)
	}
}

func (m *Meter) Stop() {
	// todo must be solved
	if err := m.measurements.Timestamp.SetTo(); err != nil {
		log.Fatalln(err)
	}
	if err := m.measurements.Timestamp.SetDuration(); err != nil {
		log.Fatalln(err)
	}
}

func (m *Meter) CommitReq(req *safehttp.Req) {
	m.commitReqCounters(req)
	m.commitReqDurations(req)
}

func (m *Meter) Summary() {
	// todo must be solved (and moved to Stop() method )
	m.commitReqsAvgDurations()

	totalDuration := m.measurements.Timestamp.GetDuration()

	fmt.Printf(
		"\n\nSummary\n"+
			"\tRequest:\n"+
			"\t\tSuccess: %d\n"+
			"\t\tFailed: %d\n"+
			"\t\tTotal: %d\n"+
			"\tDuration:\n"+
			"\t\tAVG[total]: %s\n"+
			"\t\tAVG[success]: %s\n"+
			"\t\tMIN[success]: %s\n"+
			"\t\tMAX[success]: %s\n"+
			"\t\tTotal: %s\n"+
			"\tWorkers:\n"+
			"\t\tConsumers: %d\n"+
			"\t\tProviders: %d\n"+
			"\t\tDisplayers: %d\n",
		m.measurements.SuccessReqs,
		m.measurements.FailedReqs,
		m.measurements.TotalReqs,
		time.Duration(m.measurements.AvgTotalReqsDur),
		time.Duration(m.measurements.AvgSuccessReqsDur),
		time.Duration(m.measurements.MinSuccessReqDur),
		time.Duration(m.measurements.MaxSuccessReqDur),
		totalDuration,
		m.cfg.Workers.Total,
		m.cfg.Providers.Total,
		m.cfg.Displayers.Total,
	)
}

func (m *Meter) commitReqCounters(req *safehttp.Req) {
	atomic.AddInt64(&m.measurements.TotalReqs, 1)
	if req.IsFailed {
		atomic.AddInt64(&m.measurements.FailedReqs, 1)
	} else {
		atomic.AddInt64(&m.measurements.SuccessReqs, 1)
	}
}

func (m *Meter) commitReqDurations(req *safehttp.Req) {
	defer m.mu.Unlock()
	m.mu.Lock()

	dur := req.Timestamp.GetDuration()
	atomic.AddInt64(&m.measurements.TotalReqsDur, dur.Nanoseconds())

	if req.IsFailed {
		atomic.AddInt64(&m.measurements.FailedReqsDur, dur.Nanoseconds())
	} else {
		atomic.AddInt64(&m.measurements.SuccessReqsDur, dur.Nanoseconds())
	}

	if m.measurements.MinSuccessReqDur == 0 {
		m.measurements.MinSuccessReqDur = dur.Nanoseconds()
	} else if m.measurements.MinSuccessReqDur > dur.Nanoseconds() {
		m.measurements.MinSuccessReqDur = dur.Nanoseconds()
	}

	if m.measurements.MaxSuccessReqDur == 0 {
		m.measurements.MaxSuccessReqDur = dur.Nanoseconds()
	} else if m.measurements.MaxSuccessReqDur < dur.Nanoseconds() {
		m.measurements.MaxSuccessReqDur = dur.Nanoseconds()
	}
}

func (m *Meter) commitReqsAvgDurations() {
	defer m.mu.Unlock()
	m.mu.Lock()

	if m.measurements.TotalReqs > 0 {
		m.measurements.AvgTotalReqsDur = float64(m.measurements.TotalReqsDur) / float64(m.measurements.TotalReqs)
	} else {
		m.measurements.AvgTotalReqsDur = 0
	}

	if m.measurements.SuccessReqs > 0 {
		m.measurements.AvgSuccessReqsDur = float64(m.measurements.SuccessReqsDur) / float64(m.measurements.SuccessReqs)
	} else {
		m.measurements.AvgSuccessReqsDur = 0
	}

	if m.measurements.FailedReqs > 0 {
		m.measurements.AvgFailedReqsDur = float64(m.measurements.FailedReqsDur) / float64(m.measurements.FailedReqs)
	} else {
		m.measurements.AvgFailedReqsDur = 0
	}
}
