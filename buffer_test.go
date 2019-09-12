package cast

import (
	"bytes"
	"math/rand"
	"testing"
)

const (
	str = `ivqEBnJ8/+3L7Cr1TN6t7yxhelXREMCB6LmXif1BaLzSOhUcLPDYqapdlKTd1D
/KH9D8pmnRZ9OGuSJ7kFkksZxLNQjuqtjvKICpJI5I2R8nmZ1VdEAGpp2xbIE1m+alSKt9zg==`
)

func BenchmarkBufferNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randNum := rand.Intn(8192) + 1
		data := make([]byte, 0, randNum)
		_, err := rand.Read(data)
		ok(b, err)
		buffer := bytes.Buffer{}
		buffer.WriteString(string(data))
	}
}

func BenchmarkBufferWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randNum := rand.Intn(8192) + 1
		data := make([]byte, 0, randNum)
		_, err := rand.Read(data)
		ok(b, err)

		buffer := getBuffer()
		buffer.WriteString(string(data))
		putBuffer(buffer)
	}
}

func BenchmarkBufferWithPoolNoPut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randNum := rand.Intn(8192) + 1
		data := make([]byte, 0, randNum)
		_, err := rand.Read(data)
		ok(b, err)
		buffer := getBuffer()
		buffer.WriteString(string(data))
	}
}

func BenchmarkGetBufferNoPoolFixedString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.Buffer{}
		buf.WriteString(str)
	}
}

func BenchmarkGetBufferWithPoolFixedString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := getBuffer()
		buf.WriteString(str)
		putBuffer(buf)
	}
}

func BenchmarkGetBufferWithPoolNoPutFixedString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := getBuffer()
		buf.WriteString(str)
	}
}
