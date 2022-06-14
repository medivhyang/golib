package log

import (
	"time"
)

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelText = map[Level]string{
	LevelDebug: "debug",
	LevelInfo:  "info",
	LevelWarn:  "warn",
	LevelError: "error",
	LevelFatal: "fatal",
}

func LevelText(l Level) string {
	return levelText[l]
}

type Level int

func (l Level) Valid() bool {
	switch l {
	case LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal:
		return true
	default:
		return false
	}
}

type Template struct {
	Prefix  string
	Modules []string
	Level   Level
	Message string
	Data    map[string]interface{}
	Time    time.Time
}

type Appender interface {
	Append(t Template) error
}

type Formatter interface {
	Format(t Template) ([]byte, error)
}

type Filter func(Template) bool

type Fields = map[string]interface{}

var Default = New(LevelDebug, "default", NewConsoleAppender(NewTextFormatter()))

func Debug(message string, data ...map[string]interface{}) {
	Default.Debug(message, data...)
}

func Debugf(format string, args ...interface{}) {
	Default.Debugf(format, args...)
}

func Info(message string, data ...map[string]interface{}) {
	Default.Info(message, data...)
}

func Infof(format string, args ...interface{}) {
	Default.Infof(format, args...)
}

func Warn(message string, data ...map[string]interface{}) {
	Default.Warn(message, data...)
}

func Warnf(format string, args ...interface{}) {
	Default.Warnf(format, args...)
}

func Error(message string, data ...map[string]interface{}) {
	Default.Error(message, data...)
}

func Errorf(format string, args ...interface{}) {
	Default.Errorf(format, args...)
}

func Fatal(message string, data ...map[string]interface{}) {
	Default.Fatal(message, data...)
}

func Fatalf(format string, args ...interface{}) {
	Default.Fatalf(format, args...)
}
