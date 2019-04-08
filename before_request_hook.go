package cast

import (
	"github.com/jtacoma/uritemplates"
)

// BeforeRequestHook 请求之前执行的函数
type BeforeRequestHook func(cast *Cast, request *Request) error

var defaultBeforeRequestHooks = []BeforeRequestHook{
	finalizePathIfAny,
	setRequestHeader,
}

func finalizePathIfAny(cast *Cast, request *Request) error {
	if len(request.pathParam) > 0 {
		tpl, err := uritemplates.Parse(request.path)
		if err != nil {
			contextLogger.WithError(err).Error("uritemplates.Parse: ", request.path)
			return err
		}
		request.path, err = tpl.Expand(request.pathParam)
		if err != nil {
			contextLogger.WithError(err).Error("tpl.Expand: ", request.pathParam)
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
