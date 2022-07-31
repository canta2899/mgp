package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

var limitMb int64
var excludedDirs []string
var pattern *regexp.Regexp

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

func (e *Entry) ShouldSkip() bool {
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

func (e *Entry) ShouldProcess() bool {
	isDir := e.IsDir()

	if isDir || e.Size() > limitMb {
		return false
	}

	return true
}

func (e *Entry) HasMatch() (bool, error) {

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

		if pattern.Match(line) {
			return true, nil
		}
	}

	return false, nil
}

func UpdateMatchingOptions(exc []string, limit int64, p *regexp.Regexp) {
	excludedDirs = exc
	limitMb = limit
	pattern = p
}
