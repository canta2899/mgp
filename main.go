package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func handler(q *Queue, wid int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		next, err := q.Dequeue()

		if err != nil {
			return
		}

		fmt.Println("Processing path:", next)
	}
}

func main() {
	var wg sync.WaitGroup

	startpath := os.Args[1]
	workers, _ := strconv.Atoi(os.Args[2])
	x := os.Args[3:]
	files := []string{}

	err := filepath.Walk(startpath,
		func(path string, info os.FileInfo, err error) error {

			if isin(path, x) {
				return filepath.SkipDir
			}

			if err != nil {
				return err
			}

			files = append(files, path)

			return nil
		})

	if err != nil {
		panic(err.Error())
	}

	q := NewQueue(len(files))

	for _, p := range files {
		q.Enqueue(p)
	}

	StartWorkers(q, &wg, workers)
	wg.Wait()
}

func StartWorkers(q *Queue, wg *sync.WaitGroup, count int) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go handler(q, i, wg)
	}
}

func isin(e string, l []string) bool {
	for _, el := range l {
		if el == e {
			return true
		}
	}
	return false
}
