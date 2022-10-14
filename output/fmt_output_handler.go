package output

import (
	"log"
	"os"

	"github.com/canta2899/mgp/model"
	"github.com/fatih/color"
)

var okColor = color.New(color.FgHiGreen).SprintFunc()
var koColor = color.New(color.FgRed).SprintFunc()

type FmtOutputHandler struct {
	Logger       *log.Logger
	ErrorLogger  *log.Logger
	Ok           string
	Ko           string
	OkColor      func(a ...interface{}) string
	KoColor      func(a ...interface{}) string
	PrintContext bool
}

func NewFmtOutputHandler(colored, printContext bool) *FmtOutputHandler {
	var olog, elog *log.Logger

	if colored {
		olog = log.New(os.Stdout, okColor(string("\u2713"+" ")), 0)
		elog = log.New(os.Stderr, koColor(string("\u00D7"+" ")), 0)
	} else {
		olog = log.New(os.Stdout, "", 0)
		elog = log.New(os.Stderr, "", 0)
	}

	return &FmtOutputHandler{
		Logger:       olog,
		ErrorLogger:  elog,
		PrintContext: printContext,
	}
}

func (f *FmtOutputHandler) AddMatches(path string, matches []*model.Match) {
	for _, m := range matches {
		if f.PrintContext {
			f.Logger.Printf("%v:%v:  %v", path, m.LineNumber, m.Content)
		} else {
			f.Logger.Printf("%v", path)
		}
	}
}

func (f *FmtOutputHandler) AddPathError(path string, e error) {
	f.ErrorLogger.Printf("%v %v\n", path, e.Error())
}
