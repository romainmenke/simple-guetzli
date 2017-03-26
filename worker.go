package main

import (
	"bytes"
	"errors"
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

func do(j *job) {
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

	j.logger.log(logForJob(j)("- start"))
	err := cmd.Start()
	if err != nil {
		j.logger.log(errors.New(logForJob(j)(err.Error())))
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
		j.logger.log(logForJob(j)("- cancelled"))
		j.done <- false
		close(j.done)
		break
	case err := <-done:
		if err != nil {
			j.logger.log(errors.New(logForJob(j)(errb.String())))
			j.done <- false
			close(j.done)
		} else {
			j.logger.log(logForJob(j)(outb.String()))
			j.done <- true
			close(j.done)
		}
		break
	}
}
