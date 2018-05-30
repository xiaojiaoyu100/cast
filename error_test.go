package cast

import "testing"

func TestCastError_Error(t *testing.T) {
	tests := [...]struct {
		err  CastError
		want string
	}{
		0: {
			err:  CastError(""),
			want: "",
		},
		1: {
			err:  CastError("123"),
			want: "123",
		},
	}

	for i, tt := range tests {
		assert(t, tt.err.Error() == tt.want, "%d unexpected Error()", i)
	}
}
