package traverse

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/canta2899/mgp/model"
)

type Inspector struct {
	File   model.FileInfo
	Config *model.Config
}

// Allows inspection of given file in search of matches
func NewInspector(file model.FileInfo, config *model.Config) *Inspector {
	return &Inspector{
		File:   file,
		Config: config,
	}
}

// Determines if the path entry should be skipped
func (i *Inspector) ShouldSkip() bool {
	isDir := i.File.IsDir()

	for _, n := range i.Config.Exclude {
		fullMatch, _ := filepath.Match(n, i.File.Path)
		envMatch, _ := filepath.Match(n, filepath.Base(i.File.Path))
		if isDir && (fullMatch || envMatch) {
			return true
		}
	}

	return false
}

// Determines if the path entry should be processed
func (e *Inspector) ShouldProcess() bool {
	isDir := e.File.IsDir()

	if isDir || e.File.Size() > int64(e.Config.LimitBytes) {
		return false
	}

	return true
}

// formatting text line in order to trim it
func formatMatchLine(line string) string {
	return strings.TrimSpace(strings.Trim(line, "\t"))
}

// opens a file in search of matches,
// all = true allows to get all matching lines from a file
func (e *Inspector) Match(all bool) ([]*model.Match, error) {

	var m []*model.Match = nil

	// non regular files should not be processed
	if !e.File.Mode().IsRegular() {
		return m, nil
	}

	file, err := os.Open(e.File.Path)

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

		if e.Config.Pattern.Match(line) {
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
