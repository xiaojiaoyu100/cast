package cast

import "testing"

func TestResponse_Success(t *testing.T) {
	tests := [...]struct {
		statusCode int
		want       bool
	}{
		0: {
			statusCode: 199,
			want:       false,
		},
		1: {
			statusCode: 200,
			want:       true,
		},
		2: {
			statusCode: 201,
			want:       true,
		},
		3: {
			statusCode: 299,
			want:       true,
		},
		4: {
			statusCode: 300,
			want:       false,
		},
		5: {
			statusCode: 350,
			want:       false,
		},
		6: {
			statusCode: -1,
			want:       false,
		},
	}

	for i, tt := range tests {
		response := new(Response)
		response.statusCode = tt.statusCode
		assert(t, response.Success() == tt.want, "%d: unexpected Success()", i)
	}
}

func TestResponse_String(t *testing.T) {
	tests := [...]struct {
		body []byte
		want string
	}{
		0: {
			body: []byte("ss"),
			want: "ss",
		},
		1: {
			body: nil,
			want: "",
		},
	}

	for i, tt := range tests {
		response := new(Response)
		response.body = tt.body
		assert(t, response.String() == tt.want, "%d: unexpected String()", i)
	}

}
