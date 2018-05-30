package cast

import "testing"

func Test_finalizePathIfAny(t *testing.T) {
	tests := [...]struct {
		path      string
		pathParam map[string]interface{}
		want      string
	}{
		0: {
			path: "/{1}/{2}/{3}",
			pathParam: map[string]interface{}{
				"1": "cd",
				"2": "to",
				"3": "home",
			},
			want: "/cd/to/home",
		},
		1: {
			path: "/{1}/u/{2}",
			pathParam: map[string]interface{}{
				"1": "are",
				"2": "ok",
			},
			want: "/are/u/ok",
		},
	}

	for i, tt := range tests {
		request := new(Request)
		request.pathParam = tt.pathParam
		request.path = tt.path
		err := finalizePathIfAny(nil, request)
		ok(t, err)

		assert(t, request.path == tt.want, "%d: finalizePathIfAny error", i)
	}
}
