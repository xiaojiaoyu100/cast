package cast

import (
	"encoding/json"
	"encoding/xml"
	"io"

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
	buffer := getBuffer()
	if err := json.NewEncoder(buffer).Encode(body.payload); err != nil {
		return nil, err
	}
	return buffer, nil
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
	buffer := getBuffer()
	buffer.WriteString(values.Encode())
	return buffer, nil
}

type reqXmlBody struct {
	payload interface{}
}

func (body reqXmlBody) Body() (io.Reader, error) {
	buffer := getBuffer()
	if err := xml.NewEncoder(buffer).Encode(body.payload); err != nil {
		return nil, err
	}
	return buffer, nil
}

func (body reqXmlBody) ContentType() string {
	return applicationXml
}

type reqPlainBody struct {
	payload string
}

func (body reqPlainBody) ContentType() string {
	return textPlain
}

func (body reqPlainBody) Body() (io.Reader, error) {
	buffer := getBuffer()
	buffer.WriteString(body.payload)
	return buffer, nil
}
