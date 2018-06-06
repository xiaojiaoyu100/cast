package cast

import (
	"fmt"
)

type responseHook func(c *Cast, response *Response) error

var defaultResponseHooks = []responseHook{
	dump,
}

func dump(cast *Cast, response *Response) error {
	if !globalLogger.debug {
		return nil
	}

	buffer := getBuffer()
	defer putBuffer(buffer)

	shouldPrintHeaders := cast.dumpFlag & fHeader != 0
	if shouldPrintHeaders {
		fmt.Fprintf(buffer, "\nHeaders\n")
		fmt.Fprintf(buffer, "Request URL: %s\n", response.request.rawRequest.URL.String())
		fmt.Fprintf(buffer, "Request Method: %s\n", response.request.method)
		fmt.Fprintf(buffer, "Remote Address: %s\n", response.request.remoteAddress)
		fmt.Fprintf(buffer, "Status Code: %s\n", response.rawResponse.Status)
		fmt.Fprintf(buffer, "Version: %s\n", response.rawResponse.Proto)
		fmt.Fprintf(buffer, "Response Headers\n")
		response.rawResponse.Header.Write(buffer)
		fmt.Fprintf(buffer, "Request Headers\n")
		response.request.rawRequest.Header.Write(buffer)
		fmt.Fprintf(buffer, "\n")
	}

	shouldPrintParams := cast.dumpFlag & fParam != 0
	if shouldPrintParams && len(response.body) <= defaultDumpBodyLimit {
		fmt.Fprintf(buffer, "Params\n")
		if response.request.body != nil {
			body, _ := response.request.body.Body()
			fmt.Fprintf(buffer, string(body))
		}
		fmt.Fprintf(buffer, "\n")

	}

	shouldPrintResponse := cast.dumpFlag & fResponse != 0
	if shouldPrintResponse && len(response.body) <= defaultDumpBodyLimit {
		fmt.Fprintf(buffer, "Response\n")
		fmt.Fprintf(buffer, string(response.body))
		fmt.Fprintf(buffer, "\n\n")
	}

	shouldPrintTimings := cast.dumpFlag & fTiming != 0
	if shouldPrintTimings {
		fmt.Fprintf(buffer, "Timings\n")
		fmt.Fprintf(buffer, "DNS resolution: %s\n", response.request.prof.dnsCost)
		fmt.Fprintf(buffer, "Connecting: %s\n", response.request.prof.connectCost)
		if response.request.prof.tlsHandshakeCost.Nanoseconds() > 0 {
			fmt.Fprintf(buffer, "TLS setup: %s\n", response.request.prof.tlsHandshakeCost)
		}
		fmt.Fprintf(buffer, "Sending: %s\n", response.request.prof.sendingCost.String())
		fmt.Fprintf(buffer, "Waiting: %s\n", response.request.prof.waitingCost.String())
		fmt.Fprintf(buffer, "Receiving: %s\n", response.request.prof.receivingCost.String())
		fmt.Fprintf(buffer, "All: %s\n", response.request.prof.requestCost.String())
	}

	if buffer.Len() > 0 {
		globalLogger.printf("%s", buffer.String())
	}

	return nil
}
