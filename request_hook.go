package cast

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type dumpRequestHook func(logger *log.Logger, request *http.Request)

func dumpRequest(logger *log.Logger, request *http.Request) {
	if logger == nil || request == nil {
		return
	}

	body, err := httputil.DumpRequest(request, true)
	if err != nil {
		return
	}

	logger.Printf("%s", string(body))
}
