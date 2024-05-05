package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type RateLimiter2 struct {
	Rate     int         `json:"rate,omitempty"`
	Bucket   int         `json:"bucket,omitempty"`
	Capacity int         `json:"capacity,omitempty"`
	Ticker   time.Ticker `json:"ticket,omitempty"`
}

func NewRateLimiter2(rate int, capacity int) *RateLimiter2 {
	rl := &RateLimiter2{
		Rate:     rate,
		Bucket:   capacity,
		Capacity: capacity,
		Ticker:   *time.NewTicker(time.Second / 10),
	}

	go rl.Refill()
	return rl
}

func (rl *RateLimiter2) Refill() {
	for range rl.Ticker.C {
		rl.Bucket += rl.Rate
		if rl.Bucket > rl.Capacity {
			rl.Bucket = rl.Capacity
		}
	}
}

func (rl *RateLimiter2) Limit(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if rl.Bucket > 0 {
			rl.Bucket--
			next(w, r)
		} else {
			message := Message{
				Status: "Failed",
				Body:   "Rate limited2 reqeust",
			}

			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&message)
			return
		}
	}

}
