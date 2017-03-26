package main

import (
	"fmt"
	"os"
)

func preflight(s *settings) *settings {

	createIfMissing(s.log)
	createIfMissing(s.output)

	version, err := guetzliVersion()
	if err != nil {
		panic(err)
	}

	s.version = version

	if s.verbose {
		fmt.Printf("Version  =>  %s\n", s.version)
	}

	return s

}

func createIfMissing(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func guetzliVersion() (string, error) {
	return "1.0.1", nil
}
