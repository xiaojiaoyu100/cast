package cast

import (
	"bytes"
	"encoding/json"
	"io"

	"strings"

	"github.com/google/go-querystring/query"
)

type ReqBody interface {
	ContentType() string
	Body() (io.Reader, error)
}

type ReqJsonBody struct {
	payload interface{}
}

func (body ReqJsonBody) ContentType() string {
	return applicaionJson
}

func (body ReqJsonBody) Body() (io.Reader, error) {
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(body.payload); err != nil {
		return nil, err
	}
	return &buffer, nil
}

type ReqFormUrlEncodedBody struct {
	payload interface{}
}

func (body ReqFormUrlEncodedBody) ContentType() string {
	return formUrlEncoded
}

func (body ReqFormUrlEncodedBody) Body() (io.Reader, error) {
	values, err := query.Values(body.payload)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(values.Encode()), nil
}
