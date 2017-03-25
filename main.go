package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/mgutz/ansi"
)

func main() {

	a := parseArgs()
	a = preflight(a)

	files, err := ioutil.ReadDir(a.source)
	if err != nil {
		panic(err)
	}

	logs := getLogs(a.log)
	logC := make(chan *guetzliLog)
	var jobCounter int

	version, err := guetzliVersion()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

FILE_ITERATOR:
	for index, f := range files {
		if !isFile(a.source + f.Name()) {
			continue FILE_ITERATOR
		}

		j := job{
			fileName: f.Name(),
			log:      logs[a.source+f.Name()],
			args:     a,
		}

		if a.verbose {
			j.color = ansi.ColorFunc(colors[index%len(colors)])
		}

		wg.Add(1)
		jobCounter++

		go func() {
			defer wg.Done()
			logC <- execute(&j)
		}()
	}

	for i := 0; i < jobCounter; i++ {
		l := <-logC
		if l != nil {
			logs[l.Path] = l
		}
	}

	wg.Wait()

	saveLogs(version, logs, a.log)
}

type job struct {
	fileName string
	log      *guetzliLog
	args     *args
	color    ColorFunc
}

func execute(j *job) *guetzliLog {

	var imgChanged bool
	var settingsChanged bool

	if j.log == nil {
		settingsChanged = true
		imgChanged = true
	}

	if j.log != nil && j.log.Version != j.args.version {
		settingsChanged = true
	}

	if j.log != nil && j.log.Quality != j.args.quality {
		settingsChanged = true
	}

	modTime := timeModified(j.args.source + j.fileName)
	if j.log != nil && modTime.Equal(j.log.ModTime) {
		imgChanged = true
	}

	sha := sha1ForFile(j.args.source + j.fileName)
	if j.log != nil && sha != j.log.Sha1 {
		imgChanged = true
	}

	if !imgChanged && !settingsChanged {
		return nil
	}

	if j.args.verbose {
		fmt.Println(j.color(fmt.Sprintf("Processing  =>  %s", j.args.source+j.fileName)))
	}

	err := guetzli(j)
	if err != nil {
		panic(err)
	}

	if j.args.verbose {
		fmt.Println(j.color(fmt.Sprintf("Done        =>  %s", j.args.source+j.fileName)))
	}

	return &guetzliLog{
		Quality: j.args.quality,
		ModTime: modTime,
		Path:    j.args.source + j.fileName,
		Sha1:    sha,
		Version: j.args.version,
	}
}

func writeFile(content []byte, fileName string, out string, quality int) {
	err := ioutil.WriteFile(out+fileName, content, 0644)
	if err != nil {
		panic(err)
	}
}

func isFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return !fileInfo.IsDir()
}

type ColorFunc func(string) string

var colors = []string{
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
}
