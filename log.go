package cast

import (
	"github.com/sirupsen/logrus"
)

// Logger return the underlying log instance.
func (c *Cast) Logger() *logrus.Logger {
	return c.logger
}

// LogHook log hook模板
type LogHook func(entry *logrus.Entry)

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
