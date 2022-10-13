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

func NewInspector(file model.FileInfo, config *model.Config) *Inspector {
	return &Inspector{
		File:   file,
		Config: config,
	}
}

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

func (e *Inspector) ShouldProcess() bool {
	isDir := e.File.IsDir()

	if isDir || e.File.Size() > int64(e.Config.LimitBytes) {
		return false
	}

	return true
}

func (e *Inspector) Match(all bool) ([]*model.Match, error) {

	var m []*model.Match = nil

	if !e.File.Mode().IsRegular() {
		return m, nil
	}

	file, err := os.Open(e.File.Path)

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
			m = append(m, &model.Match{
				LineNumber: count,
				Content:    formatMatchLine(string(line)),
			})

			if !all {
				return m, nil
			}
		}

		count += 1
	}

	return m, nil
}

func formatMatchLine(line string) string {
	return strings.TrimSpace(strings.Trim(line, "\t"))
}
