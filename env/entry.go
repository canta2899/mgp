package traverse

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
  "github.com/canta2899/mgp/model"
)

type Entry struct {
  env  *Env
  Node *model.FileInfo
}

func NewEntry(info os.FileInfo, path string, env *Env) *Entry {
  return &Entry{
    env:  env,
    Node: &model.FileInfo{
      FileInfo: info,
      Path: path,
    },
  }
}

func (e *Entry) GetPath() string {
  return e.Node.Path
}

func (e *Entry) ShouldSkip() bool {
  isDir := e.Node.IsDir()

  for _, n := range e.env.Exclude {
    fullMatch, _ := filepath.Match(n, e.Node.Path)
    envMatch, _ := filepath.Match(n, filepath.Base(e.Node.Path))
    if isDir && (fullMatch || envMatch) {
      return true
    }
  }

  return false
}

func (e *Entry) ShouldProcess() bool {
  isDir := e.Node.IsDir()

  if isDir || e.Node.Size() > int64(e.env.LimitBytes) {
    return false
  }

  return true
}

func (e *Entry) MatchFirst() (*model.Match, error) {

  if !e.Node.Mode().IsRegular() {
    return nil, nil
  }

  file, err := os.Open(e.Node.Path)

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

    if e.env.Pattern.Match(line) {
      return &model.Match{LineNumber: count, Content: formatMatchLine(string(line))}, nil
    }
    count += 1
  }

	return nil, nil
}

func (e *Entry) MatchAll() ([]*model.Match, error) {

  var m []*model.Match = nil

  if !e.Node.Mode().IsRegular() {
    return m, nil
  }

  file, err := os.Open(e.Node.Path)

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

    if e.env.Pattern.Match(line) {
      m = append(m, &model.Match{ 
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

