package middleware

import (
	"sync"
	"time"
)

const (
	threshold = 20
	ttl       = time.Minute * 1
)

type (
	TokenInfo struct {
		RemainingTokens int32
		Expiry          time.Time
	}

	TokenBucket struct {
		threshold     int32
		ttl           time.Duration
		tokens        map[string]TokenInfo
		mx            sync.RWMutex
		cleanupTicker *time.Ticker
		stopCleanup   chan bool
	}

	Metadata struct {
		allowed   bool
		remaining int32
		expiry    time.Time
		limit     int32
		used      int32
	}
)

func NewTokenBucket(threshold int32, ttl time.Duration) *TokenBucket {
	tb := &TokenBucket{
		threshold:     threshold,
		ttl:           ttl,
		tokens:        make(map[string]TokenInfo),
		cleanupTicker: time.NewTicker(ttl),
		stopCleanup:   make(chan bool),
	}

	go tb.cleanupExpiredTokens()

	return tb
}

func DefaultTokenBucket() *TokenBucket {
	return NewTokenBucket(threshold, ttl)
}

func (t *TokenBucket) Allow(key string) bool {
	t.mx.Lock()
	defer t.mx.Unlock()

	now := time.Now()
	info, ok := t.tokens[key]
	if ok && now.After(info.Expiry) {
		info = TokenInfo{
			RemainingTokens: t.threshold - 1,
			Expiry:          now.Add(t.ttl),
		}

		t.tokens[key] = info

		return true
	}

	if !ok {
		info = TokenInfo{
			RemainingTokens: t.threshold - 1,
			Expiry:          now.Add(t.ttl),
		}

		t.tokens[key] = info

		return true
	}

	if info.RemainingTokens > 0 {
		info.RemainingTokens--
		t.tokens[key] = info

		return true
	}

	return false
}

func (t *TokenBucket) cleanupExpiredTokens() {
	for {
		select {
		case <-t.cleanupTicker.C:
			t.mx.Lock()
			now := time.Now()
			for key, info := range t.tokens {
				if now.After(info.Expiry) {
					delete(t.tokens, key)
				}
			}
			t.mx.Unlock()
		case <-t.stopCleanup:
			t.cleanupTicker.Stop()

			return
		}
	}
}

func (t *TokenBucket) Stop() {
	close(t.stopCleanup)
}

func (t *TokenBucket) GetTokenInfo(
	key string,
) (remaining int32, expiry time.Time, limit int32) {
	t.mx.RLock()
	defer t.mx.RUnlock()

	now := time.Now()
	info, ok := t.tokens[key]

	if !ok || now.After(info.Expiry) {
		return t.threshold, now.Add(t.ttl), t.threshold
	}

	return info.RemainingTokens, info.Expiry, t.threshold
}

func (t *TokenBucket) AllowAndGetInfo(key string) Metadata {
	t.mx.Lock()
	defer t.mx.Unlock()

	now := time.Now()
	info, ok := t.tokens[key]

	if ok && now.After(info.Expiry) {
		info = TokenInfo{
			RemainingTokens: t.threshold - 1,
			Expiry:          now.Add(t.ttl),
		}
		t.tokens[key] = info

		return Metadata{
			allowed:   true,
			remaining: info.RemainingTokens,
			expiry:    info.Expiry,
			limit:     t.threshold,
			used:      1,
		}
	}

	if !ok {
		info = TokenInfo{
			RemainingTokens: t.threshold - 1,
			Expiry:          now.Add(t.ttl),
		}
		t.tokens[key] = info

		return Metadata{
			allowed:   true,
			remaining: info.RemainingTokens,
			expiry:    info.Expiry,
			limit:     t.threshold,
			used:      1,
		}
	}

	if info.RemainingTokens > 0 {
		info.RemainingTokens--
		t.tokens[key] = info
		usedTokens := t.threshold - info.RemainingTokens

		return Metadata{
			allowed:   true,
			remaining: info.RemainingTokens,
			expiry:    info.Expiry,
			limit:     t.threshold,
			used:      usedTokens,
		}
	}

	usedTokens := t.threshold - info.RemainingTokens

	return Metadata{
		allowed:   false,
		remaining: 0,
		expiry:    info.Expiry,
		limit:     t.threshold,
		used:      usedTokens,
	}
}
