package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mgutz/ansi"
)

func main() {

	var (
		settings   *settings
		logger     *logger
		reports    map[string]guetzliReport
		newReports map[string]guetzliReport
		jobs       []*job
	)

	settings = parseArgs()
	settings = preflight(settings)
	logger = newLogger(settings.verbose)
	reports = getReports(settings.log)
	jobs, newReports = getJobs(settings, reports, logger)

	if len(jobs) == 0 {
		return
	}

	settings = adjustSettingsBasedOnJobs(settings, len(jobs))

	if settings.force {
		saveReports(settings.version, newReports, settings.log)
	} else if settings.forceQuality {
		for path, r := range reports {
			if r.Quality == settings.quality {
				newReports[path] = r
			}
		}
		saveReports(settings.version, newReports, settings.log)
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
						newReports[j.report.Path] = j.report
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

	saveReports(settings.version, newReports, settings.log)
}

func getJobs(settings *settings, reports map[string]guetzliReport, logger *logger) ([]*job, map[string]guetzliReport) {

	newReports := make(map[string]guetzliReport)

	files, err := ioutil.ReadDir(settings.source)
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
			newReports[j.report.Path] = j.report
			fmt.Printf("%s %s : \n- skipped\n", time.Now().Format("15:04:05"), j.settings.source+j.fileName)
			continue FILE_ITERATOR
		}

		jobs = append(jobs, j)
	}

	return jobs, newReports
}
