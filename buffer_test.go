package cast

import (
	"bytes"
	"testing"
)

const (
	str = `ivqEBnJ8/+3L7Cr1TN6t7yxhelXREMCB6LmXif1BaLzSOhUcLPDYqapdlKTd1D/KH9D8pmnRZ9OGuSJ7kFkksZxLNQjuqtjvKICpJI5I2R8nmZ1VdEAGpp2xbIE1m+alSKt9zg==`
)

func BenchmarkGetBufferNoPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.Buffer{}
		buf.WriteString(str)
	}
}

func BenchmarkGetBufferWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := getBuffer()
		buf.WriteString(str)
		putBuffer(buf)
	}
}

func BenchmarkGetBufferWithPoolNoPut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := getBuffer()
		buf.WriteString(str)
	}
}
