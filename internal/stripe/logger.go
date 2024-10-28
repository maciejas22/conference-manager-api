package stripe

import (
	"log/slog"
)

type StripeLogger struct {
	logger *slog.Logger
}

func NewStripeLogger(l *slog.Logger) *StripeLogger {
	return &StripeLogger{logger: l}
}

func (l *StripeLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debug(format, v...)
}

func (l *StripeLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error(format, v...)
}

func (l *StripeLogger) Infof(format string, v ...interface{}) {
	l.logger.Info(format, v...)
}

func (l *StripeLogger) Warnf(format string, v ...interface{}) {
	l.logger.Warn(format, v...)
}
