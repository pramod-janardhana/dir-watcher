package zlog

import (
	"log"
	"sync"
)

var logger Logger

type Logger interface {
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

type Closer interface {
	Close()
}

var once *sync.Once = &sync.Once{}

func InitLogger(l Logger) {
	once.Do(func() {
		logger = l
	})
}

func CloseLogger() {
	if closable, ok := logger.(Closer); ok {
		closable.Close()
	}
}

func handlePanicf(format string, v ...interface{}) {
	if a := recover(); a != nil {
		// ignore any panic errors in logger package and log to default log.
		log.Printf(format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	defer handlePanicf(format, v...)
	logger.Infof(format, v...)
}

func Errorf(format string, v ...interface{}) {
	defer handlePanicf(format, v...)
	logger.Errorf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	defer handlePanicf(format, v...)
	logger.Debugf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	defer handlePanicf(format, v...)
	logger.Fatalf(format, v...)
}
