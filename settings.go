package main

import (
	"fmt"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func parseArgs() *settings {
	quality := kingpin.Flag("quality", "Quality in units equivalent to libjpeg quality").Short('q').Default("84").Int()
	verbose := kingpin.Flag("verbose", "Verbose mode").Short('v').Bool()
	maxThreads := kingpin.Flag("threads", "Max concurrent threads").Short('t').Default("3").Int()

	source := kingpin.Arg("source", "Source directory").Default("./").String()
	output := kingpin.Arg("output", "Output directory").Default("./").String()
	log := kingpin.Arg("log", "Log directory, the log is used to prevent duplicate compressions").Default("").String()

	kingpin.Parse()

	if *quality < 84 {
		*quality = 84
	}

	if *log == "" {
		*log = *output
	}

	*log = strings.TrimSuffix(*log, "/") + "/"
	*output = strings.TrimSuffix(*output, "/") + "/"
	*source = strings.TrimSuffix(*source, "/") + "/"

	s := &settings{
		quality: *quality,
		verbose: *verbose,

		log:    *log,
		output: *output,
		source: *source,

		maxThreads: *maxThreads,
	}

	if s.verbose {
		fmt.Printf("Quality  =>  %d\n", s.quality)
		fmt.Printf("Source   =>  %s\n", s.source)
		fmt.Printf("Output   =>  %s\n", s.output)
		fmt.Printf("Log      =>  %s\n", s.log)
	}

	return s
}

type settings struct {
	quality int
	verbose bool

	log    string
	output string
	source string

	version    string
	maxThreads int
}
