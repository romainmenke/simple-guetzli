package main

import (
	"fmt"
	"time"
)

type logger struct {
	pipe    chan string
	verbose bool
}

func newLogger(verbose bool) *logger {
	l := &logger{
		pipe:    make(chan string),
		verbose: verbose,
	}

	go func() {
		for s := range l.pipe {
			if l.verbose {
				fmt.Print(s)
			}
		}
	}()

	return l
}

func (l *logger) Close() {
	close(l.pipe)
}

func (l *logger) log(s string) {
	l.pipe <- s
}

type colorFunc func(string) string

var colors = []string{
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
}

func logForJob(j *job) func(string) string {
	return func(s string) string {
		return j.color(fmt.Sprintf("%s %s : \n%s", time.Now().Format("15:04:05"), j.settings.source+j.fileName, s))
	}
}