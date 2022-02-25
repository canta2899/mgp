/**
    Multigrep
    =========

    A CLI tool that allows faster recursive search of files
    whose content matches the given pattern.

    It returns the equivalent of grep -E -r -l "pattern" "path".

    Run multigrep --help for a brief usage guide
**/

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

func GetMatches(patterns []string) []string {
    matches := []string{}

    for _, pattern := range patterns {
        m, _ := filepath.Glob(pattern)
        matches = append(matches, m...)
    }

    return matches 
}

func handler(q *Queue, wg *sync.WaitGroup, r *regexp.Regexp) {
    defer wg.Done()

    for {

        filepath, err := q.Dequeue()

        if err != nil {
            return 
        }

        fi, err := os.Stat(filepath)

        if err != nil || !fi.Mode().IsRegular() {
            continue // Skips
        }

        filedata, err := os.ReadFile(filepath)

        if err != nil {
            continue // Skips
        }

        if r.Match(filedata) {
            fmt.Println(filepath)
        }
    }
}

func ProcessPath(info *os.FileInfo, pathname string, q *Queue, excludes []string) error {
    isdir := (*info).IsDir()

    for _, n := range excludes {
        m, err := filepath.Match(n, pathname)

        if err != nil {
            return err
        }

        if isdir && m {
           return filepath.SkipDir 
        }
    }

    if !isdir {
        q.Enqueue(pathname)
    }

    return nil
}

func handle(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    var wg sync.WaitGroup

    params, err := ParseArgs()

    handle(err)

    r, _ := regexp.Compile(*params.pattern)
    matches := GetMatches(*params.exclude)

    q := NewQueue()

    wg.Add(*params.workers)
    for i := 0; i < *params.workers; i++ {
        go handler(q, &wg, r)
    }

    err = filepath.Walk(*params.startpath,
        func(pathname string, info os.FileInfo, err error) error {

            if err != nil {
                return err
            }

            // Processes path in search of matches with the given
            // pattern or the folders that excluded folders
            return ProcessPath(&info, pathname, q, matches)
        })

    // Panic if issues occour while traversing path
    handle(err)

    // Closes the queue in order to sync with goroutines
    q.Done()

    // Waits for goroutines to finish
    wg.Wait()
}

