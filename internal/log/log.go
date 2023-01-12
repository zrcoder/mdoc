package log

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func Enable() {
	logger = newLogger()
}

func Info(v ...any) {
	trace(func() { logger.Infoln(v...) })
}
func Infof(format string, v ...any) {
	trace(func() { logger.Infof(format, v...) })
}
func Error(v ...any) {
	trace(func() { logger.Error(v...) })
}
func Errorf(format string, v ...any) {
	trace(func() { logger.Errorf(format, v...) })
}

func trace(f func()) {
	if logger == nil {
		return
	}
	fmt.Println(makeTraceInfo())
	f()
}

func newLogger() *logrus.Logger {
	lg := logrus.New()
	lg.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	return lg
}

func makeTraceInfo() string {
	_, filename, line, _ := runtime.Caller(3)
	return fmt.Sprintf("%s:%d", filename, line)
}
