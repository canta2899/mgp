package main

import (
	"log"
	"os"

	"github.com/fatih/color"
)

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

func NewFmtOutputHandler() *FmtOutputHandler {
	return &FmtOutputHandler{
		Logger:        log.New(os.Stdout, "", 0),
    ErrorLogger:   log.New(os.Stderr, "", 0),
		Ok:            string("\u2713"),
		Ko:            string("\u00D7"),
		OkColor:       color.New(color.FgHiGreen).SprintFunc(),
		KoColor:       color.New(color.FgRed).SprintFunc(),
	}
}

func (f *FmtOutputHandler) AddMatch(path string) {

}

func (f *FmtOutputHandler) AddPathError(path string) {

}
