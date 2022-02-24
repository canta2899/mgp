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


func isin(el string, list []string) bool {

    // Utility that checks if file/dir should be skipped

    for _, entry := range list {
        m, _ := filepath.Match(filepath.Clean(entry), filepath.Clean(el))
        n, _ := filepath.Glob(el + "/"+ entry)
        if m || len(n) > 0 {
            return true
        }
    }

    return false
}

func PrintUsageGuide() {
    msg := `Usage: mgrep [pattern] [path]

    You can use the flags
        -w for workers
        -x to exclude paths`
    fmt.Println(msg)
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

    // Allows synchronization in order to let the
    // main function wait for the other goroutines
    control := NewFlag(false)

    var wg sync.WaitGroup

    params, err := ParseArgs()

    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    r, _ := regexp.Compile(*params.pattern)

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

            // Should exlude some paths if specified
            if isin(pathname, *params.exclude) { 
                return filepath.SkipDir
            }
            
            // The main routine produces paths that are buffered 
            // and consumed by the N workers 
            q.Enqueue(pathname)
            return nil
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

