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

	if s.maxThreads > numberOfJobs && numberOfJobs > 0 {
		s.maxThreads = numberOfJobs
	}

	originalMemLimit := s.memlimit
	adjustedMemLimit := s.memlimit / s.maxThreads
	s.memlimit = s.memlimit / s.maxThreads

	if s.verbose {
		fmt.Printf("Quality     =>  %d\n", s.quality)
		fmt.Printf("NoMemLimit  =>  %t\n", s.nomemlimit)
		fmt.Printf("MemLimit    =>  %d / %d => %d\n", originalMemLimit, s.maxThreads, adjustedMemLimit)
		fmt.Printf("Source      =>  %s\n", s.source)
		fmt.Printf("Output      =>  %s\n", s.output)
		fmt.Printf("Log         =>  %s\n", s.log)
		fmt.Printf("Force       =>  %t\n", s.force)
		fmt.Printf("Force Q     =>  %t\n", s.forceQuality)
		fmt.Printf("Threads     =>  %d\n", s.maxThreads)
		fmt.Printf("Version     =>  %s\n", s.version)
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
