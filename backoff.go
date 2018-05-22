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
	c := float64(backoff.cap)
	b := float64(backoff.base)
	r := float64(retry)
	return math.Min(c, math.Exp2(r)*b)
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
	u := uniform(0, v/2.0)
	return time.Duration(v/2.0 + u)
}

type exponentialBackoffFullJitterStrategy struct {
	exponentialBackoff
}

func (strat exponentialBackoffFullJitterStrategy) backoff(retry int) time.Duration {
	v := strat.expo(retry)
	u := uniform(0, v)
	return time.Duration(u)
}

type exponentialBackoffDecorrelatedJitterStrategy struct {
	exponentialBackoff
	sleep time.Duration
}

// uniform returns number in [min, max)
func uniform(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (strat exponentialBackoffDecorrelatedJitterStrategy) backoff(retry int) time.Duration {
	c := float64(strat.cap)
	b := float64(strat.base)
	s := float64(strat.sleep)
	u := uniform(b, 3*s)
	s = math.Min(c, u)
	return time.Duration(s)
}
