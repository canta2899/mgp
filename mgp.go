package main

import (
	"errors"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"
)

// WaitGroup for parallel goroutines
var wg sync.WaitGroup

// Limits the max number of goroutines running
var sChan chan bool

// Handles outputs
var m *MessageHandler

// Process path and enqueues if ok for match checking
func processPath(e *Entry) error {

	if e.ShouldSkip() {
		return filepath.SkipDir
	}

	if !e.ShouldProcess() {
		return nil
	}

	// hangs if the buffer is full
	sChan <- true
	// adds one goroutine to the wait group
	wg.Add(1)
	go func() {
		match, _ := e.HasMatch()

		if match {
			m.PrintSuccess(e.Path)
		}

		// frees one position in the buffer
		<-sChan
		// signals goroutine finished
		wg.Done()
	}()

	return nil
}

// Evaluates error for path and returns action to perform
func handlePathError(e *Entry, err error) error {

	// Prints error line for current path
	m.PrintError(e.Path)
	m.PrintInfo(err.Error())

	if e.IsDir() {
		return filepath.SkipDir
	}
	return nil
}

// Handler for sigterm (ctrl + c from cli)
func setSignalHandlers(stopWalk *bool) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch
		*stopWalk = true
	}()
}

func compileRegex(pattern string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		pattern = "(?i)" + pattern
	}

	if r, err := regexp.Compile(pattern); err == nil {
		return r, nil
	}

	return nil, errors.New("unable to compile regex pattern")
}

func Run(
	out io.Writer,
	workers int,
	caseInsensitive bool,
	colors bool,
	startpath string,
	pattern string,
	exc []string,
	limitMb int) {

	// Regex compilation
	r, err := compileRegex(pattern, caseInsensitive)

	if err != nil {
		m.PrintFatal("Invalid regex pattern")
	}

	if _, err := os.Stat(startpath); os.IsNotExist(err) {
		m.PrintFatal("Path does not exists")
	}

	stopWalk := false
	sChan = make(chan bool, workers)
	m = NewMessageHandler(colors, out)
	UpdateMatchingOptions(exc, int64(limitMb), r)

	setSignalHandlers(&stopWalk)

	// Traversing filepath
	filepath.Walk(startpath,

		func(pathname string, info os.FileInfo, err error) error {

			if stopWalk {
				// If the termination is requested, the path Walking
				// stops and the function returns with an error
				return errors.New("user requested termination")
			}

			e := NewEntry(info, pathname)

			// Checking permission and access errors
			if err != nil {
				return handlePathError(e, err)
			}

			// Processes path in search of matches with the given
			// pattern or the excluded directories
			return processPath(e)
		})

	// Waits for goroutines to finish
	wg.Wait()

	if stopWalk {
		m.PrintInfo("Ended by user")
	}
}
