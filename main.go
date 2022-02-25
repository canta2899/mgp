/**

    Multigrep, a tool for efficient recursive search
    of file content using regexps and goroutines

    Usage
    ======

    mgrep [pattern] [startpath]

    Flags
    ======
        • -w  (degree of parallelism [default 8])
        • -x  (directories to be excluded)

    Example
    =======

    mgrep -w 20 -x ./docs/w1 -x ./docs/w2 ^.+end$ ./docs/

**/

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

func GetMatches(patterns []string) []string {
    matches := []string{}

    for _, pattern := range patterns {
        m, _ := filepath.Glob(pattern)
        matches = append(matches, m...)
    }

    return matches 
}

func handler(q *Queue, wg *sync.WaitGroup, r *regexp.Regexp, control *SyncFlag) {
    defer wg.Done()

    for {

        filepath, err := q.Dequeue()

        if err != nil {

            if control.Get() {
                return
            }

            time.Sleep(100 * time.Nanosecond)
        }

        fi, err := os.Stat(filepath)

        if err != nil || !fi.Mode().IsRegular() {
            // Skips
            continue
        }

        filedata, err := os.ReadFile(filepath)

        if err != nil {
            // Skips
            continue
        }

        if r.Match(filedata) {
            fmt.Println(filepath)
        }
    }
}

func main() {
    var wg sync.WaitGroup

    // Allows synchronization in order to let the
    // main function wait for the other goroutines
    control := NewFlag(false)

    params, err := ParseArgs()
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    r, _ := regexp.Compile(*params.pattern)
    matches := GetMatches(*params.exclude)

    q := NewQueue()

    wg.Add(*params.workers)
    for i := 0; i < *params.workers; i++ {
        go handler(q, &wg, r, control)
    }

    err = filepath.Walk(*params.startpath,
        func(pathname string, info os.FileInfo, err error) error {

            if err != nil {
                return err
            }

            return Process(&info, pathname, q, matches)
        })

    if err != nil {
        panic(err.Error())
    }

    // States that there are no more path to be
    // enqueued for synchronization purposes
    control.Set(true)

    // Waits for goroutines to finish
    wg.Wait()
}

func Process(info *os.FileInfo, pathname string, q *Queue, excludes []string) error {
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
