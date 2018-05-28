package cast

import (
	"time"

	"github.com/jtacoma/uritemplates"
)

type BeforeRequestHook func(cast *Cast) error

var defaultBeforeRequestHooks = []BeforeRequestHook{
	requestStart,
	finalizePathIfAny,
}

func requestStart(cast *Cast) error {
	if cast != nil {
		cast.start = time.Now().In(time.UTC)
	}
	return nil
}

func finalizePathIfAny(cast *Cast) error {
	if len(cast.pathParam) > 0 {
		tpl, err := uritemplates.Parse(cast.path)
		if err != nil {
			cast.logger.Printf("ERROR [%v]", err)
			return err
		}
		cast.path, err = tpl.Expand(cast.pathParam)
		if err != nil {
			cast.logger.Printf("ERROR [%v]", err)
			return err
		}
	}
	return nil
}
