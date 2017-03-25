package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type guetzliLog struct {
	Quality int       `json:"quality"`
	ModTime time.Time `json:"revisionTime"`
	Path    string    `json:"path"`
	Sha1    string    `json:"checksumSHA1"`
	Version string    `json:"guetzli-version"`
}

func getLogs(path string) map[string]*guetzliLog {
	buf := bytes.NewBuffer(nil)
	file, err := os.Open(path + "guetzli.json")
	if err != nil {
		return make(map[string]*guetzliLog)
	}

	_, err = io.Copy(buf, file)
	if err != nil {
		return make(map[string]*guetzliLog)
	}

	file.Close()
	logs := make(map[string]*guetzliLog)

	err = json.Unmarshal(buf.Bytes(), &logs)
	if err != nil {
		return make(map[string]*guetzliLog)
	}

	return logs
}

func saveLogs(version string, logs map[string]*guetzliLog, path string) {
	b, err := json.Marshal(logs)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path+"guetzli.json", out.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func timeModified(filePath string) time.Time {
	info, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	return info.ModTime()
}

func sha1ForFile(filePath string) string {
	var returnSHA1String string
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	hash := sha1.New()
	if _, err = io.Copy(hash, file); err != nil {
		panic(err)
	}

	hashInBytes := hash.Sum(nil)[:20]
	returnSHA1String = hex.EncodeToString(hashInBytes)

	return returnSHA1String
}

func guetzli(j *job) error {
	args := []string{
		"--quality",
		fmt.Sprintf("%d", j.args.quality),
	}

	if j.args.verbose {
		args = append(args, "--verbose")
	}

	args = append(args, j.args.source+j.fileName)
	args = append(args, j.args.output+j.fileName)

	_, err := exec.Command("guetzli", args...).Output()
	if err != nil {
		return err
	}
	return nil
}
