package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Info(msg string)
	Error(msg string, err error)
}

type StdLogger struct {
	info  *log.Logger
	error *log.Logger
}

func New() Logger {
	return &StdLogger{
		info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lmsgprefix),
		error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lmsgprefix),
	}
}

func (l *StdLogger) Info(msg string) {
	l.info.Println(msg)
}

func (l *StdLogger) Error(msg string, err error) {
	l.error.Println(fmt.Sprintf("%s: %v", msg, err))
}
