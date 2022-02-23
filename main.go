package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
)

func handler(c <-chan string, wg *sync.WaitGroup, r *regexp.Regexp) {
	defer wg.Done()

	for {

		if len(c) == 0 {
			return
		}

		next := <-c

		fi, err := os.Stat(next)

		if err != nil || !fi.Mode().IsRegular() {
			continue
		}

		filedata, err := os.ReadFile(next)

		if err != nil {
			fmt.Println(err.Error())
			continue
		} else {
			if r.Match(filedata) {
				fmt.Println(next)
			}
		}
	}
}

func main() {

	var wg sync.WaitGroup

	startpath := os.Args[1]
	pattern := os.Args[2]
	workers, _ := strconv.Atoi(os.Args[3])

	r, _ := regexp.Compile(pattern)

	// x := os.Args[4:]
	files := []string{}

	err := filepath.Walk(startpath,
		func(path string, info os.FileInfo, err error) error {

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

	c := make(chan string, len(files)+1)

	for _, p := range files {
		c <- p
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go handler(c, &wg, r)
	}
	close(c)

	wg.Wait()
}

func isin(e string, l []string) bool {
	for _, el := range l {
		if el == e {
			return true
		}
	}
	return false
}
