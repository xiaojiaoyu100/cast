package cast

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRetryWhenTooManyRequests(t *testing.T) {
	tests := [...]struct {
		handler func(w http.ResponseWriter, r *http.Request)
		want    bool
	}{
		0: {
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTooManyRequests)
			},
			want: true,
		},
		1: {
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			want: false,
		},
	}

	for i, tt := range tests {
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		tt.handler(w, req)
		resp := w.Result()
		err := RetryWhenTooManyRequests(resp)
		assert(t, tt.want == err, "%d: unexpected RetryWhenTooManyRequests return", i)
	}
}

func TestRetryWhenInternalServerError(t *testing.T) {
	tests := [...]struct {
		handler func(w http.ResponseWriter, r *http.Request)
		want    bool
	}{
		0: {
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			want: true,
		},
		1: {
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			want: false,
		},
	}

	for i, tt := range tests {
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		tt.handler(w, req)
		resp := w.Result()
		err := RetryWhenInternalServerError(resp)
		assert(t, tt.want == err, "%d: unexpected RetryWhenInternalServerError return", i)
	}
}
