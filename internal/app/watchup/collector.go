package watchup

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Collector is a struct that contains the state of the collector.
type Collector struct {
	codeLog map[string]lastCode
	sum     map[int]time.Duration
	count   int64
}

type lastCode struct {
	code  int
	time  time.Time
	count int64
}

// WatchCodeStats is a map of status codes to the total time spent on them.
type WatchCodeStats map[int]time.Duration

// Start starts the collector and listens for results on the bus.
func (c *Collector) Start(bus ChanBus) {

	c.codeLog = make(map[string]lastCode)
	c.sum = make(WatchCodeStats)

	requestSetLen := 100
	var requestSet = make([]WatchResult, 0)

	t := time.Now()
	for result := range bus.WatchResults {
		c.sum[result.StatusCode] += time.Since(t)
		bus.CodeStats <- c.sum
		t = time.Now()
		c.count++

		requestSet = append(requestSet, result)
		if len(requestSet) > requestSetLen {
			requestSet = requestSet[1:]
		}

		bus.RequestStats <- c.getRequestStats(requestSet)

		last, ok := c.codeLog[result.Url]
		if !ok {
			// first time we've seen this URL
			c.codeLog[result.Url] = lastCode{code: result.StatusCode, time: t, count: c.count}
			result.Count = c.count
			bus.ChangeResults <- result
			continue
		}

		if last.code != result.StatusCode {
			c.codeLog[result.Url] = lastCode{code: result.StatusCode, time: t, count: c.count}
			result.DeltaTime = time.Since(last.time)
			result.Count = c.count - last.count
			bus.ChangeResults <- result
		}
	}
}

// Stop stops the collector and prints the total time spent on each status code.
func (c *Collector) Stop() {
	for code, duration := range c.sum {
		fmt.Printf("Total time for %s: %s\n", FormatStatusCode(code), duration)
	}
}

// RequestStats is a struct that contains statistics about the requests.
type RequestStats struct {
	TotalRequests int64
	TotalTime     time.Duration
	AvgTime       time.Duration
	Jitter        time.Duration
}

func (c *Collector) getRequestStats(requests []WatchResult) RequestStats {
	var totalTime int64
	var responseTimes []int64
	for _, result := range requests {
		totalTime += result.Total.Microseconds()
		responseTimes = append(responseTimes, result.Total.Microseconds())
	}

	avgTime := calculateAverage(responseTimes)
	jitter := calculateJitter(responseTimes, avgTime)
	logrus.Infof("Total Avg time: %v, Jitter: %v", time.Duration(avgTime*int64(time.Microsecond)), time.Duration(jitter*int64(time.Microsecond)))

	return RequestStats{
		TotalRequests: c.count,
		AvgTime:       time.Duration(avgTime * int64(time.Microsecond)),
		Jitter:        time.Duration(jitter * int64(time.Microsecond)),
	}
}
