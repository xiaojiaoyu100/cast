package cast

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var contextLogger = log.WithFields(log.Fields{
	"source": "cast",
})

func init() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     true})
	log.SetReportCaller(true)
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

// AddLogHook adds a log hook.
func AddLogHook(hook log.Hook) {
	log.AddHook(hook)
}
