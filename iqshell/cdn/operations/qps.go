package operations

import "time"

var prefetchTimeTicker *time.Ticker

func createQpsLimitIfNeeded(limit int) {
	if limit < 1 {
		return
	}

	if prefetchTimeTicker == nil {
		prefetchTimeTicker = time.NewTicker(time.Second / time.Duration(limit))
	}
}

func waiterIfNeeded() {
	if prefetchTimeTicker != nil {
		<-prefetchTimeTicker.C
	}
}

