package application

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/canta2899/mgp/model"
)

type Application struct {
	Wg       sync.WaitGroup
	Msg      model.OutputHandler
	Explorer model.PathWalk
	Options  *model.Options
}

func (app *Application) Run() {
	app.Explorer.Walk(app.getWalkFunction())
	app.Wg.Wait()
}

func (app *Application) getWalkFunction() filepath.WalkFunc {
	return func(pathname string, info os.FileInfo, err error) error {
		select {
		case <-app.Options.StopWalk:
			return errors.New("walk ended")
		default:
			e := model.FileInfo{FileInfo: info, Path: pathname}

			if err == nil {
				return app.processEntry(e)
			}

			app.Msg.AddPathError(e.Path, err)
			return nil
		}
	}
}

// Determines if the path entry should be skipped
func (app *Application) shouldSkip(f model.FileInfo) bool {
	isExcluded := matchCriteria(f, app.Options.Exclude)
	isIncluded := matchCriteria(f, app.Options.Include)
	exceedSize := f.Size() > int64(app.Options.LimitBytes) && !f.IsDir()

	if exceedSize || isExcluded {
		return true
	}

	if len(app.Options.Include) != 0 && !f.IsDir() {
		return !isIncluded
	}

	return false
}

// formatting text line in order to trim it
func formatMatchLine(line string) string {
	return strings.TrimSpace(strings.Trim(line, "\t"))
}

// Process path and enqueues if ok for match checking
func (app *Application) processEntry(f model.FileInfo) error {
	if app.shouldSkip(f) {
		if f.IsDir() {
			return filepath.SkipDir
		} else {
			return nil
		}
	}

	if f.IsDir() {
		return nil
	}

	app.Options.Running <- true // hangs if the buffer is full
	app.Wg.Add(1)
	go func() {
		match, err := app.match(f, app.Options.MatchAll)
		if err == nil && match != nil && len(match) != 0 {
			app.Msg.AddMatches(f.Path, match)
		}
		<-app.Options.Running // frees one position in the buffer
		app.Wg.Done()
	}()

	return nil
}

func (app *Application) match(f model.FileInfo, all bool) ([]*model.Match, error) {

	var m []*model.Match = nil

	// skipping regular files
	if !f.Mode().IsRegular() {
		return m, nil
	}

	file, err := os.Open(f.Path)

	// generally this is due to permission errors
	if err != nil {
		return m, err
	}

	defer file.Close()

	bufread := bufio.NewReader(file)

	m = []*model.Match{}
	count := 1 // counts line

	for {
		line, err := bufread.ReadBytes('\n')

		if err == io.EOF {
			break
		}

		if app.Options.Pattern.Match(line) {
			m = append(m, model.NewMatch(count, formatMatchLine(string(line))))
			if !all {
				// just return the first one if all is false
				return m, nil
			}
		}

		count += 1
	}

	return m, nil
}

// matches a file with a set of patterns and returns true when
// a match is found, false otherwise
func matchCriteria(f model.FileInfo, criteria []string) bool {
	for _, n := range criteria {
		fullMatch, _ := filepath.Match(n, f.Path)
		envMatch, _ := filepath.Match(n, filepath.Base(f.Path))
		if fullMatch || envMatch {
			return true
		}
	}

	return false
}
