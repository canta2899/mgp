package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
  env *Env
  node *FileInfo
}

func NewEntry(info os.FileInfo, path string, env *Env) *Entry {
  return &Entry{
    env:  env,
    node: &FileInfo{
      FileInfo: info,
      Path: path,
    },
  }
}

func (e *Entry) GetPath() string {
  return e.node.Path
}

func (e *Entry) ShouldSkip() bool {
  isDir := e.node.IsDir()

  for _, n := range e.env.exclude {
    fullMatch, _ := filepath.Match(n, e.node.Path)
    baseMatch, _ := filepath.Match(n, filepath.Base(e.node.Path))
    if isDir && (fullMatch || baseMatch) {
      return true
    }
  }

  return false
}

func (e *Entry) ShouldProcess() bool {
  isDir := e.node.IsDir()

  if isDir || e.node.Size() > int64(e.env.limitBytes) {
    return false
  }

  return true
}

func (e *Entry) MatchFirst() (*Match, error) {

  if !e.node.Mode().IsRegular() {
    return nil, nil
  }

  file, err := os.Open(e.node.Path)

  if err != nil {
    return nil, err
  }

  defer file.Close()

  bufread := bufio.NewReader(file)

  count := 1
  for {
    line, err := bufread.ReadBytes('\n')

    if err == io.EOF {
      break
    }

    if e.env.pattern.Match(line) {
      return &Match{LineNumber: count, Content: formatMatchLine(string(line))}, nil
    }
    count += 1
  }

	return nil, nil
}

func (e *Entry) MatchAll() ([]*Match, error) {

  var m []*Match = nil

  if !e.node.Mode().IsRegular() {
    return m, nil
  }

  file, err := os.Open(e.node.Path)

  if err != nil {
    return m, err
  }

  defer file.Close()

  bufread := bufio.NewReader(file)

  m = []*Match{}
  count := 1

  for {
    line, err := bufread.ReadBytes('\n')

    if err == io.EOF {
      break
    }

    if e.env.pattern.Match(line) {
      m = append(m, &Match{ 
        LineNumber: count, 
        Content:    formatMatchLine(string(line)),
      })
    }

    count += 1
  }

  return m, nil
}

func formatMatchLine(line string) string {
  return strings.TrimSpace(strings.Trim(line, "\t"))
}

