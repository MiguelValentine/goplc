package lib

import "time"

func Cron(d time.Duration, f func()) {
	go func() {
		f()
		for range time.Tick(d) {
			f()
		}
	}()
}
