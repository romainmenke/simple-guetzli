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

	version, err := guetzliVersion()
	if err != nil {
		panic(err)
	}

	jobs := []*job{}

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

		jobs = append(jobs, j)
	}

	jobsQueue := make(chan *job, len(jobs))
	for _, j := range jobs {
		jobsQueue <- j
	}
	close(jobsQueue)

	var wg sync.WaitGroup
	var cancels []chan bool
	for i := 0; i < settings.maxThreads; i++ {
		wg.Add(1)

		go func() {
		JOB_QUEUE:
			for j := range jobsQueue {

				cancel := make(chan bool, 1)
				cancels = append(cancels, cancel)

				go func() {
					do(j)
				}()

				select {
				case success := <-j.done:
					if success {
						reports[j.report.Path] = j.report
					}
				case <-cancel:
					close(j.quit)
					break JOB_QUEUE
				}
			}
			wg.Done()
		}()
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		_ = <-signalChannel
		for _, c := range cancels {
			close(c)
		}
	}()

	wg.Wait()

	saveReports(version, reports, settings.log)
}
