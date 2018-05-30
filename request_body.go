package cast

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"
	"os"

	"path/filepath"

	"github.com/google/go-querystring/query"
)

type requestBody interface {
	ContentType() string
	Body() (io.Reader, error)
}

type requestJsonBody struct {
	payload interface{}
}

func (body *requestJsonBody) ContentType() string {
	return applicaionJson
}

func (body *requestJsonBody) Body() (io.Reader, error) {
	buffer := getBuffer()
	if err := json.NewEncoder(buffer).Encode(body.payload); err != nil {
		return nil, err
	}
	return buffer, nil
}

type requestFormUrlEncodedBody struct {
	payload interface{}
}

func (body *requestFormUrlEncodedBody) ContentType() string {
	return formUrlEncoded
}

func (body *requestFormUrlEncodedBody) Body() (io.Reader, error) {
	values, err := query.Values(body.payload)
	if err != nil {
		return nil, err
	}
	buffer := getBuffer()
	buffer.WriteString(values.Encode())
	return buffer, nil
}

type requestXmlBody struct {
	payload interface{}
}

func (body *requestXmlBody) Body() (io.Reader, error) {
	buffer := getBuffer()
	if err := xml.NewEncoder(buffer).Encode(body.payload); err != nil {
		return nil, err
	}
	return buffer, nil
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

func (body *requestPlainBody) Body() (io.Reader, error) {
	buffer := getBuffer()
	buffer.WriteString(body.payload)
	return buffer, nil
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

func (body *requestMultipartFormDataBody) Body() (io.Reader, error) {
	buffer := getBuffer()

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
	return buffer, nil
}
