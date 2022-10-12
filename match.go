package main

import (
  "strconv"
)

type Match struct {
  Path    string
  Context map[string]string
}

func NewMatch(path string) *Match {
  return &Match{
    Path: path,
    Context: nil,
  }
}

func (m *Match) AddContext(lineNumb int, content string) {
  lineNumbStr := strconv.Itoa(lineNumb)
  m.Context[lineNumbStr] = content
}

