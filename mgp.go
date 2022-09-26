package main

import (
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"
)

type env struct {
	wg         sync.WaitGroup
	sChan      chan bool
	msg        *MessageHandler
	pattern    *regexp.Regexp
	startpath  string
	exclude    []string
	limitBytes int
}

// Process path and enqueues if ok for match checking
func (env *env) processPath(e *Entry) error {

	if e.ShouldSkip() {
		return filepath.SkipDir
	}

	if !e.ShouldProcess() {
		return nil
	}

	// hangs if the buffer is full
	env.sChan <- true
	// adds one goroutine to the wait group
	env.wg.Add(1)
	go func() {
		match, _ := e.HasMatch()

		if match {
			env.msg.PrintSuccess(e.Path)
		}

		// frees one position in the buffer
		<-env.sChan
		// signals goroutine finished
		env.wg.Done()
	}()

	return nil
}

// Evaluates error for path and returns action to perform
func (env *env) handlePathError(e *Entry, err error) error {

	// Prints error line for current path
	env.msg.PrintError(e.Path)
	env.msg.PrintInfo(err.Error())

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

func (env *env) Run() {

	if _, err := os.Stat(env.startpath); os.IsNotExist(err) {
		env.msg.PrintFatal("Path does not exists")
	}

	stopWalk := false

	setSignalHandlers(&stopWalk)

	// Traversing filepath
	filepath.Walk(env.startpath,

		func(pathname string, info os.FileInfo, err error) error {

			if stopWalk {
				// If the termination is requested, the path Walking
				// stops and the function returns with an error
				return errors.New("user requested termination")
			}

			e := env.NewEntry(info, pathname)

			// Checking permission and access errors
			if err != nil {
				return env.handlePathError(e, err)
			}

			// Processes path in search of matches with the given
			// pattern or the excluded directories
			return env.processPath(e)
		})

	// Waits for goroutines to finish
	env.wg.Wait()

	if stopWalk {
		env.msg.PrintInfo("Ended by user")
	}
}
