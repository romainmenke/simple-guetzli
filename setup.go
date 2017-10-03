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

	return s
}

func adjustSettingsBasedOnJobs(s *settings, numberOfJobs int) *settings {

	new := settings{}
	new = *s

	if new.maxThreads > numberOfJobs && numberOfJobs > 0 {
		s.maxThreads = numberOfJobs
	}

	originalMemLimit := new.memlimit
	adjustedMemLimit := new.memlimit / new.maxThreads
	s.memlimit = new.memlimit / new.maxThreads

	if new.verbose {
		fmt.Printf("Quality     =>  %d\n", new.quality)
		fmt.Printf("NoMemLimit  =>  %t\n", new.nomemlimit)
		fmt.Printf("MemLimit    =>  %d / %d => %d\n", originalMemLimit, new.maxThreads, adjustedMemLimit)
		fmt.Printf("Source      =>  %s\n", new.source)
		fmt.Printf("Output      =>  %s\n", new.output)
		fmt.Printf("Log         =>  %s\n", new.log)
		fmt.Printf("Force       =>  %t\n", new.force)
		fmt.Printf("Force Q     =>  %t\n", new.forceQuality)
		fmt.Printf("Threads     =>  %d\n", new.maxThreads)
		fmt.Printf("Dont Grow   =>  %t\n", new.dontGrow)
		fmt.Printf("Copy        =>  %t\n", new.copy)
		fmt.Printf("Version     =>  %s\n", new.version)
	}

	return &new

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
