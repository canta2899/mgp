package traverse

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/canta2899/mgp/model"
)

// Process path and enqueues if ok for match checking
func ProcessEntry(f model.FileInfo, config *model.Config) error {

	i := NewInspector(f, config)

	if i.ShouldSkip() {
		return filepath.SkipDir
	}

	if !i.ShouldProcess() {
		return nil
	}

	// hangs if the buffer is full
	config.Schan <- true
	// adds one goroutine to the wait group
	config.Wg.Add(1)
	go func() {
		match, err := i.Match(config.MatchAll)
		if err == nil && match != nil && len(match) != 0 {
			config.Msg.AddMatches(i.File.Path, match)
		}

		// frees one position in the buffer
		<-config.Schan
		// signals goroutine finished
		config.Wg.Done()
	}()

	return nil
}

func TraversePath(config *model.Config) {

	if _, err := os.Stat(config.StartPath); os.IsNotExist(err) {
		config.Msg.AddPathError(config.StartPath, errors.New("path does not exists"))
		os.Exit(1)
	}

	// Traversing filepath
	filepath.Walk(config.StartPath,

		func(pathname string, info os.FileInfo, err error) error {

			if *config.StopWalk {
				// If the termination is requested, the path Walking
				// stops and the function returns with an error
				return errors.New("user requested termination")
			}

			e := model.FileInfo{FileInfo: info, Path: pathname}

			// Processes path in search of matches with the given
			// pattern or the excluded directories
			if err == nil {
				return ProcessEntry(e, config)
			}

			// Checking permission and access errors
			config.Msg.AddPathError(e.Path, err)

			if e.IsDir() {
				return filepath.SkipDir
			}

			return nil
		})

	// Waits for goroutines to finish
	config.Wg.Wait()
}
