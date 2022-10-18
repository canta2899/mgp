package application

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	Explorer   model.PathWalk
	Exclude    []string
	LimitBytes int
}

func (app *Application) Run() {
	app.Explorer.Walk(app.getWalkFunction())
	app.Wg.Wait()
}

// Traversing filepath
func (app *Application) getWalkFunction() func(pathname string, info os.FileInfo, err error) error {
	return func(pathname string, info os.FileInfo, err error) error {
		select {
		case <-app.StopWalk:
			return errors.New("walk ended")
		default:
			e := model.FileInfo{FileInfo: info, Path: pathname}

			// Processes path in search of matches with the given
			// pattern or the excluded directories
			if err == nil {
				return app.processEntry(e)
			}

			// Checking permission and access errors
			app.Msg.AddPathError(e.Path, err)

			if e.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}
	}
}

// Determines if the path entry should be skipped
func (app *Application) shouldSkip(f model.FileInfo) bool {
	isDir := f.IsDir()

	for _, n := range app.Exclude {
		fullMatch, _ := filepath.Match(n, f.Path)
		envMatch, _ := filepath.Match(n, filepath.Base(f.Path))
		if isDir && (fullMatch || envMatch) {
			return true
		}
	}

	return false
}

// Determines if the path entry should be processed
func (app *Application) shouldProcess(f model.FileInfo) bool {
	isDir := f.IsDir()

	if isDir || f.Size() > int64(app.LimitBytes) {
		return false
	}

	return true
}

// formatting text line in order to trim it
func formatMatchLine(line string) string {
	return strings.TrimSpace(strings.Trim(line, "\t"))
}

// Process path and enqueues if ok for match checking
func (app *Application) processEntry(f model.FileInfo) error {
	if app.shouldSkip(f) {
		return filepath.SkipDir
	}

	if !app.shouldProcess(f) {
		return nil
	}

	// hangs if the buffer is full
	app.Running <- true
	// adds one goroutine to the wait group
	app.Wg.Add(1)
	go func() {
		match, err := app.match(f, app.MatchAll)
		if err == nil && match != nil && len(match) != 0 {
			app.Msg.AddMatches(f.Path, match)
		}

		// frees one position in the buffer
		<-app.Running
		// signals goroutine finished
		app.Wg.Done()
	}()

	return nil
}

// opens a file in search of matches,
// all = true allows to get all matching lines from a file
func (app *Application) match(f model.FileInfo, all bool) ([]*model.Match, error) {

	var m []*model.Match = nil

	// non regular files should not be processed
	if !f.Mode().IsRegular() {
		return m, nil
	}

	file, err := os.Open(f.Path)

	// might happend in case of permission errors
	if err != nil {
		return m, err
	}

	defer file.Close()

	bufread := bufio.NewReader(file)

	m = []*model.Match{}
	count := 1

	for {
		line, err := bufread.ReadBytes('\n')

		if err == io.EOF {
			break
		}

		if app.Pattern.Match(line) {
			// build match according to regexp match output
			m = append(m, model.NewMatch(count, formatMatchLine(string(line))))

			if !all {
				// just return the first one if all is false
				return m, nil
			}
		}

		count += 1 // counts the line of the file
	}

	return m, nil
}
