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

func (stg linearBackoffStrategy) backoff(retry int) time.Duration {
	return time.Duration(retry) * stg.slope
}

type constantBackOffStrategy struct {
	interval time.Duration
}

func (stg constantBackOffStrategy) backoff(retry int) time.Duration {
	return stg.interval
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

func (stg exponentialBackoffStrategy) backoff(retry int) time.Duration {
	return time.Duration(stg.expo(retry))
}

type exponentialBackoffEqualJitterStrategy struct {
	exponentialBackoff
}

func (stg exponentialBackoffEqualJitterStrategy) backoff(retry int) time.Duration {
	v := stg.expo(retry)
	u := uniform(0, v/2.0)
	return time.Duration(v/2.0 + u)
}

type exponentialBackoffFullJitterStrategy struct {
	exponentialBackoff
}

func (stg exponentialBackoffFullJitterStrategy) backoff(retry int) time.Duration {
	v := stg.expo(retry)
	u := uniform(0, v)
	return time.Duration(u)
}

type exponentialBackoffDecorrelatedJitterStrategy struct {
	exponentialBackoff
	sleep time.Duration
}

// uniform returns a number in [min, max)
func uniform(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (stg exponentialBackoffDecorrelatedJitterStrategy) backoff(retry int) time.Duration {
	c := float64(stg.cap)
	b := float64(stg.base)
	s := float64(stg.sleep)
	u := uniform(b, 3*s)
	s = math.Min(c, u)
	return time.Duration(s)
}
