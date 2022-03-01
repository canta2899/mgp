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
    "log"
    "os"
    "os/signal"
    "path/filepath"
    "regexp"
    "sync"
    "syscall"
    "github.com/fatih/color"
)

// One megabyte
const MEGABYTE int64 = 1048576

// Runes for emoji
const OK string = string('\u2713')
const KO string = string('\u00D7')

// Colors for printing
var green =  color.New(color.FgHiGreen).SprintFunc()
var red   =  color.New(color.FgRed).SprintFunc()
var cyan  =  color.New(color.FgCyan).SprintFunc()


// Checks paths matching the exclude pattern
func GetMatches(patterns []string) []string {
    matches := []string{}

    for _, pattern := range patterns {
        m, _ := filepath.Glob(pattern)
        matches = append(matches, m...)
    }

    return matches 
}

// Routine performed by each worker
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
            log.Printf("%v %v\n", green(OK), filepath)
        }
    }
}

// Process path and enqueues if valid for match checking
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

    if !isdir && (*info).Size() < MEGABYTE {
        q.Enqueue(pathname)
    }

    return nil
}

// Handles fatal errors
func handle(err error) {
    if err != nil {
        log.Fatal(red(err.Error()))
    }
}

// Handler for sigterm (ctrl + c from cli)
func setHandlers() {
    sigch := make(chan os.Signal)
    signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-sigch 
        log.Fatal(cyan("\nClosing..."))
    }()
}

func setupLogger() {
    log.SetFlags(0)
}

func main() {
    var wg sync.WaitGroup

    setupLogger()
    setHandlers()
    params, err := ParseArgs()
    handle(err)

    color.NoColor = (*params.nocolor)

    r, _ := regexp.Compile(*params.pattern)
    matches := GetMatches(*params.exclude)

    q := NewQueue()

    wg.Add(*params.workers)
    for i := 0; i < *params.workers; i++ {
        go handler(q, &wg, r)
    }

    err = filepath.Walk(*params.startpath,
        func(pathname string, info os.FileInfo, err error) error {

            // Checking permission and access errors
            if err != nil {
                log.Printf("%v %v\n", red(KO), pathname)
                if info.IsDir() {
                    return filepath.SkipDir
                } else {
                    return nil
                }
            }

            // Processes path in search of matches with the given
            // pattern or the folders that excluded folders
            return ProcessPath(&info, pathname, q, matches)
        })

    // Closes the queue in order to sync with goroutines
    q.Done()

    // Waits for goroutines to finish
    wg.Wait()
}

