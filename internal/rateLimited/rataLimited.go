package rateLimited

import (
	"context"
	"time"
)

type TokenBucket struct {
	tokens chan struct{}
	ticker *time.Ticker
}

func (t *TokenBucket) Allow() bool {
	select {
	case <-t.tokens:
		return true
	default:
		return false
	}
}

func NewTokenBucket(ctx context.Context, rate time.Duration, capacity int) *TokenBucket {
	tkn := &TokenBucket{
		tokens: make(chan struct{}, capacity),
		ticker: time.NewTicker(rate),
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
