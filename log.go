package cast

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// LogHook log hook模板
type LogHook func(entry *logrus.Entry)

var contextLogger = log.WithFields(logrus.Fields{
	"source": "cast",
})

func init() {
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetReportCaller(true)
	log.SetOutput(os.Stderr)
	log.SetLevel(logrus.InfoLevel)
}

// AddLogHook add a log reporter.
func AddLogHook(f LogHook) {
	m := NewMonitor(f)
	log.AddHook(m)
}

// Monitor 信息监控
type Monitor struct {
	Callback LogHook
}

// NewMonitor 返回一个实例
func NewMonitor(l LogHook) *Monitor {
	m := new(Monitor)
	m.Callback = l
	return m
}

// Levels 这些级别的日志会被回调
func (m *Monitor) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}

// Fire 实际执行了回调
func (m *Monitor) Fire(entry *logrus.Entry) error {
	m.Callback(entry)
	return nil
}
