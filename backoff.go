package cast

import (
	"math"
	"math/rand"
	"time"
)

type backoffStrategy interface {
	backoff(retry int) time.Duration
}

type linearBackoffStrategy struct {
	slope time.Duration
}

func (strat linearBackoffStrategy) backoff(retry int) time.Duration {
	return time.Duration(retry) * strat.slope
}

type constantBackOffStrategy struct {
	interval time.Duration
}

func (strat constantBackOffStrategy) backoff(retry int) time.Duration {
	return strat.interval
}

type exponentialBackoff struct {
	base time.Duration
	cap  time.Duration
}

func (backoff exponentialBackoff) expo(retry int) float64 {
	b := float64(backoff.base)
	r := float64(retry)
	return math.Min(b, math.Exp2(r)*b)
}

type exponentialBackoffStrategy struct {
	exponentialBackoff
}

func (strat exponentialBackoffStrategy) backoff(retry int) time.Duration {
	return time.Duration(strat.expo(retry))
}

type exponentialBackoffEqualJitterStrategy struct {
	exponentialBackoff
}

func (strat exponentialBackoffEqualJitterStrategy) backoff(retry int) time.Duration {
	v := strat.expo(retry)
	return time.Duration(v/2.0 + rand.Float64()*v/2.0)
}

type exponentialBackoffFullJitterStrategy struct {
	exponentialBackoff
}

func (strat exponentialBackoffFullJitterStrategy) backoff(retry int) time.Duration {
	v := strat.expo(retry)
	return time.Duration(rand.Float64() * v)
}

type exponentialBackoffDecorrelatedJitterStrategy struct {
	exponentialBackoff
	sleep time.Duration
}

func (strat exponentialBackoffDecorrelatedJitterStrategy) backoff(retry int) time.Duration {
	return time.Duration(math.Min(float64(strat.cap), float64(strat.base)+rand.Float64()*(float64(strat.sleep)*3-float64(strat.base))))
}
