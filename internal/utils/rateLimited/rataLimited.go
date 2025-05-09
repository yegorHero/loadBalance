package rateLimited

import (
	"context"
	"time"
)

type TokenBucket struct {
	tokens chan struct{}
	ticker *time.Ticker
	cancel context.CancelFunc
}

func (t *TokenBucket) Allow() bool {
	select {
	case <-t.tokens:
		return true
	default:
		return false
	}
}

func (t *TokenBucket) Stop() {
	t.cancel()
}

func NewTokenBucket(rate time.Duration, capacity int) *TokenBucket {
	ctx, cancel := context.WithCancel(context.Background())

	tkn := &TokenBucket{
		tokens: make(chan struct{}, capacity),
		ticker: time.NewTicker(rate),
		cancel: cancel,
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
			case <-ctx.Done():
				tkn.ticker.Stop()
				return
			}
		}
	}()

	return tkn
}
