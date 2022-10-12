package main

import (
  "log"
  "os"

  "github.com/fatih/color"
)

var okColor = color.New(color.FgHiGreen).SprintFunc()
var koColor = color.New(color.FgRed).SprintFunc()

type OutputHandler interface {
  AddMatch(path string)
  AddPathError(path string)
}

type FmtOutputHandler struct {
  Logger        *log.Logger
  ErrorLogger   *log.Logger
  Ok            string
  Ko            string
  OkColor       func(a ...interface{}) string
  KoColor       func(a ...interface{}) string
}

func NewFmtOutputHandler(colored bool) *FmtOutputHandler {
  var olog, elog *log.Logger

  if colored {
    olog = log.New(os.Stdout, okColor(string("\u2713" + " ")), 0)
    elog = log.New(os.Stderr, koColor(string("\u00D7" + " ")), 0)
  } else {
    olog = log.New(os.Stdout, "", 0)
    elog = log.New(os.Stderr, "", 0)
  }

  return &FmtOutputHandler{
    Logger: olog,
    ErrorLogger: elog,
  }
}

func (f *FmtOutputHandler) AddMatch(path string) {
  f.Logger.Println(path)
}

func (f *FmtOutputHandler) AddPathError(path string) {
  f.ErrorLogger.Println(path)
}


type TestOutputHandler struct {
  Matches []string 
  Errors  []string
}

func NewTestOutputHandler() *TestOutputHandler {
  return &TestOutputHandler{
    Matches: []string{},
    Errors:  []string{},
  }
}

func (f *TestOutputHandler) AddMatch(path string) {
  f.Matches = append(f.Matches, path)
}

func (f *TestOutputHandler) AddPathError(path string) {
  f.Errors = append(f.Errors, path)
}





















