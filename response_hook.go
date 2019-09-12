package cast

import (
	"bytes"
	"fmt"
)

type responseHook func(c *Cast, response *Response) error

var defaultResponseHooks = []responseHook{
	dump,
}

func dump(cast *Cast, response *Response) error {
	buffer := getBuffer()
	defer putBuffer(buffer)

	var err error
	prt := func(buffer *bytes.Buffer, format string, a ...interface{}) {
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(buffer, format, a...)
	}

	shouldPrintHeaders := cast.dumpFlag&fHeader != 0
	if shouldPrintHeaders {
		prt(buffer, "\nHeaders\n")
		prt(buffer, "Request URL: %s\n", response.request.rawRequest.URL.String())
		prt(buffer, "Request Method: %s\n", response.request.method)
		prt(buffer, "Remote Address: %s\n", response.request.remoteAddress)
		prt(buffer, "Status Code: %s\n", response.rawResponse.Status)
		prt(buffer, "Version: %s\n", response.rawResponse.Proto)
		prt(buffer, "Response Headers\n")
		if err != nil {
			return err
		}

		err = response.rawResponse.Header.Write(buffer)
		if err != nil {
			return err
		}

		prt(buffer, "Request Headers\n")
		if err != nil {
			return err
		}

		err = response.request.rawRequest.Header.Write(buffer)
		if err != nil {
			return err
		}

		prt(buffer, "\n")
		if err != nil {
			return err
		}
	}

	shouldPrintParams := cast.dumpFlag&fParam != 0
	if shouldPrintParams && len(response.body) <= defaultDumpBodyLimit {
		prt(buffer, "Params\n")
		if err != nil {
			return err
		}

		if response.request.body != nil {
			body, _ := response.request.body.Body()
			prt(buffer, string(body))
		}

		prt(buffer, "\n")
		if err != nil {
			return err
		}

	}

	shouldPrintResponse := cast.dumpFlag&fResponse != 0
	if shouldPrintResponse && len(response.body) <= defaultDumpBodyLimit {
		prt(buffer, "Response\n")
		prt(buffer, string(response.body))
		prt(buffer, "\n\n")
		if err != nil {
			return err
		}
	}

	shouldPrintTimings := cast.dumpFlag&fTiming != 0
	if shouldPrintTimings {
		prt(buffer, "Timings\n")
		prt(buffer, "DNS resolution: %s\n", response.request.prof.dnsCost)
		prt(buffer, "Connecting: %s\n", response.request.prof.connectCost)
		if err != nil {
			return err
		}
		if response.request.prof.tlsHandshakeCost.Nanoseconds() > 0 {
			prt(buffer, "TLS setup: %s\n", response.request.prof.tlsHandshakeCost)
			if err != nil {
				return err
			}
		}
		prt(buffer, "Sending: %s\n", response.request.prof.sendingCost.String())
		prt(buffer, "Waiting: %s\n", response.request.prof.waitingCost.String())
		prt(buffer, "Receiving: %s\n", response.request.prof.receivingCost.String())
		prt(buffer, "All: %s\n", response.request.prof.requestCost.String())
		if err != nil {
			return err
		}
	}

	if buffer.Len() > 0 {
		cast.Logger().Info(buffer.String())
	}

	return nil
}
