package cast

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"crypto/tls"
	"net/http/httptrace"

	"github.com/google/go-querystring/query"
)

type RequestHook func(cast *Cast, request *Request) error

var defaultRequestHooks = []RequestHook{
	finalizeQueryParamIfAny,
	finalizeAuthorization,
	addCookies,
	finalizeHeaderIfAny,
	setTimeoutIfAny,
	clientTrace,
}

func finalizeQueryParamIfAny(cast *Cast, request *Request) error {
	values, err := url.ParseQuery(request.rawRequest.URL.RawQuery)
	if err != nil {
		cast.Logger().WithError(err).Error("url.ParseQuery")
		return err
	}

	qValues, err := query.Values(request.queryParam)
	if err != nil {
		cast.Logger().WithError(err).Error(" query.Values")
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

func setTimeoutIfAny(_ *Cast, request *Request) error {
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

func clientTrace(_ *Cast, request *Request) error {
	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			request.prof.waitingDone = time.Now().In(time.UTC)
			request.prof.waitingCost = request.prof.waitingDone.Sub(request.prof.waitingStart)
			request.prof.receivingSart = time.Now().In(time.UTC)
		},
		DNSStart: func(httptrace.DNSStartInfo) {
			request.prof.dnsStart = time.Now().In(time.UTC)
		},
		DNSDone: func(httptrace.DNSDoneInfo) {
			request.prof.dnsDone = time.Now().In(time.UTC)
			request.prof.dnsCost = request.prof.dnsDone.Sub(request.prof.dnsStart)
		},
		ConnectStart: func(network, addr string) {
			request.prof.connectStart = time.Now().In(time.UTC)
		},
		ConnectDone: func(network, addr string, err error) {
			request.prof.connectDone = time.Now().In(time.UTC)
			request.prof.connectCost = request.prof.connectDone.Sub(request.prof.connectStart)
			request.remoteAddress = addr
		},
		TLSHandshakeStart: func() {
			request.prof.tlsHandshakeStart = time.Now().In(time.UTC)
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			request.prof.tlsHandshakeDone = time.Now().In(time.UTC)
			request.prof.tlsHandshakeCost = request.prof.tlsHandshakeDone.Sub(request.prof.tlsHandshakeStart)
		},
		WroteHeaders: func() {
			request.prof.sendingStart = time.Now().In(time.UTC)
		},
		WroteRequest: func(httptrace.WroteRequestInfo) {
			request.prof.sendingDone = time.Now().In(time.UTC)
			request.prof.sendingCost = request.prof.sendingDone.Sub(request.prof.sendingStart)
			request.prof.waitingStart = time.Now().In(time.UTC)
		},
	}
	request.rawRequest = request.rawRequest.WithContext(httptrace.WithClientTrace(request.rawRequest.Context(), trace))
	return nil
}
