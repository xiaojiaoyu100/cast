package cast

import (
	"net/http"
)

type ResponseHook func(cast *Cast, response *http.Response) error



