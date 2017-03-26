package main

import (
	"fmt"
	"os"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func parseArgs() *settings {
	quality := kingpin.Flag("quality", "Quality in units equivalent to libjpeg quality").Short('q').Default("95").Int()
	verbose := kingpin.Flag("verbose", "Verbose mode").Bool()
	force := kingpin.Flag("force", "Force recompression").Short('f').Bool()
	forceQuality := kingpin.Flag("force-quality", "Force recompression if quality changed").Bool()
	maxThreads := kingpin.Flag("threads", "Max concurrent threads").Short('t').Default("3").Int()

	source := kingpin.Arg("source", "Source directory").Default("./").String()
	output := kingpin.Arg("output", "Output directory").Default("./").String()
	log := kingpin.Arg("log", "Log directory, the log is used to prevent duplicate compressions").Default("").String()

	version := kingpin.Flag("version", "Guetzli Version").Short('v').Bool()

	kingpin.Parse()

	if *version {
		v, err := guetzliVersion()
		if err != nil {
			panic(err)
		}
		fmt.Println("Using Guetzli : " + v)
		os.Exit(0)
	}

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
		quality:      *quality,
		force:        *force,
		forceQuality: *forceQuality,
		verbose:      *verbose,

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
		fmt.Printf("Force    =>  %t\n", s.force)
		fmt.Printf("Force Q  =>  %t\n", s.forceQuality)
		fmt.Printf("Threads  =>  %d\n", s.maxThreads)
	}

	return s
}

type settings struct {
	quality      int
	force        bool
	forceQuality bool
	verbose      bool

	log    string
	output string
	source string

	version    string
	maxThreads int
}
