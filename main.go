package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
)

func handler(c <-chan interface{}, wg *sync.WaitGroup, r *regexp.Regexp) {
	defer wg.Done()

	for {

		if len(c) == 0 {
			return
		}

		next := <-c

		filepath := next.(string) // type assertion

		fi, err := os.Stat(filepath)

		if err != nil || !fi.Mode().IsRegular() {
			continue
		}

		filedata, err := os.ReadFile(filepath)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if r.Match(filedata) {
			fmt.Println(next)
		}
	}
}

func main() {

	// Allows synchronization in order to let the
	// main function wait for the other coroutines
	var wg sync.WaitGroup

	startpath := os.Args[1]
	pattern := os.Args[2]
	workers, _ := strconv.Atoi(os.Args[3])

	r, _ := regexp.Compile(pattern)

	files := []string{}

	err := filepath.Walk(startpath,
		func(path string, info os.FileInfo, err error) error {

			// Should exlude some paths if specified
			// if isin(path, x) {
			// 	return filepath.SkipDir
			// }

			if err != nil {
				return err
			}

			files = append(files, path)

			return nil
		})

	if err != nil {
		panic(err.Error())
	}

	// Using a channel in order to implement a thread
	// safe queue that allows multiple consumers to
	// "dequeue" each path and process the content of
	// the referred file
	c := make(chan interface{}, len(files))

	for _, p := range files {
		c <- p
	}

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go handler(c, &wg, r)
	}

	wg.Wait()
}

// Utility that checks whether el is in list
func isin(el string, list []string) bool {
	for _, entry := range list {
		if el == entry {
			return true
		}
	}
	return false
}
