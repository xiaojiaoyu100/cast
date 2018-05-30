package cast

import (
	"net/http/httputil"
)

type responseHook func(c *Cast, response *Response) error

var defaultResponseHooks = []responseHook{
	dumpResponse,
}

func dumpResponse(c *Cast, response *Response) error {
	if !globalLogger.debug {
		return nil
	}

	if response.rawResponse.ContentLength > c.dumpBodyLimit {
		return nil
	}

	bytes, err := httputil.DumpResponse(response.rawResponse, true)
	if err != nil {
		globalLogger.printf("ERROR [%v]", err)
		return err
	}

	globalLogger.printf("%s\n\n%s took %s upto %d time(s)", bytes, response.rawResponse.Request.URL.String(), response.cost, response.times)
	return nil
}
