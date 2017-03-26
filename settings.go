package main

import (
	"fmt"
	"os"
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

	maxThreads := kingpin.Flag("threads", fmt.Sprintf("Max concurrent threads. Default limit is %d", defaultMaxThreads)).
		Short('t').
		Default(fmt.Sprint(defaultMaxThreads)).
		Uint()

	source := kingpin.Arg("source", "Source directory").
		Default("./").
		String()

	output := kingpin.Arg("output", "Output directory").
		Default("./").
		String()

	log := kingpin.Arg("log", "Log directory, the log is used to prevent duplicate compressions").
		Default("").
		String()

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
	}

	if s.verbose {
		fmt.Printf("Quality     =>  %d\n", s.quality)
		fmt.Printf("NoMemLimit  =>  %t\n", s.nomemlimit)
		fmt.Printf("MemLimit    =>  %d\n", s.memlimit)
		fmt.Printf("Source      =>  %s\n", s.source)
		fmt.Printf("Output      =>  %s\n", s.output)
		fmt.Printf("Log         =>  %s\n", s.log)
		fmt.Printf("Force       =>  %t\n", s.force)
		fmt.Printf("Force Q     =>  %t\n", s.forceQuality)
		fmt.Printf("Threads     =>  %d\n", s.maxThreads)
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
}
