package cast

import (
	"bytes"
	"encoding/json"
	"io"

	"strings"

	"github.com/google/go-querystring/query"
)

type reqBody interface {
	ContentType() string
	Body() (io.Reader, error)
}

type reqJsonBody struct {
	payload interface{}
}

func (body reqJsonBody) ContentType() string {
	return applicaionJson
}

func (body reqJsonBody) Body() (io.Reader, error) {
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(body.payload); err != nil {
		return nil, err
	}
	return &buffer, nil
}

type reqFormUrlEncodedBody struct {
	payload interface{}
}

func (body reqFormUrlEncodedBody) ContentType() string {
	return formUrlEncoded
}

func (body reqFormUrlEncodedBody) Body() (io.Reader, error) {
	values, err := query.Values(body.payload)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(values.Encode()), nil
}
