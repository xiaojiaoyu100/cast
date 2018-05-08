package cast

import (
	"testing"
	"time"
)

func TestCast_WithConstantBackoffStrategy(t *testing.T) {
	c := New()
	interval := time.Second
	c.WithConstantBackoffStrategy(interval)
	for i := 0; i < 100; i++ {
		t.Log(c.strat.backoff(i))
		if c.strat.backoff(i) != interval {
			t.Fatal("unexpected")
		}
	}
}

func TestCast_WithLinearBackoffStrategy(t *testing.T) {
	c := New()
	slope := 30 * time.Millisecond
	c.WithLinearBackoffStrategy(slope)
	for i := 0; i < 100; i++ {
		t.Log(c.strat.backoff(i))
		if c.strat.backoff(i) != (time.Duration(i) * slope) {
			t.Fatal("unexpected")
		}
	}
}

func TestCast_WithExponentialBackoffStrategy(t *testing.T) {
	c := New()
	c.WithExponentialBackoffStrategy(10*time.Millisecond, 20*time.Millisecond)
	for i := 0; i < 100; i++ {
		t.Log(c.strat.backoff(i))
	}
}

func TestCast_WithExponentialBackoffEqualJitterStrategy(t *testing.T) {
	c := New()
	c.WithExponentialBackoffEqualJitterStrategy(10*time.Millisecond, 20*time.Millisecond)
	for i := 0; i < 100; i++ {
		t.Log(c.strat.backoff(i))
	}
}

func TestCast_WithExponentialBackoffFullJitterStrategy(t *testing.T) {
	c := New()
	c.WithExponentialBackoffFullJitterStrategy(10*time.Millisecond, 20*time.Millisecond)
	for i := 0; i < 100; i++ {
		t.Log(c.strat.backoff(i))
	}
}

func TestCast_WithExponentialBackoffDecorrelatedJitterStrategy(t *testing.T) {
	c := New()
	c.WithExponentialBackoffDecorrelatedJitterStrategy(100*time.Millisecond, 200*time.Millisecond)
	for i := 0; i < 100; i++ {
		t.Log(c.strat.backoff(i))
	}
}
