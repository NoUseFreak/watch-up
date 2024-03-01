package watchup

import "time"

type AppSettings struct {
	LoopDuration time.Duration
	URL          string
}
