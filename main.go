package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"
	"github.com/fatih/color"
)

type Entry struct {
    Path string
    Info *os.FileInfo
}

type MessageType int64

const (
    Match MessageType = iota
    ErrorMatch
    Log
)

var coloredOutput bool = true

// Runes for emoji
const OK string = string('\u2713')
const KO string = string('\u00D7')

// Colors for printing
var green =  color.New(color.FgHiGreen).SprintFunc()
var red   =  color.New(color.FgRed).SprintFunc()

func printHandler(message string, messageType MessageType) {
    switch messageType {
    case Match:
        if coloredOutput {
            log.Printf("%v %v\n", green(OK), message)
        } else {
            log.Printf("%v\n", message)
        }
        return;
    case ErrorMatch:
        if coloredOutput {
            log.Printf("%v %v\n", red(KO), message)
        }
        return;
    }

    log.Printf(message)
}


// Routine performed by each worker
func handler(ch <-chan *Entry, wg *sync.WaitGroup, r *regexp.Regexp) {
    defer wg.Done()

    for {

        e, more := <-ch

        if !more {
            return
        }

        info, fullpath := e.Info, e.Path

        if !(*info).Mode().IsRegular() {
            continue // Skips
        }

        file, err := os.Open(fullpath)

        if err != nil {
            continue // Skips
        }

        bufread := bufio.NewReader(file)

        for {
            line, err := bufread.ReadBytes('\n')

            if err == io.EOF {
                break
            }

            if r.Match(line) {
                printHandler(fullpath, Match)
                break
            }
        }
        file.Close()
    }
}

// Process path and enqueues if ok for match checking
func processPath(info *os.FileInfo, pathname string, c chan *Entry, params *Parameters) error {
    isdir := (*info).IsDir()
    exc := params.GetExcludedDirs()
    lim := params.limitMb

    for _, n := range exc {
        fullMatch, _ := filepath.Match(n, pathname)
        baseMatch, _ := filepath.Match(n, filepath.Base(pathname))
        if isdir && (fullMatch || baseMatch) {
            return filepath.SkipDir 
        }
    }

    if !isdir && (*info).Size() < int64(lim) {
        c <- &Entry{Path: pathname, Info: info} 
    }

    return nil
}

// Evaluates error for path and returns action to perform
func handlePathError(info *os.FileInfo, pathname string, err error) error {

    if os.IsNotExist(err) { log.Fatal("Invalid path") }

    // Prints error line for current path
    printHandler(pathname, ErrorMatch)
    printHandler(err.Error(), Log)
    
    if (*info).IsDir() {
        return filepath.SkipDir
    } else {
        return nil
    }
}

// Handler for sigterm (ctrl + c from cli)
func setSignalHandlers(closed *bool, wg *sync.WaitGroup) {
    sigch := make(chan os.Signal, 1)
    signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-sigch 
        *closed = true
    }()
}

func setupLogger() {
    // Message wih no flags printed to stdout
    log.SetFlags(0)
    log.SetOutput(os.Stdout)
}

func main() {
    var wg sync.WaitGroup
    ch := make(chan *Entry, 5000)
    closed := false

    setupLogger()
    setSignalHandlers(&closed, &wg)
    params := ParseArgs()

    pattern := params.pattern
    coloredOutput = !params.nocolor
    
    if params.icase {
        pattern = "(?i)" + pattern
    }

    r, _ := regexp.Compile(pattern)

    wg.Add(params.workers)
    for i := 0; i < params.workers; i++ {
        go handler(ch, &wg, r)
    }

    filepath.Walk(params.startpath,

        func(pathname string, info os.FileInfo, err error) error {

            if closed {
                // If the termination is requested, the path Walking
                // stops and the function returns with an error
                return errors.New("User requested termination")
            }

            // Checking permission and access errors
            if err != nil {
                return handlePathError(&info, pathname, err)
            }

            // Processes path in search of matches with the given
            // pattern or the excluded directories 
            return processPath(&info, pathname, ch, params)

        })

    // The channel is closed, this communicates that
    // no more values will be enqueued
    close(ch)

    // Waits for goroutines to finish
    wg.Wait()
}

