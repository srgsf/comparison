package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	freqFile  = "2015-2017-spoken-frequency.txt"
	transFile = "translation.txt"
)
const freqLimit = 798

var baseFolder = filepath.FromSlash("data/a")

func main() {
	start := time.Now()
	files, err := ioutil.ReadDir(baseFolder)
	if err != nil {
		log.Fatalf("Unable to read directory %s\r\n", baseFolder)
	}
	wordsChan := make(chan string)
	done := make(chan bool, 1)
	var words []string
	go func() {
		for w := range wordsChan {
			words = append(words, w)
		}
		done <- true
	}()

	var rLimit syscall.Rlimit
	if err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Fatal("Unable to read file limits")
	}
	guard := make(chan bool, rLimit.Cur)

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, f := range files {
		guard <- true
		go func(f fs.FileInfo) {
			defer wg.Done()
			defer func() { <-guard }()
			freq, err := os.Open(path.Join(baseFolder, f.Name(), freqFile))
			if err != nil {
				return
			}
			defer func() {
				_ = freq.Close()
			}()
			scanner := bufio.NewScanner(freq)
			if !scanner.Scan() {
				return
			}
			i, err := strconv.Atoi(scanner.Text())
			if err != nil || freqLimit > i {
				return
			}
			trans, err := os.Open(path.Join(baseFolder, f.Name(), transFile))
			if err != nil {
				return
			}
			defer func() {
				_ = trans.Close()
			}()
			scanner = bufio.NewScanner(trans)
			scanner.Split(bufio.ScanWords)
			if scanner.Scan() {
				wordsChan <- scanner.Text()
			}
		}(f)
	}

	wg.Wait()
	close(wordsChan)
	<-done
	close(done)
	fmt.Printf("selected %d words from %d, took %f seconds\r\n", len(words), len(files),
		time.Since(start).Seconds())
}
