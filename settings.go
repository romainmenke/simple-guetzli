package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var defaultJPEGQuality = 95
var defaultMemlimitMB = 6000
var defaultMaxThreads = 3

func parseArgs() *settings {
	quality := kingpin.Flag("quality", fmt.Sprintf("Visual quality to aim for, expressed as a JPEG quality value. Default value is %d.", defaultJPEGQuality)).
		Short('q').
		Default(fmt.Sprint(defaultJPEGQuality)).
		Uint()

	verbose := kingpin.Flag("verbose", "Print a verbose trace of all attempts to standard output.").
		Bool()

	memlimit := kingpin.Flag("memlimit", fmt.Sprintf("Memory limit in MB. Guetzli will fail if unable to stay under the limit. Default limit is %d", defaultMemlimitMB)).
		Short('m').
		Default(fmt.Sprint(defaultMemlimitMB)).
		Uint()

	nomemlimit := kingpin.Flag("nomemlimit", "Do not limit memory usage.").
		Bool()

	force := kingpin.Flag("force", "Force recompression").
		Short('f').
		Bool()

	forceQuality := kingpin.Flag("force-quality", "Force recompression if quality changed").
		Bool()

	maxThreads := kingpin.Flag("threads", fmt.Sprintf("Max concurrent threads. Default limit is %d", maxParallelism()-1)).
		Short('t').
		Default(fmt.Sprint(defaultMaxThreads)).
		Uint()

	dontGrow := kingpin.Flag("dontgrow", "Delete compressed files that got bigger").
		Short('g').
		Bool()

	copy := kingpin.Flag("copy", "Copy all files from source folder to output, without overwriting compression results").
		Short('c').
		Bool()

	source := kingpin.Arg("source", "Source directory").
		Default("./").
		String()

	output := kingpin.Arg("output", "Output directory").
		Default("./").
		String()

	log := kingpin.Arg("log", "Log directory, the log is used to prevent duplicate compressions").
		Default("").
		String()

	interval := kingpin.Flag("interval", "").
		Short('i').
		Hidden().
		Int()

	version := kingpin.Flag("version", "Guetzli Version").Short('v').Bool()

	_ = memlimit
	_ = nomemlimit

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
		quality: int(*quality),
		verbose: *verbose,

		memlimit:   int(*memlimit),
		nomemlimit: *nomemlimit,

		force:        *force,
		forceQuality: *forceQuality,

		log:    *log,
		output: *output,
		source: *source,

		maxThreads: int(*maxThreads),

		dontGrow: *dontGrow,
		copy:     *copy,

		interval: *interval,
	}

	return s
}

type settings struct {
	quality int
	verbose bool

	memlimit   int
	nomemlimit bool

	force        bool
	forceQuality bool

	log    string
	output string
	source string

	version    string
	maxThreads int

	dontGrow bool
	copy     bool

	interval int
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
