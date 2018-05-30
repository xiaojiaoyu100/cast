package cast

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestSetDebug(t *testing.T) {
	tests := [...]struct {
		l logger
	}{
		0: {
			l: log.New(os.Stderr, "", log.LstdFlags),
		},
		1: {
			l: nil,
		},
	}

	for i, tt := range tests {
		SetDebug(tt.l)
		assert(t, globalLogger.l == tt.l, "%d: unexpected SetDebug", i)
		assert(t, globalLogger.debug == true, "%d: unexpected SetDebug", i)
	}
}

func TestQuickDebug(t *testing.T) {
	QuickDebug()
	assert(t, globalLogger.l != nil, "%d: unexpected QuickDebug")
	assert(t, globalLogger.debug == true, "%d: unexpected QuickDebug")
}

func ExampleDebugLoggerPrintf() {
	var (
		logger debugLogger
		buffer bytes.Buffer
	)
	logger.debug = true
	logger.l = log.New(&buffer, "", 0)
	logger.printf("%s: %d", "hello", 100)
	fmt.Printf(buffer.String())
	// Output:
	// hello: 100
}
