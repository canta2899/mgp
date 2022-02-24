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
	"flag"
    "time"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"
)


type SyncControl struct {
    done bool
    mutex sync.Mutex
}

func (s *SyncControl) get() bool {
    defer s.mutex.Unlock()
    s.mutex.Lock()
    return s.done 
}

func (s *SyncControl) set(val bool) {
    defer s.mutex.Unlock() 
    s.mutex.Lock()
    s.done = val
}


type ExcludedDirs []string

func (ed *ExcludedDirs) String() string {
    return fmt.Sprintln(*ed)
}

func (ed *ExcludedDirs) Set(s string) error {
    fullpath, err := filepath.Abs(s)
    if err != nil {
        return err
    }
    *ed = append(*ed, fullpath)
    return nil
}


type Parameters struct {
    startpath string
    pattern   string
    workers   int
    excluded  ExcludedDirs
}


func ParseArgs() *Parameters {
    var exc ExcludedDirs

    w := flag.Int("w", 8, "degree of parallelism")
    flag.Var(&exc, "x", "excluded paths")

    flag.Parse()

    p := flag.Args()

    if len(p) < 2 {
        PrintUsageGuide()
        os.Exit(0)
    }

    return &Parameters{workers: *w, excluded: exc, startpath: p[1], pattern: p[0]}
}


func isin(el string, list ExcludedDirs) bool {

    // Utility that checks whether el is in list

    for _, entry := range list {
        if el == entry {
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


func handler(q *Queue, wg *sync.WaitGroup, r *regexp.Regexp, control *SyncControl) {
    defer wg.Done()

    for {

        filepath, err := q.Dequeue()

        if err != nil {

            if control.get() {
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
    var scanningdone *SyncControl = &SyncControl{done: false}
    var wg sync.WaitGroup

    params := ParseArgs()

    r, _ := regexp.Compile(params.pattern)

    q := NewQueue()

    wg.Add(params.workers)

    for i := 0; i < params.workers; i++ {
        go handler(q, &wg, r, scanningdone)
    }

    err := filepath.Walk(params.startpath,
        func(pathname string, info os.FileInfo, err error) error {

            // Should exlude some paths if specified
            if isin(path.Dir(pathname), params.excluded) {
                return filepath.SkipDir
            }

            if err != nil {
                return err
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
    scanningdone.set(true)

    // Waits for goroutines to finish
    wg.Wait()
}


