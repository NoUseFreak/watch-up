package watchup

import "time"

// AppSettings is a struct that contains the settings for the application.
type AppSettings struct {
	LoopDuration time.Duration
	URL          string
}
