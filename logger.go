package main

import (
	"fmt"
	"time"
)

type logger struct {
	pipe     chan interface{}
	logLevel string
}

func newLogger(logLevel string) *logger {
	l := &logger{
		pipe:     make(chan interface{}),
		logLevel: logLevel,
	}

	go func() {
		for s := range l.pipe {
			switch s.(type) {
			case error:
				fmt.Println(s)
			default:
				if l.logLevel != "" {
					fmt.Println(s)
				}
			}
		}
	}()

	return l
}

func (l *logger) Close() {
	close(l.pipe)
}

func (l *logger) log(s interface{}) {
	l.pipe <- s
}

type colorFunc func(string) string

var colors = []string{
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"red+h",
	"green+h",
	"yellow+h",
	"blue+h",
	"magenta+h",
	"cyan+h",
}

func logForJob(j *job) func(string) string {
	return func(s string) string {
		return j.color(fmt.Sprintf("%s %s : \n%s", time.Now().Format("15:04:05"), j.settings.source+j.fileName, s))
	}
}
