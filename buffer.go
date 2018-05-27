package cast

import (
	"sync"
	"bytes"
	"io"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func getBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

func putBuffer(reader io.Reader) {
	if reader == nil {
		return
	}

	buffer, ok := reader.(*bytes.Buffer)
	if !ok {
		return
	}

	buffer.Reset()
	bufferPool.Put(buffer)
}




