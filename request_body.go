package cast

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"
	"os"

	"path/filepath"

	"github.com/google/go-querystring/query"
	"bytes"
)

type requestBody interface {
	ContentType() string
	Body() ([]byte, error)
}

type requestJsonBody struct {
	payload interface{}
}

func (body *requestJsonBody) ContentType() string {
	return applicaionJson
}

func (body *requestJsonBody) Body() ([]byte, error) {
	return json.Marshal(body.payload)
}

type requestFormUrlEncodedBody struct {
	payload interface{}
}

func (body *requestFormUrlEncodedBody) ContentType() string {
	return formUrlEncoded
}

func (body *requestFormUrlEncodedBody) Body() ([]byte, error) {
	values, err := query.Values(body.payload)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

type requestXmlBody struct {
	payload interface{}
}

func (body *requestXmlBody) Body() ([]byte, error) {
	return xml.Marshal(body.payload)
}

func (body *requestXmlBody) ContentType() string {
	return applicationXml
}

type requestPlainBody struct {
	payload string
}

func (body *requestPlainBody) ContentType() string {
	return textPlain
}

func (body *requestPlainBody) Body() ([]byte, error) {
	return []byte(body.payload), nil
}

type FormData struct {
	FieldName string
	Value     string
	FileName  string
	Path      string
	Reader    io.Reader
}

type requestMultipartFormDataBody struct {
	formData    []*FormData
	contentType string
}

func (body *requestMultipartFormDataBody) ContentType() string {
	return body.contentType
}

func (body *requestMultipartFormDataBody) Body() ([]byte, error) {
	buffer := &bytes.Buffer{}
	w := multipart.NewWriter(buffer)
	defer w.Close()

	for _, data := range body.formData {
		switch {
		case len(data.Value) != 0:
			if err := w.WriteField(data.FieldName, data.Value); err != nil {
				return nil, err
			}
		default:
			if len(data.FileName) == 0 || len(data.FileName) == 0 {
				continue
			}

			if len(data.Path) == 0 && data.Reader == nil {
				continue
			}

			fw, err := w.CreateFormFile(data.FieldName, data.FileName)
			if err != nil {
				return nil, err
			}

			switch {
			case len(data.Path) > 0:
				path, err := filepath.Abs(data.Path)
				if err != nil {
					return nil, err
				}
				f, err := os.Open(path)
				if err != nil {
					return nil, err
				}
				defer f.Close()
				_, err = io.Copy(fw, f)
				if err != nil {
					return nil, err
				}

			case data.Reader != nil:
				_, err := io.Copy(fw, data.Reader)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	body.contentType = w.FormDataContentType()
	return buffer.Bytes(), nil
}
