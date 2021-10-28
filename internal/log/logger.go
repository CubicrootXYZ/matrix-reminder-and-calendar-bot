package log

import (
	"go.uber.org/zap"
)

// ANSI colors: https://gist.github.com/JBlond/2fea43a3049b38287e5e9cefc87b2124
var (
	None   = "0"
	Red    = "31"
	Green  = "32"
	Yellow = "33"
	Blue   = "34"
)

var logger *zap.SugaredLogger

// InitLogger initializes a new logger. Make sure to call defer logger.Sync().
func InitLogger() *zap.SugaredLogger {
	log, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger = log.Sugar()
	return logger
}

// Debug logs with tag debug
func Debug(msg string) {
	logger.Debug("msg")
}

// UInfo unstructured info log
func UInfo(msg string, args ...interface{}) {
	logger.Infow(msg, args...)
}

// Info logs with tag info and in blue
func Info(msg string) {
	logger.Info(msg)
}

// Warn logs with tag warn and in yellow
func Warn(msg string) {
	logger.Warn(msg)
}

// Error logs with tag error and in red
func Error(msg string) {
	logger.Error(msg)
}
