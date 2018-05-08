package cast

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type dumpResponseHook func(logger *log.Logger, request *http.Response)

func dumpResponse(logger *log.Logger, response *http.Response) {
	if logger == nil || response == nil {
		return
	}

	body, err := httputil.DumpResponse(response, true)
	if err != nil {
		return
	}

	logger.Printf("%s", body)
}
