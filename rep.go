package cast

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"time"
)

type Reply struct {
	url        string
	statusCode int
	body       []byte
	size       int64
	header     http.Header
	cookies    []*http.Cookie
	start      time.Time
	end        time.Time
	cost       time.Duration
	times      int
}

func (rep *Reply) Url() string {
	return rep.url
}

func (rep *Reply) DecodeFromJson(v interface{}) error {
	return json.Unmarshal(rep.body, &v)
}

func (rep *Reply) DecodeFromXml(v interface{}) error {
	return xml.Unmarshal(rep.body, &v)
}

func (rep *Reply) Body() []byte {
	return rep.body
}

func (rep *Reply) Size() int64 {
	return rep.size
}

func (rep *Reply) Header() http.Header {
	return rep.header
}

func (rep *Reply) Cookies() []*http.Cookie {
	return rep.cookies
}

func (rep *Reply) StatusOk() bool {
	return rep.statusCode == http.StatusOK
}

func (rep *Reply) StatusCode() int {
	return rep.statusCode
}

func (rep *Reply) Start() time.Time {
	return rep.start
}

func (rep *Reply) End() time.Time {
	return rep.end
}

func (rep *Reply) Cost() time.Duration {
	return rep.cost
}
