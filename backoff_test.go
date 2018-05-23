package cast

import (
	"fmt"
	"testing"
	"time"
)

func ExampleLinearBackoff() {
	s := linearBackoffStrategy{
		slope: 1 * time.Second,
	}
	for i := 1; i <= 3; i++ {
		fmt.Println(s.backoff(i))
	}
	// Output:
	// 1s
	// 2s
	// 3s
}

func ExampleConstantBackOff() {
	s := constantBackOffStrategy{
		interval: 2 * time.Second,
	}
	for i := 1; i <= 3; i++ {
		fmt.Println(s.backoff(i))
	}
	// Output:
	// 2s
	// 2s
	// 2s
}

func ExampleExpo() {
	expoBackoff := exponentialBackoff{
		base: time.Nanosecond,
		cap:  10 * time.Nanosecond,
	}
	for i := 1; i <= 3; i++ {
		fmt.Printf("%f\n", expoBackoff.expo(i))
	}
	// Output:
	// 2.000000
	// 4.000000
	// 8.000000
}

func ExampleExponentialBackoff() {
	s := exponentialBackoffStrategy{
		exponentialBackoff: exponentialBackoff{
			base: time.Second,
			cap:  10 * time.Second,
		},
	}
	for i := 1; i <= 5; i++ {
		fmt.Println(s.backoff(i))
	}
	// Output:
	// 2s
	// 4s
	// 8s
	// 10s
	// 10s
}

func BenchmarkExampleExponentialBackoffEqualJitter(b *testing.B) {
	s := exponentialBackoffEqualJitterStrategy{
		exponentialBackoff: exponentialBackoff{
			base: time.Second,
			cap:  10 * time.Second,
		},
	}
	for i := 1; i <= 5; i++ {
		fmt.Println(s.backoff(i))
	}
}

func BenchmarkExponentialBackoffFullJitter(b *testing.B) {
	s := exponentialBackoffFullJitterStrategy{
		exponentialBackoff: exponentialBackoff{
			base: time.Second,
			cap:  10 * time.Second,
		},
	}
	for i := 1; i <= 5; i++ {
		b.Log(s.backoff(i))
	}
}

func BenchmarkExponentialBackoffDecorrelatedJitter(b *testing.B) {
	base := time.Second
	s := exponentialBackoffDecorrelatedJitterStrategy{
		exponentialBackoff: exponentialBackoff{
			base: base,
			cap:  10 * time.Second,
		},
		sleep: base,
	}
	for i := 1; i <= 5; i++ {
		b.Log(s.backoff(i))
	}
}

func BenchmarkUniform(b *testing.B) {
	for n := b.N; n >= 0; n-- {
		min := 3.0
		max := 10.0
		u := uniform(min, max)
		if u < min || u >= max {
			b.Fatal("unexpected uniform")
		}
	}
}
