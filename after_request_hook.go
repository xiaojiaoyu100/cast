package cast

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

type AfterRequestHook func(cast *Cast, request *http.Request) (*http.Request, error)

var defaultAfterRequestHooks = []AfterRequestHook{
	finalizeQueryParamIfAny,
	finalizeHeaderIfAny,
	setBearerTokenIfAny,
	setBasicAuthIfAny,
	setTimeoutIfAny,
}

func finalizeQueryParamIfAny(cast *Cast, request *http.Request) (*http.Request, error) {
	values, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		cast.logger.Printf("ERROR [%v]", err)
		return nil, err
	}

	qValues, err := query.Values(cast.queryParam)
	if err != nil {
		cast.logger.Printf("ERROR [%v]", err)
		return nil, err
	}
	for k, vv := range qValues {
		for _, v := range vv {
			values.Add(k, v)
		}
	}
	request.URL.RawQuery = values.Encode()
	return request, nil
}

func finalizeHeaderIfAny(cast *Cast, request *http.Request) (*http.Request, error) {
	if cast.body != nil && len(cast.body.ContentType()) > 0 {
		cast.SetHeader(http.Header{
			contentType: []string{cast.body.ContentType()},
		})
	}

	for k, vv := range cast.header {
		for _, v := range vv {
			request.Header.Add(k, v)
		}
	}
	return request, nil
}

func setBearerTokenIfAny(cast *Cast, request *http.Request) (*http.Request, error) {
	if len(cast.bearerToken) > 0 {
		request.Header.Set(authorization, fmt.Sprintf("Bearer %s", cast.bearerToken))
	}
	return request, nil
}

func setBasicAuthIfAny(cast *Cast, request *http.Request) (*http.Request, error) {
	if cast.basicAuth != nil {
		request.SetBasicAuth(cast.basicAuth.info())
	}
	return request, nil
}

func setTimeoutIfAny(cast *Cast, request *http.Request) (*http.Request, error) {
	if cast.timeout > 0 {
		ctx, cancel := context.WithCancel(context.TODO())
		_ = time.AfterFunc(cast.timeout, func() {
			cancel()
		})
		request = request.WithContext(ctx)
	}
	return request, nil
}
