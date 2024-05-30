package log

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/gookit/slog"
	"github.com/gookit/slog/rotatefile"
	log "github.com/sirupsen/logrus"
)

func Info(msg ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Info(msg...)
}

func Warning(msg ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Warning(msg...)
}

func Error(err ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Error(err...)
}

func Debug(value ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Debug(value...)
}

func Fatal(value ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Fatal(value...)
}

func Println(value ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Println(value...)
}

func Infof(format string, msg ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Infof(format, msg...)
}

func Warningf(format string, msg ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Warningf(format, msg...)
}

func Errorf(format string, err ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Errorf(format, err...)
}

func Debugf(format string, value ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Debugf(format, value...)
}

func Fatalf(format string, value ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Fatalf(format, value...)
}

func InitLogger(level string, logFile string) {
	logLevel := slog.DebugLevel
	switch level {
	case "debug":
		logLevel = slog.DebugLevel
	case "info":
		logLevel = slog.InfoLevel
	case "error":
		logLevel = slog.ErrorLevel
	case "warn":
		logLevel = slog.WarnLevel
	}
	slog.SetLogLevel(logLevel)
	logTemplate := "[{{level}}] [{{datetime}}] [{{meta}}] Message: {{message}} {{data}} \n"

	slog.SetFormatter(slog.NewTextFormatter(logTemplate).WithEnableColor(true))
	writer, err := rotatefile.NewConfig(logFile).Create()
	if err != nil {
		panic(err)
	}

	log.SetOutput(writer)
}
