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

var scanningdone *SyncControl = &SyncControl{done: false}

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

    pattern := flag.Arg(0)
    startpath := flag.Arg(1)

    return &Parameters{workers: *w, excluded: exc, startpath: startpath, pattern: pattern}
}

func handler(c <-chan interface{}, wg *sync.WaitGroup, r *regexp.Regexp) {
    defer wg.Done()

    for {

        if len(c) == 0 && scanningdone.get() {
            return
        }

        next := <-c

        filepath := next.(string) // type assertion

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
            fmt.Println(next)
        }
    }
}

func main() {

    // Allows synchronization in order to let the
    // main function wait for the other coroutines
    var wg sync.WaitGroup

    params := ParseArgs()

    r, _ := regexp.Compile(params.pattern)

    c := make(chan interface{}, 1000)

    wg.Add(params.workers)

    for i := 0; i < params.workers; i++ {
        go handler(c, &wg, r)
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
            
                c <- pathname

            return nil
        })

    if err != nil {
        panic(err.Error())
    }

    // Using a channel in order to implement a thread
    // safe queue that allows multiple consumers to
    // "dequeue" each path and process the content of
    // the referred file

    scanningdone.set(true)

    wg.Wait()
}

// Utility that checks whether el is in list
func isin(el string, list ExcludedDirs) bool {
    for _, entry := range list {
        if el == entry {
            return true
        }
    }
    return false
}
