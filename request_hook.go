package cast

import (
	"context"
	"net/url"
	"time"

	"fmt"

	"net/http/httputil"

	"github.com/google/go-querystring/query"
)

type requestHook func(cast *Cast, request *Request) error

var defaultRequestHooks = []requestHook{
	finalizeQueryParamIfAny,
	finalizeAuthorization,
	addCookies,
	finalizeHeaderIfAny,
	setTimeoutIfAny,
	dumpRequest,
}

func finalizeQueryParamIfAny(cast *Cast, request *Request) error {
	values, err := url.ParseQuery(request.rawRequest.URL.RawQuery)
	if err != nil {
		globalLogger.printf("ERROR [%v]", err)
		return err
	}

	qValues, err := query.Values(request.queryParam)
	if err != nil {
		globalLogger.printf("ERROR [%v]", err)
		return err
	}
	for k, vv := range qValues {
		for _, v := range vv {
			values.Add(k, v)
		}
	}
	request.rawRequest.URL.RawQuery = values.Encode()
	return nil
}

func finalizeAuthorization(cast *Cast, request *Request) error {
	switch {
	case len(cast.bearerToken) > 0:
		request.rawRequest.Header.Set(authorization, fmt.Sprintf("Bearer %s", cast.bearerToken))
	case cast.basicAuth != nil:
		request.rawRequest.SetBasicAuth(cast.basicAuth.info())
	}
	return nil
}

func addCookies(cast *Cast, request *Request) error {
	for _, cookie := range cast.cookies {
		request.rawRequest.AddCookie(cookie)
	}
	return nil
}

func finalizeHeaderIfAny(cast *Cast, request *Request) error {
	for k, vv := range cast.header {
		for _, v := range vv {
			request.rawRequest.Header.Add(k, v)
		}
	}

	for k, vv := range request.header {
		for _, v := range vv {
			request.rawRequest.Header.Add(k, v)
		}
	}

	return nil
}

func setTimeoutIfAny(cast *Cast, request *Request) error {
	if request.timeout == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(context.TODO())
	_ = time.AfterFunc(request.timeout, func() {
		cancel()
	})
	request.rawRequest = request.rawRequest.WithContext(ctx)

	return nil
}

func dumpRequest(cast *Cast, request *Request) error {
	if !globalLogger.debug {
		return nil
	}

	if request.rawRequest.ContentLength > cast.dumpBodyLimit {
		return nil
	}

	bytes, err := httputil.DumpRequest(request.rawRequest, true)
	if err != nil {
		return err
	}

	globalLogger.printf("%s", bytes)

	return nil
}
