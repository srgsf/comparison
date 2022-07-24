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
	"time"
)

const (
	freqFile  = "2015-2017-spoken-frequency.txt"
	transFile = "translation.txt"
)
const freqLimit = 798

//var baseFolder = filepath.FromSlash("./data/a")

var baseFolder = filepath.FromSlash("../data/a")

func main() {
	start := time.Now()
	var words []string
	var procs int = 0
	files, err := ioutil.ReadDir(baseFolder)
	if err != nil {
		log.Fatalf("Unable to read directory %s\r\n", baseFolder)
	}
	wordsChan := make(chan string)
	var wg sync.WaitGroup
	for _, f := range files {
		wg.Add(1)
		procs += 1
		go func(f fs.FileInfo) {
			defer wg.Done()

			freq, err := os.Open(path.Join(baseFolder, f.Name(), freqFile))
			if err != nil {
				return
			}
			scanner := bufio.NewScanner(freq)
			if !scanner.Scan() {
				//return
			}
			i, err := strconv.Atoi(scanner.Text())
			if err != nil || freqLimit > i {
				return
			}
			trans, err := os.Open(path.Join(baseFolder, f.Name(), transFile))
			if err != nil {
				return
			}
			scanner = bufio.NewScanner(trans)
			scanner.Split(bufio.ScanWords)
			word := ""
			if scanner.Scan() {
				word = scanner.Text()
			} else {
				word = f.Name()
			}

			//word := f.Name()
			//fmt.Printf("%s\n", word)
			wordsChan <- word
		}(f)
	}
	done := make(chan struct{})
	go func(channel chan string, done chan struct{}) {
		for data := range channel {
			//fmt.Printf("%s\n", data)
			words = append(words, data)
		}
		close(done)
	}(wordsChan, done)
	wg.Wait()
	close(wordsChan)
	<-done

	fmt.Printf("selected %d words from %d, took %f seconds\r\n", len(words), len(files),
		time.Since(start).Seconds())

	/*
		sort.Strings(words)
		i := 0
		for _, word := range words {
			i++
			fmt.Printf("%d. %s\n", i, word)
		}
	*/
}
