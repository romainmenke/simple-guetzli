package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mgutz/ansi"
)

func main() {

	settings := parseArgs()
	settings = preflight(settings)

	logger := newLogger(settings.verbose)

	files, err := ioutil.ReadDir(settings.source)
	if err != nil {
		panic(err)
	}

	reports := getReports(settings.log)
	var jobCounter int

	version, err := guetzliVersion()
	if err != nil {
		panic(err)
	}

	jobs := []*job{}

	var wg sync.WaitGroup

FILE_ITERATOR:
	for index, f := range files {
		if !isFile(settings.source + f.Name()) {
			continue FILE_ITERATOR
		}

		j := &job{
			fileName: f.Name(),
			report:   reports[settings.source+f.Name()],
			settings: settings,
			quit:     make(chan bool, 1),
			done:     make(chan bool, 1),
			logger:   logger,
		}

		j.color = ansi.ColorFunc(colors[index%len(colors)])

		if !needsProc(j) {
			reports[j.report.Path] = j.report
			j.logger.log(logForJob(j)("- skipped"))
			continue FILE_ITERATOR
		}

		wg.Add(1)
		jobCounter++

		jobs = append(jobs, j)

		go func() {
			work(j)
		}()

		go func() {
			select {
			case success := <-j.done:
				if success {
					reports[j.report.Path] = j.report
				}
				wg.Done()
				break
			}
		}()

	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		_ = <-signalChannel
		for _, j := range jobs {
			j.quit <- true
		}
	}()

	wg.Wait()

	saveReports(version, reports, settings.log)
}
