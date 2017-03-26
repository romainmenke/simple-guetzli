package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type job struct {
	fileName string
	report   *guetzliReport
	settings *settings
	color    colorFunc
	logger   *logger
	quit     chan bool
	done     chan bool
}

func work(j *job) {
	var (
		outb bytes.Buffer
		errb bytes.Buffer
	)

	args := []string{
		"--quality",
		fmt.Sprintf("%d", j.settings.quality),
	}

	if j.settings.verbose {
		args = append(args, "--verbose")
	}

	args = append(args, j.settings.source+j.fileName)
	args = append(args, j.settings.output+j.fileName)

	cmd := exec.Command("guetzli", args...)

	cmd.Stdout = &outb
	cmd.Stderr = &errb

	j.logger.log(logForJob(j)("start\n"))
	cmd.Start()

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-j.quit:
		err := cmd.Process.Kill()
		if err != nil {
			j.logger.log(logForJob(j)(err.Error()))
		}

		j.logger.log(logForJob(j)("cancelled\n"))
		j.done <- false
		break
	case err := <-done:
		if err != nil {
			j.logger.log(logForJob(j)(errb.String()))
			j.done <- false
		} else {
			j.logger.log(logForJob(j)(outb.String()))
			j.done <- true
		}
		break
	}
}
