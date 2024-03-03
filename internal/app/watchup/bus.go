package watchup

import "os"

// ChanBus is a struct that contains channels for communication between the
// different parts of the application.
type ChanBus struct {
	ChangeResults chan WatchResult
	WatchResults  chan WatchResult
	CodeStats     chan WatchCodeStats
	Shutdown      chan os.Signal
	RequestStats  chan RequestStats
}
