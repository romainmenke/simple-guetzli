package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type job struct {
	fileName string
	report   guetzliReport
	settings *settings
	color    colorFunc
	logger   *logger
	errored  bool
	quit     chan bool
	done     chan bool
}

func do(j *job) {
	var (
		outb bytes.Buffer
		errb bytes.Buffer
	)

	args := []string{
		"--quality",
		fmt.Sprint(j.settings.quality),
		"--memlimit",
		fmt.Sprint(j.settings.memlimit),
	}

	if j.settings.nomemlimit {
		args = append(args, "--nomemlimit")
	}
	if j.settings.logLevel == "debug" {
		args = append(args, "--verbose")
	}

	outputFileName := j.fileName
	outputFileName = strings.TrimSuffix(strings.TrimSuffix(outputFileName, filepath.Ext(outputFileName)), ".") + ".jpg"

	args = append(args, j.settings.source+j.fileName)
	args = append(args, j.settings.output+outputFileName)

	j.fileName = outputFileName

	cmd := exec.Command("guetzli", args...)

	cmd.Stdout = &outb
	cmd.Stderr = &errb

	j.logger.log(logForJob(j)("start"))
	err := cmd.Start()
	if err != nil {
		j.logger.log(errors.New(logForJob(j)(err.Error())))
		j.errored = true
		j.done <- false
		close(j.done)
		return
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
		close(done)
	}()

	select {
	case <-j.quit:
		err := cmd.Process.Kill()
		if err != nil {
			j.logger.log(errors.New(logForJob(j)(err.Error())))
		}
		j.logger.log(logForJob(j)("cancelled"))
		j.done <- false
		close(j.done)
		break
	case err := <-done:
		if err != nil {
			j.errored = true
			j.logger.log(errors.New(logForJob(j)(errb.String())))
			j.done <- false
			close(j.done)
		} else {
			doneMsg := outb.String()
			if doneMsg == "" {
				doneMsg = "done"
			}
			j.logger.log(logForJob(j)(doneMsg))
			j.done <- true
			close(j.done)
		}
		break
	}
}
