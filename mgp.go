package main

import (
	"errors"
	"os"
	"path/filepath"
)

// Process path and enqueues if ok for match checking
func (env *Env) processEntry(e *Entry) error {

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
			env.msg.PrintSuccess(e.GetPath())
		}

		// frees one position in the buffer
		<-env.sChan
		// signals goroutine finished
		env.wg.Done()
	}()

	return nil
}

// Evaluates error for path and returns action to perform
func (env *Env) handleEntryError(e *Entry, err error) error {

	// Prints error line for current path
	env.msg.PrintError(e.GetPath())
	env.msg.PrintInfo(err.Error())

	if e.node.IsDir() {
		return filepath.SkipDir
	}
	return nil
}

func (env *Env) Run() {

	if _, err := os.Stat(env.startpath); os.IsNotExist(err) {
		env.msg.PrintFatal("Path does not exists")
	}

	// Traversing filepath
	filepath.Walk(env.startpath,

		func(pathname string, info os.FileInfo, err error) error {

			if *env.stopWalk {
				// If the termination is requested, the path Walking
				// stops and the function returns with an error
				return errors.New("user requested termination")
			}
			e := NewEntry(info, pathname, env)

			// Checking permission and access errors
			if err != nil {
				return env.handleEntryError(e, err)
			}

			// Processes path in search of matches with the given
			// pattern or the excluded directories
			return env.processEntry(e)
		})

	// Waits for goroutines to finish
	env.wg.Wait()

	if *env.stopWalk {
		env.msg.PrintInfo("Ended by user")
	}
}
