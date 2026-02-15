package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	limiter := New(10)
	if limiter == nil {
		t.Error("New() should return a non-nil limiter")
	}

	limiter.Stop()
}

func TestNewWithZeroRate(t *testing.T) {
	limiter := New(0)
	if limiter == nil {
		t.Error("New() should return a non-nil limiter even with 0 rate")
	}

	limiter.Stop()
}

func TestNewWithInterval(t *testing.T) {
	limiter := NewWithInterval(100 * time.Millisecond)
	if limiter == nil {
		t.Error("NewWithInterval() should return a non-nil limiter")
	}

	if limiter.interval != 100*time.Millisecond {
		t.Errorf("Expected interval 100ms, got %v", limiter.interval)
	}

	limiter.Stop()
}

func TestLimiterWait(t *testing.T) {
	limiter := New(10) // 10 requests per second = 100ms per request
	defer limiter.Stop()

	ctx := context.Background()
	start := time.Now()

	// First wait should be immediate
	err := limiter.Wait(ctx)
	if err != nil {
		t.Fatalf("Wait() error = %v", err)
	}

	// Second wait should take ~100ms
	err = limiter.Wait(ctx)
	if err != nil {
		t.Fatalf("Wait() error = %v", err)
	}

	elapsed := time.Since(start)
	if elapsed < 90*time.Millisecond {
		t.Errorf("Expected at least 90ms elapsed, got %v", elapsed)
	}
}

func TestLimiterWaitWithCancelledContext(t *testing.T) {
	limiter := New(1) // 1 request per second
	defer limiter.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := limiter.Wait(ctx)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

func TestLimiterStop(t *testing.T) {
	limiter := New(10)
	limiter.Stop()
	limiter.Stop() // Should be idempotent

	if !limiter.stopped {
		t.Error("Expected limiter to be stopped")
	}
}

func TestNewTokenBucket(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	defer tb.Stop()

	if tb == nil {
		t.Error("NewTokenBucket() should return a non-nil token bucket")
	}

	if tb.capacity != 10 {
		t.Errorf("Expected capacity 10, got %d", tb.capacity)
	}

	if tb.refillRate != 5 {
		t.Errorf("Expected refill rate 5, got %d", tb.refillRate)
	}

	// Should start with full bucket
	if tb.tokens != 10 {
		t.Errorf("Expected initial tokens 10, got %d", tb.tokens)
	}
}

func TestTokenBucketTryTake(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	defer tb.Stop()

	// Should be able to take 5 tokens
	if !tb.TryTake(5) {
		t.Error("Should be able to take 5 tokens")
	}

	// Should have 5 tokens left
	if tb.tokens != 5 {
		t.Errorf("Expected 5 tokens remaining, got %d", tb.tokens)
	}

	// Should not be able to take 10 tokens (only 5 remaining)
	if tb.TryTake(10) {
		t.Error("Should not be able to take 10 tokens")
	}
}

func TestTokenBucketTake(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	defer tb.Stop()

	ctx := context.Background()

	// Should be able to take 10 tokens immediately
	err := tb.Take(ctx, 10)
	if err != nil {
		t.Fatalf("Take() error = %v", err)
	}

	// Now bucket is empty, should have 0 tokens
	if tb.tokens != 0 {
		t.Errorf("Expected 0 tokens remaining, got %d", tb.tokens)
	}
}

func TestTokenBucketRefill(t *testing.T) {
	tb := NewTokenBucket(10, 10) // Refill 10 per second
	defer tb.Stop()

	// Empty the bucket
	tb.TryTake(10)

	// Wait for refill (1 second + some buffer)
	time.Sleep(1100 * time.Millisecond)

	// Should have refilled
	available := tb.Available()
	if available < 8 { // Allow some margin
		t.Errorf("Expected at least 8 tokens after refill, got %d", available)
	}
}

func TestTokenBucketAvailable(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	defer tb.Stop()

	if tb.Available() != 10 {
		t.Errorf("Expected 10 available tokens, got %d", tb.Available())
	}

	tb.TryTake(3)

	if tb.Available() != 7 {
		t.Errorf("Expected 7 available tokens, got %d", tb.Available())
	}
}

func TestTokenBucketStop(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	tb.Stop()
	tb.Stop() // Should be idempotent

	if !tb.stopped {
		t.Error("Expected token bucket to be stopped")
	}

	// Operations after stop should fail gracefully
	if tb.TryTake(1) {
		t.Error("TryTake should fail after stop")
	}
}

func TestNewAdaptiveLimiter(t *testing.T) {
	al := NewAdaptiveLimiter(1, 100, 10)
	defer al.Stop()

	if al == nil {
		t.Error("NewAdaptiveLimiter() should return a non-nil limiter")
	}

	if al.minRate != 1 {
		t.Errorf("Expected minRate 1, got %d", al.minRate)
	}

	if al.maxRate != 100 {
		t.Errorf("Expected maxRate 100, got %d", al.maxRate)
	}

	if al.currentRate != 10 {
		t.Errorf("Expected currentRate 10, got %d", al.currentRate)
	}
}

func TestAdaptiveLimiterOnSuccess(t *testing.T) {
	al := NewAdaptiveLimiter(1, 100, 10)
	defer al.Stop()

	initialRate := al.CurrentRate()

	// Simulate success
	al.OnSuccess()

	newRate := al.CurrentRate()
	if newRate <= initialRate {
		t.Errorf("Expected rate to increase after success, got %d -> %d", initialRate, newRate)
	}
}

func TestAdaptiveLimiterOnFailure(t *testing.T) {
	al := NewAdaptiveLimiter(1, 100, 10)
	defer al.Stop()

	initialRate := al.CurrentRate()

	// Simulate failure
	al.OnFailure()

	newRate := al.CurrentRate()
	if newRate >= initialRate {
		t.Errorf("Expected rate to decrease after failure, got %d -> %d", initialRate, newRate)
	}
}

func TestAdaptiveLimiterRateBounds(t *testing.T) {
	al := NewAdaptiveLimiter(5, 20, 10)
	defer al.Stop()

	// Try to increase rate beyond max
	for range 20 {
		al.OnSuccess()
	}

	rate := al.CurrentRate()
	if rate > 20 {
		t.Errorf("Rate should not exceed maxRate, got %d", rate)
	}

	// Try to decrease rate below min
	for range 20 {
		al.OnFailure()
	}

	rate = al.CurrentRate()
	if rate < 5 {
		t.Errorf("Rate should not go below minRate, got %d", rate)
	}
}

func TestAdaptiveLimiterWait(t *testing.T) {
	al := NewAdaptiveLimiter(10, 100, 10)
	defer al.Stop()

	ctx := context.Background()

	err := al.Wait(ctx)
	if err != nil {
		t.Fatalf("Wait() error = %v", err)
	}
}

func TestLimiterMultipleWaits(t *testing.T) {
	limiter := New(100) // 100 requests per second
	defer limiter.Stop()

	ctx := context.Background()

	// Should be able to wait multiple times
	for i := range 5 {
		err := limiter.Wait(ctx)
		if err != nil {
			t.Fatalf("Wait() iteration %d error = %v", i, err)
		}
	}
}

func TestTokenBucketTakeWithContext(t *testing.T) {
	tb := NewTokenBucket(10, 1)
	defer tb.Stop()

	// Empty the bucket
	tb.TryTake(10)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Try to take more tokens than available with short timeout
	err := tb.Take(ctx, 5)
	if err == nil {
		t.Error("Expected timeout error when tokens not available")
	}
}
