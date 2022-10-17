package services

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/canta2899/mgp/pkg/model"
)

type Application struct {
	Wg         sync.WaitGroup
	Running    chan bool
	StopWalk   chan bool
	MatchAll   bool
	Msg        model.OutputHandler
	Pattern    *regexp.Regexp
	Explorer   model.PathExplorer
	Exclude    []string
	LimitBytes int
}

func (c *Application) Run() {
	c.Explorer.Walk(evalPath(c))
	c.Wg.Wait()
}

// Process path and enqueues if ok for match checking
func processEntry(f model.FileInfo, config *Application) error {
	i := NewInspector(f, config)

	if i.ShouldSkip() {
		return filepath.SkipDir
	}

	if !i.ShouldProcess() {
		return nil
	}

	// hangs if the buffer is full
	config.Running <- true
	// adds one goroutine to the wait group
	config.Wg.Add(1)
	go func() {
		match, err := i.Match(config.MatchAll)
		if err == nil && match != nil && len(match) != 0 {
			config.Msg.AddMatches(i.File.Path, match)
		}

		// frees one position in the buffer
		<-config.Running
		// signals goroutine finished
		config.Wg.Done()
	}()

	return nil
}

// Traversing filepath
func evalPath(config *Application) func(pathname string, info os.FileInfo, err error) error {
	return func(pathname string, info os.FileInfo, err error) error {
		select {
		case <-config.StopWalk:
			return errors.New("walk ended")
		default:
			e := model.FileInfo{FileInfo: info, Path: pathname}

			// Processes path in search of matches with the given
			// pattern or the excluded directories
			if err == nil {
				return processEntry(e, config)
			}

			// Checking permission and access errors
			config.Msg.AddPathError(e.Path, err)

			if e.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}
	}
}
