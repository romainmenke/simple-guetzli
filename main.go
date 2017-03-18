package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func main() {

	source := flag.String("source", "./", "source directory")
	out := flag.String("out", "./", "output directory")
	log := flag.String("log", "./", "log directory, the log is used to prevent duplicate compressions")
	level := flag.Int("level", 84, "compression level")
	flag.Parse()

	sourceDir := strings.TrimSuffix(*source, "/") + "/"
	outDir := strings.TrimSuffix(*out, "/") + "/"
	logDir := strings.TrimSuffix(*log, "/") + "/"
	createIfMissing(outDir)

	exclude := flag.Args()

	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		panic(err)
	}

	logs := getLogs(logDir)
	logC := make(chan *guetzliLog, len(files))

	var wg sync.WaitGroup

FILE_ITERATOR:
	for _, f := range files {
		if !isFile(sourceDir + f.Name()) {
			continue FILE_ITERATOR
		}

		for _, exc := range exclude {
			if strings.Contains(f.Name(), exc) {
				continue FILE_ITERATOR
			}
		}

		j := job{
			sourceDir: sourceDir,
			outDir:    outDir,
			fileName:  f.Name(),
			level:     *level,
			log:       logs[sourceDir+f.Name()],
		}

		if j.level < 84 {
			j.level = 84
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			logC <- execute(&j)
		}()
	}

	wg.Wait()

	for _ = range files {
		l := <-logC
		if l != nil {
			logs[l.Path] = l
		}
	}

	saveLogs(logs, logDir)
}

type job struct {
	sourceDir string
	outDir    string
	fileName  string
	level     int
	log       *guetzliLog
}

func execute(j *job) *guetzliLog {
	modTime := timeModified(j.sourceDir + j.fileName)
	if j.log != nil && modTime.Equal(j.log.ModTime) {
		return nil
	}

	sha := sha1ForFile(j.sourceDir + j.fileName)
	if j.log != nil && sha == j.log.Sha1 {
		return nil
	}

	err := guetzli(j)
	if err != nil {
		panic(err)
	}

	return &guetzliLog{
		ModTime: modTime,
		Sha1:    sha,
		Path:    j.sourceDir + j.fileName,
	}
}

type guetzliLog struct {
	ModTime time.Time `json:"revisionTime"`
	Sha1    string    `json:"checksumSHA1"`
	Path    string    `json:"path"`
}

func writeFile(content []byte, fileName string, out string, level int) {

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

func createIfMissing(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
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

func saveLogs(logs map[string]*guetzliLog, path string) {

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
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnSHA1String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := sha1.New()

	//Copy the file in the hash interface and check for any error
	if _, err = io.Copy(hash, file); err != nil {
		panic(err)
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:20]

	//Convert the bytes to a string
	returnSHA1String = hex.EncodeToString(hashInBytes)

	return returnSHA1String

}

func guetzli(j *job) error {
	args := []string{
		"--quality",
		fmt.Sprintf("%d", j.level),
		j.sourceDir + j.fileName,
		j.outDir + j.fileName,
	}
	_, err := exec.Command("guetzli", args...).Output()
	if err != nil {
		return err
	}
	return nil
}
