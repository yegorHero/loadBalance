package rateLimited

import "time"

type TokenBucket struct {
	tokens chan struct{}
	ticker *time.Ticker
	stop   chan struct{}
}

func NewTokenBucket(rate time.Duration, capacity int) *TokenBucket {
	tkn := &TokenBucket{
		tokens: make(chan struct{}, capacity),
		ticker: time.NewTicker(rate),
		stop:   make(chan struct{}),
	}

	for i := 0; i < capacity; i++ {
		tkn.tokens <- struct{}{}
	}

	go func() {
		for {
			select {
			case <-tkn.ticker.C:
				select {
				case tkn.tokens <- struct{}{}:
				default:
				}
			case <-tkn.stop:
				tkn.ticker.Stop()
				return
			}
		}
	}()

	return tkn
}
