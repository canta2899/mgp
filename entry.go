package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

type Entry struct {
	env *Env
	os.FileInfo
	Path string
}

func NewEntry(info os.FileInfo, path string, env *Env) *Entry {
	return &Entry{
		env:      env,
		FileInfo: info,
		Path:     path,
	}
}

func (e *Entry) ShouldSkip() bool {
	isDir := e.IsDir()

	for _, n := range e.env.exclude {
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

	if isDir || e.Size() > int64(e.env.limitBytes) {
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

		if e.env.pattern.Match(line) {
			return true, nil
		}
	}

	return false, nil
}
