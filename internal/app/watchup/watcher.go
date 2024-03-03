package watchup

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tcnksm/go-httpstat"
)

// Watcher is a struct that contains the state of the watcher.
type Watcher struct {
	url    string
	loop   time.Duration
	client *http.Client

	stop chan struct{}
}

// NewWatcher returns a new Watcher.
func NewWatcher(url string, loop time.Duration) *Watcher {
	return &Watcher{
		url:  url,
		loop: loop,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Start starts the watcher and sends results to the bus.
func (w *Watcher) Start(bus ChanBus) {
	w.stop = make(chan struct{})

	ticker := time.NewTicker(w.loop)
	bus.WatchResults <- w.doRequest()
	for {
		select {
		case <-ticker.C:
			go func() {
				bus.WatchResults <- w.doRequest()
			}()
		case <-w.stop:
			ticker.Stop()
			return
		}
	}
}

// Stop stops the watcher.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) createHttpRequest() *http.Request {
	req, err := http.NewRequest("GET", w.url, nil)
	if err != nil {
		log.Fatal(err)
	}

	return req
}

func (w *Watcher) doRequest() WatchResult {
	req := w.createHttpRequest()

	start := time.Now()
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error: %v", err)
		result.End(time.Now())

		code := codeFromError(err)

		return WatchResult{
			Url:        w.url,
			HttpStats:  result,
			StatusCode: code,
			Total:      time.Since(start),
			Time:       time.Now(),
		}
	}

	if _, err := io.Copy(io.Discard, res.Body); err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	now := time.Now()
	result.End(now)

	logrus.Debugf("result: %v, %v", now.Sub(start), res.StatusCode)
	return WatchResult{
		Url:        w.url,
		HttpStats:  result,
		StatusCode: res.StatusCode,
		Total:      now.Sub(start),
		Time:       now,
	}
}

// WatchResult is a struct that contains the result of a request.
type WatchResult struct {
	Url        string
	HttpStats  httpstat.Result
	StatusCode int
	Total      time.Duration
	Time       time.Time
	DeltaTime  time.Duration
	Count      int64
}
