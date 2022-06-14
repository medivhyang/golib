package log

import (
	"fmt"
	"time"
)

type Logger struct {
	Level     Level
	Prefix    string
	Locale    *time.Location
	Appenders []Appender
}

func New(level Level, prefix string, appenders ...Appender) *Logger {
	l := &Logger{
		Prefix:    prefix,
		Level:     level,
		Locale:    time.Local,
		Appenders: appenders,
	}
	return l
}

func (l *Logger) New(module ...string) *ModuleLogger {
	return &ModuleLogger{root: l, module: module}
}

func (l *Logger) Append(level Level, message string, data ...map[string]interface{}) {
	if !l.Level.Valid() || level < l.Level {
		return
	}
	t := Template{
		Prefix:  l.Prefix,
		Level:   level,
		Message: message,
		Time:    time.Now(),
	}
	if l.Locale != nil {
		t.Time = t.Time.In(l.Locale)
	}
	if len(data) > 0 {
		t.Data = data[0]
	}
	for _, appender := range l.Appenders {
		_ = appender.Append(t)
	}
}

func (l *Logger) Appendf(level Level, format string, args ...interface{}) {
	l.Append(level, fmt.Sprintf(format, args...))
}

func (l *Logger) Debug(message string, data ...map[string]interface{}) {
	l.Append(LevelDebug, message, data...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Appendf(LevelDebug, format, args...)
}

func (l *Logger) Info(message string, data ...map[string]interface{}) {
	l.Append(LevelInfo, message, data...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Appendf(LevelInfo, format, args...)
}

func (l *Logger) Warn(message string, data ...map[string]interface{}) {
	l.Append(LevelWarn, message, data...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Appendf(LevelWarn, format, args...)
}

func (l *Logger) Error(message string, data ...map[string]interface{}) {
	l.Append(LevelError, message, data...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Appendf(LevelError, format, args...)
}

func (l *Logger) Fatal(message string, data ...map[string]interface{}) {
	l.Append(LevelFatal, message, data...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Appendf(LevelFatal, format, args...)
}

type ModuleLogger struct {
	root   *Logger
	module []string
}

func (l *ModuleLogger) New(module ...string) *ModuleLogger {
	result := &ModuleLogger{root: l.root}
	result.module = append([]string{}, l.module...)
	result.module = append(result.module, module...)
	return result
}

func (l *ModuleLogger) Append(level Level, message string, data ...map[string]interface{}) {
	if !l.root.Level.Valid() || level < l.root.Level {
		return
	}
	e := Template{
		Prefix:  l.root.Prefix,
		Modules: l.module,
		Level:   level,
		Message: message,
		Time:    time.Now(),
	}
	if l.root.Locale != nil {
		e.Time = e.Time.In(l.root.Locale)
	}
	if len(data) > 0 {
		e.Data = data[0]
	}
	for _, appender := range l.root.Appenders {
		_ = appender.Append(e)
	}
	return
}

func (l *ModuleLogger) Appendf(level Level, format string, args ...interface{}) {
	l.Append(level, fmt.Sprintf(format, args...))
}

func (l *ModuleLogger) Debug(message string, data ...map[string]interface{}) {
	l.Append(LevelDebug, message, data...)
}

func (l *ModuleLogger) Debugf(format string, args ...interface{}) {
	l.Appendf(LevelDebug, format, args...)
}

func (l *ModuleLogger) Info(message string, data ...map[string]interface{}) {
	l.Append(LevelInfo, message, data...)
}

func (l *ModuleLogger) Infof(format string, args ...interface{}) {
	l.Appendf(LevelInfo, format, args...)
}

func (l *ModuleLogger) Warn(message string, data ...map[string]interface{}) {
	l.Append(LevelWarn, message, data...)
}

func (l *ModuleLogger) Warnf(format string, args ...interface{}) {
	l.Appendf(LevelWarn, format, args...)
}

func (l *ModuleLogger) Error(message string, data ...map[string]interface{}) {
	l.Append(LevelError, message, data...)
}

func (l *ModuleLogger) Errorf(format string, args ...interface{}) {
	l.Appendf(LevelError, format, args...)
}

func (l *ModuleLogger) Fatal(message string, data ...map[string]interface{}) {
	l.Append(LevelFatal, message, data...)
}

func (l *ModuleLogger) Fatalf(format string, args ...interface{}) {
	l.Appendf(LevelFatal, format, args...)
}
