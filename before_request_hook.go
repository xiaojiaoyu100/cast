package cast

import (
	"github.com/jtacoma/uritemplates"
)

type beforeRequestHook func(cast *Cast, request *Request) error

var defaultBeforeRequestHooks = []beforeRequestHook{
	finalizePathIfAny,
	setRequestHeader,
}

func finalizePathIfAny(cast *Cast, request *Request) error {
	if len(request.pathParam) > 0 {
		tpl, err := uritemplates.Parse(request.path)
		if err != nil {
			contextLogger.WithError(err)
			return err
		}
		request.path, err = tpl.Expand(request.pathParam)
		if err != nil {
			contextLogger.WithError(err)
			return err
		}
	}
	return nil
}

func setRequestHeader(cast *Cast, request *Request) error {
	if request.body != nil && len(request.body.ContentType()) > 0 {
		request.SetHeader(contentType, request.body.ContentType())
	}
	return nil
}
