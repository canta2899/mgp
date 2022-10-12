package main

import (
  "bufio"
  "io"
  "os"
  "path/filepath"
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

func (e *Entry) HasMatch() (bool, error) {

  if !e.node.Mode().IsRegular() {
    return false, nil
  }

  file, err := os.Open(e.node.Path)

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
