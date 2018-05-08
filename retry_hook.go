package cast

import "net/http"

type RetryHook func(resp *http.Response) error
