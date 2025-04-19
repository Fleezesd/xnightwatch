package x

import "github.com/fleezesd/xnightwatch/pkg/log"

type cronLogger struct{}

// cronLogger implement the cron.Logger interface.
func NewLogger() *cronLogger {
	return &cronLogger{}
}

// Debug logs routine messages about cron's operation.
func (l *cronLogger) Debug(msg string, keysAndValues ...interface{}) {
	log.Debugw(msg, keysAndValues...)
}

// Info logs routine messages about cron's operation.
func (l *cronLogger) Info(msg string, keysAndValues ...any) {
	log.Infow(msg, keysAndValues...)
}

// Error logs an error condition.
func (l *cronLogger) Error(err error, msg string, keysAndValues ...any) {
	log.Errorw(err, msg, keysAndValues...)
}
