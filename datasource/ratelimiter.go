package datasource

import (
	"time"
)

type RateLimiter struct {
	ticker *time.Ticker
}

func NewRateLimiter(reqPerSec int) *RateLimiter {
	interval := time.Second / time.Duration(reqPerSec)
	return &RateLimiter{ticker: time.NewTicker(interval)}
}

func (r *RateLimiter) Wait() {
	<-r.ticker.C
}

func (r *RateLimiter) Stop() {
	r.ticker.Stop()
}
