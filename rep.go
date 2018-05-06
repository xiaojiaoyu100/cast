package cast

import "encoding/json"

type Reply struct {
	statusCode int
	body []byte
}

func (rep *Reply) DecodeFromJson(v interface{}) error {
	return json.Unmarshal(rep.body, &v)
}
