package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

type Entry struct {
	os.FileInfo
	Path string
}

func NewEntry(info os.FileInfo, path string) *Entry {
	return &Entry{
		FileInfo: info,
		Path:     path,
	}
}

func (e *Entry) ShouldSkip(excludedDirs []string) bool {
	isDir := e.IsDir()

	for _, n := range excludedDirs {
		fullMatch, _ := filepath.Match(n, e.Path)
		baseMatch, _ := filepath.Match(n, filepath.Base(e.Path))
		if isDir && (fullMatch || baseMatch) {
			return true
		}
	}

	return false
}

func (e *Entry) ShouldProcess(limitMb int) bool {
	isDir := e.IsDir()

	if isDir || e.Size() > int64(limitMb) {
		return false
	}

	return true
}

func (e *Entry) HasMatch(r *regexp.Regexp) (bool, error) {

	if !e.Mode().IsRegular() {
		return false, nil
	}

	file, err := os.Open(e.Path)

	if err != nil {
		return false, err
	}

	defer file.Close()

	bufread := bufio.NewReader(file)

	for {
		line, err := bufread.ReadBytes('\n')

		if err == io.EOF {
			break
		}

		if r.Match(line) {
			return true, nil
		}
	}

	return false, nil
}
