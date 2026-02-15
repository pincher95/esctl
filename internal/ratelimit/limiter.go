package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Limiter controls the rate of operations
type Limiter struct {
	ticker   *time.Ticker
	interval time.Duration
	mu       sync.Mutex
	stopped  bool
}

// New creates a new rate limiter that allows requestsPerSecond operations
func New(requestsPerSecond int) *Limiter {
	if requestsPerSecond <= 0 {
		requestsPerSecond = 1
	}

	interval := time.Second / time.Duration(requestsPerSecond)
	return &Limiter{
		ticker:   time.NewTicker(interval),
		interval: interval,
	}
}

// NewWithInterval creates a new rate limiter with a specific interval between operations
func NewWithInterval(interval time.Duration) *Limiter {
	if interval <= 0 {
		interval = time.Second
	}

	return &Limiter{
		ticker:   time.NewTicker(interval),
		interval: interval,
	}
}

// Wait blocks until the next operation is allowed or context is cancelled
func (l *Limiter) Wait(ctx context.Context) error {
	l.mu.Lock()
	stopped := l.stopped
	l.mu.Unlock()

	if stopped {
		return context.Canceled
	}

	select {
	case <-l.ticker.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop stops the rate limiter and releases resources
func (l *Limiter) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.stopped {
		return
	}

	l.stopped = true
	if l.ticker != nil {
		l.ticker.Stop()
	}
}

// TokenBucket implements a token bucket rate limiter
// This allows for bursts while maintaining an average rate
type TokenBucket struct {
	capacity    int
	tokens      int
	refillRate  int // tokens per second
	lastRefill  time.Time
	mu          sync.Mutex
	refillTimer *time.Ticker
	stopped     bool
}

// NewTokenBucket creates a new token bucket rate limiter
// capacity: maximum number of tokens (burst size)
// refillRate: tokens added per second
func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	if capacity <= 0 {
		capacity = 10
	}
	if refillRate <= 0 {
		refillRate = 1
	}

	tb := &TokenBucket{
		capacity:   capacity,
		tokens:     capacity, // Start with full bucket
		refillRate: refillRate,
		lastRefill: time.Now(),
	}

	// Start refill timer
	tb.refillTimer = time.NewTicker(time.Second)
	go tb.refillLoop()

	return tb
}

// refillLoop continuously refills tokens
func (tb *TokenBucket) refillLoop() {
	for range tb.refillTimer.C {
		tb.mu.Lock()
		if tb.stopped {
			tb.mu.Unlock()
			return
		}

		// Add tokens based on time elapsed
		tokensToAdd := tb.refillRate
		tb.tokens += tokensToAdd

		// Cap at capacity
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}

		tb.lastRefill = time.Now()
		tb.mu.Unlock()
	}
}

// TryTake attempts to take n tokens, returns true if successful
func (tb *TokenBucket) TryTake(n int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.stopped {
		return false
	}

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}

	return false
}

// Take blocks until n tokens are available or context is cancelled
func (tb *TokenBucket) Take(ctx context.Context, n int) error {
	if n <= 0 {
		return nil
	}

	for {
		if tb.TryTake(n) {
			return nil
		}

		// Wait a bit before retrying
		select {
		case <-time.After(50 * time.Millisecond):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Available returns the number of available tokens
func (tb *TokenBucket) Available() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.tokens
}

// Stop stops the token bucket and releases resources
func (tb *TokenBucket) Stop() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.stopped {
		return
	}

	tb.stopped = true
	if tb.refillTimer != nil {
		tb.refillTimer.Stop()
	}
}

// AdaptiveLimiter adjusts rate based on success/failure
type AdaptiveLimiter struct {
	minRate     int
	maxRate     int
	currentRate int
	limiter     *Limiter
	mu          sync.Mutex
}

// NewAdaptiveLimiter creates an adaptive rate limiter
func NewAdaptiveLimiter(minRate, maxRate, startRate int) *AdaptiveLimiter {
	if minRate <= 0 {
		minRate = 1
	}
	if maxRate <= minRate {
		maxRate = minRate * 10
	}
	if startRate <= 0 || startRate < minRate || startRate > maxRate {
		startRate = minRate
	}

	return &AdaptiveLimiter{
		minRate:     minRate,
		maxRate:     maxRate,
		currentRate: startRate,
		limiter:     New(startRate),
	}
}

// Wait blocks until the next operation is allowed
func (al *AdaptiveLimiter) Wait(ctx context.Context) error {
	return al.limiter.Wait(ctx)
}

// OnSuccess indicates a successful operation, may increase rate
func (al *AdaptiveLimiter) OnSuccess() {
	al.mu.Lock()
	defer al.mu.Unlock()

	// Increase rate by 10% up to max
	newRate := min(int(float64(al.currentRate)*1.1), al.maxRate)

	if newRate != al.currentRate {
		al.updateRate(newRate)
	}
}

// OnFailure indicates a failed operation, decreases rate
func (al *AdaptiveLimiter) OnFailure() {
	al.mu.Lock()
	defer al.mu.Unlock()

	// Decrease rate by 50% down to min
	newRate := max(al.currentRate/2, al.minRate)

	if newRate != al.currentRate {
		al.updateRate(newRate)
	}
}

// updateRate changes the current rate (must be called with lock held)
func (al *AdaptiveLimiter) updateRate(newRate int) {
	al.currentRate = newRate

	// Replace limiter with new rate
	if al.limiter != nil {
		al.limiter.Stop()
	}
	al.limiter = New(newRate)
}

// CurrentRate returns the current rate
func (al *AdaptiveLimiter) CurrentRate() int {
	al.mu.Lock()
	defer al.mu.Unlock()
	return al.currentRate
}

// Stop stops the adaptive limiter
func (al *AdaptiveLimiter) Stop() {
	al.mu.Lock()
	defer al.mu.Unlock()

	if al.limiter != nil {
		al.limiter.Stop()
	}
}
