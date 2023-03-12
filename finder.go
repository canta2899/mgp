package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var wg sync.WaitGroup = sync.WaitGroup{}
var running chan bool = nil
var stopWalk chan bool = nil

type Finder struct {
	Msg      OutputHandler
	Explorer *PathWalk
	Options  *Options
}

type Options struct {
	MatchAll   bool
	Pattern    *regexp.Regexp
	Exclude    []string
	Include    []string
	LimitBytes int
}

func (app *Finder) Run(maxWorkers int) {
	running = make(chan bool, maxWorkers)
	stopWalk = make(chan bool)
	app.Explorer.Walk(app.getWalkFunction())
	wg.Wait()
	close(stopWalk)
}

func (app *Finder) Stop() {
	stopWalk <- true
	<-stopWalk
}

func (app *Finder) getWalkFunction() filepath.WalkFunc {
	return func(pathname string, info os.FileInfo, err error) error {
		select {
		case <-stopWalk:
			return errors.New("walk ended")
		default:
			e := FileInfo{FileInfo: info, Path: pathname}

			if err == nil {
				return app.processEntry(e)
			}

			app.Msg.AddPathError(e.Path, err)
			return nil
		}
	}
}

// Determines if the path entry should be skipped
func (app *Finder) shouldSkip(f FileInfo) bool {
	// skipping regular files
	if !f.IsDir() && !f.Mode().IsRegular() {
		return true
	}

	isExcluded := matchCriteria(f, app.Options.Exclude)
	isIncluded := matchCriteria(f, app.Options.Include)
	exceedSize := (int64(app.Options.LimitBytes) != 0 && f.Size() > int64(app.Options.LimitBytes) && !f.IsDir())

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
func (app *Finder) processEntry(f FileInfo) error {
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

	running <- true // hangs if the buffer is full
	wg.Add(1)
	go func() {
		match, err := app.match(f, app.Options.MatchAll)
		if err == nil && match != nil && len(match) != 0 {
			app.Msg.AddMatches(f.Path, match)
		}
		<-running // frees one position in the buffer
		wg.Done()
	}()

	return nil
}

func (app *Finder) match(f FileInfo, all bool) ([]*Match, error) {
	var m []*Match = nil

	file, err := os.Open(f.Path)

	// generally this is due to permission errors
	if err != nil {
		return m, err
	}

	defer file.Close()

	bufread := bufio.NewReader(file)

	m = []*Match{}
	count := 1 // counts line

	for {
		line, err := bufread.ReadBytes('\n')

		if err == io.EOF {
			break
		}

		if app.Options.Pattern.Match(line) {
			m = append(m, NewMatch(count, formatMatchLine(string(line))))
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
func matchCriteria(f FileInfo, criteria []string) bool {
	for _, n := range criteria {
		fullMatch, _ := filepath.Match(n, f.Path)
		envMatch, _ := filepath.Match(n, filepath.Base(f.Path))
		if fullMatch || envMatch {
			return true
		}
	}

	return false
}
