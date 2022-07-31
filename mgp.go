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

var wg sync.WaitGroup
var sChan chan int
var m *MessageHandler

// Process path and enqueues if ok for match checking
func processPath(e *Entry, exc []string, limitMb int, r *regexp.Regexp) error {

	if e.ShouldSkip(exc) {
		return filepath.SkipDir
	}

	if !e.ShouldProcess(limitMb) {
		return nil
	}

	wg.Add(1)
	sChan <- 0
	go func() {
		match, _ := e.HasMatch(r)

		if match {
			m.printSuccess(e.Path)
		}

		<-sChan
		wg.Done()
	}()

	return nil
}

// Evaluates error for path and returns action to perform
func handlePathError(e *Entry, err error) error {

	if os.IsNotExist(err) {
		m.printFatal("Invalid path")
	}

	// Prints error line for current path
	m.printError(e.Path)
	m.printInfo(err.Error())

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

func Run(
	out io.Writer,
	workers int,
	caseInsensitive bool,
	colors bool,
	startpath string,
	pattern string,
	exludedDirs []string,
	limitMb int) {

	sChan = make(chan int, workers)
	m = NewMessageHandler(colors, out)

	// Regex compilation
	if caseInsensitive {
		pattern = "(?i)" + pattern
	}
	r, err := regexp.Compile(pattern)

	if err != nil {
		m.printInfo("Error in regex pattern")
		os.Exit(1)
	}

	stopWalk := false
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
			return processPath(e, exludedDirs, limitMb, r)
		})

	// Waits for goroutines to finish
	wg.Wait()

	if stopWalk {
		m.printInfo("Ended by user")
	}
}
