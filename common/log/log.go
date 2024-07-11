package log

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/gookit/slog"
	"github.com/gookit/slog/rotatefile"
)

func Info(msg ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	slog.WithFields(slog.M{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Info(msg...)
}

func Warn(msg ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	slog.WithFields(slog.M{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Warn(msg...)
}

func Error(err ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	slog.WithFields(slog.M{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Error(err...)
}

func Debug(value ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	slog.WithFields(slog.M{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Debug(value...)
	slog.WithExtra(slog.M{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Debug(value...)
}

func Fatal(value ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	slog.WithFields(slog.M{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Fatal(value...)
}

func Println(value ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	slog.WithFields(slog.M{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Println(value...)
}

func Infof(format string, msg ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	slog.WithFields(slog.M{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Infof(format, msg...)
}

func Warningf(format string, msg ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	slog.WithFields(slog.M{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Warnf(format, msg...)
}

func Errorf(format string, err ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	slog.WithFields(slog.M{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Errorf(format, err...)
}

func Debugf(format string, value ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	slog.WithFields(slog.M{
		"meta": fmt.Sprintf("%s:%d", srcFile, numLine),
	}).Debugf(format, value...)
}

func Fatalf(format string, value ...interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	slog.WithFields(slog.M{
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
