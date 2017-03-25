package main

import (
	"fmt"
	"os"
)

func preflight(a *args) *args {

	createIfMissing(a.log)
	createIfMissing(a.output)

	version, err := guetzliVersion()
	if err != nil {
		panic(err)
	}

	a.version = version

	if a.verbose {
		fmt.Printf("Version  =>  %s\n", a.version)
	}

	return a

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
