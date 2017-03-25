package main

import (
	"fmt"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func parseArgs() *args {
	quality := kingpin.Flag("quality", "Quality in units equivalent to libjpeg quality").Short('q').Default("84").Int()
	verbose := kingpin.Flag("verbose", "Verbose mode").Short('v').Bool()

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

	a := args{
		quality: *quality,
		verbose: *verbose,

		log:    *log,
		output: *output,
		source: *source,
	}

	if a.verbose {
		fmt.Printf("Quality  =>  %d\n", a.quality)
		fmt.Printf("Source   =>  %s\n", a.source)
		fmt.Printf("Output   =>  %s\n", a.output)
		fmt.Printf("Log      =>  %s\n", a.log)
	}

	return &a
}

type args struct {
	quality int
	verbose bool

	log    string
	output string
	source string

	version string
}
