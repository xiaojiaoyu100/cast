package cast

import (
	"net/http"
	"reflect"
	"testing"
)

func TestRequest_WithHeader(t *testing.T) {

	tests := [...]struct {
		header http.Header
		want   http.Header
	}{
		0: {
			header: nil,
			want:   nil,
		},
		1: {
			header: map[string][]string{
				"Content-Type": []string{"application/json"},
			},
			want: map[string][]string{
				"Content-Type": []string{"application/json"},
			},
		},
	}

	for i, tt := range tests {
		request := new(Request)
		request.WithHeader(tt.header)
		assert(t, reflect.DeepEqual(request.header, tt.want), "%d: unexpected WithHeader", i)
	}
}
