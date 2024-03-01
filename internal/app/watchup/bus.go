package watchup

import "os"

type ChanBus struct {
	ChangeResults chan WatchResult
	WatchResults  chan WatchResult
	CodeStats     chan WatchCodeStats
	Shutdown      chan os.Signal
	RequestStats  chan RequestStats
}
