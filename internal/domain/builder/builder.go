package builder

import (
	"net/http"
	"time"
)

type TaskBuilderInterface interface {
	BuildFromRaw(
		req *http.Request,
		workersNum int,
		requestsNum int64,
		duration time.Duration,
		reqDurThresholds map[int64]string,
	)
}
