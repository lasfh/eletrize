package command

import (
	"sync"
	"time"
)

func debounce(delay time.Duration, fn func(string)) func(string) {
	var mu sync.Mutex
	var timer *time.Timer

	return func(event string) {
		mu.Lock()
		defer mu.Unlock()

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(delay, func() { fn(event) })
	}
}
