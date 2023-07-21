package model

import "github.com/Borislavv/ddos/internal/tester/domain/vo"

type Measurements struct {
	Timestamp vo.Timestamp

	// reqs timestamps
	AvgSuccessReqsDur float64
	SuccessReqsDur    int64
	AvgFailedReqsDur  float64
	FailedReqsDur     int64
	AvgTotalReqsDur   float64
	TotalReqsDur      int64
	MinSuccessReqDur  int64
	MaxSuccessReqDur  int64

	// counters
	SuccessReqs int64
	FailedReqs  int64
	TotalReqs   int64
	PerSecond   float64
	// threshold percentages
}
