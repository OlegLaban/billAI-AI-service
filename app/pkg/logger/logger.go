package logger

import "fmt"

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) Error(msg string, err error) {
	fmt.Println(msg, err.Error())
}
