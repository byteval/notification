package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	With(fields ...interface{}) Logger
}

type appLogger struct {
	*log.Logger
	fields []interface{}
}

func New(out io.Writer) Logger {
	if out == nil {
		out = os.Stderr
	}
	return &appLogger{
		Logger: log.New(out, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l *appLogger) formatFields(fields []interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	var str string
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			str += fmt.Sprintf(" %s=%v", fields[i], fields[i+1])
		} else {
			str += fmt.Sprintf(" %s=", fields[i])
		}
	}
	return str
}

func (l *appLogger) log(level, msg string, fields []interface{}) {
	allFields := append(l.fields, fields...)
	l.Printf("[%s] %s%s", level, msg, l.formatFields(allFields))
}

func (l *appLogger) Info(msg string, fields ...interface{}) {
	l.log("INFO", msg, fields)
}

func (l *appLogger) Error(msg string, fields ...interface{}) {
	l.log("ERROR", msg, fields)
}

func (l *appLogger) Fatal(msg string, fields ...interface{}) {
	l.log("FATAL", msg, fields)
	os.Exit(1)
}

func (l *appLogger) With(fields ...interface{}) Logger {
	return &appLogger{
		Logger: l.Logger,
		fields: append(l.fields, fields...),
	}
}
