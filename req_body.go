package cast

import (
	"io"
	"bytes"
	"encoding/json"
)

type ReqBody interface {
	ContentType() string
	Body() (io.Reader, error)
}

type ReqJsonBody struct {
	payload interface{}
}

func (body ReqJsonBody) Body() (io.Reader, error){
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(body.payload); err != nil {
		return nil, err
	}
	return &buffer, nil
}
