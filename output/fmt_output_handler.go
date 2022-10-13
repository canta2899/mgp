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
	Logger      *log.Logger
	ErrorLogger *log.Logger
	Ok          string
	Ko          string
	OkColor     func(a ...interface{}) string
	KoColor     func(a ...interface{}) string
}

func NewFmtOutputHandler(colored bool) *FmtOutputHandler {
	var olog, elog *log.Logger

	if colored {
		olog = log.New(os.Stdout, okColor(string("\u2713"+" ")), 0)
		elog = log.New(os.Stderr, koColor(string("\u00D7"+" ")), 0)
	} else {
		olog = log.New(os.Stdout, "", 0)
		elog = log.New(os.Stderr, "", 0)
	}

	return &FmtOutputHandler{
		Logger:      olog,
		ErrorLogger: elog,
	}
}

func (f *FmtOutputHandler) AddMatch(path string, match *model.Match) {
	f.Logger.Println(path)
}

func (f *FmtOutputHandler) AddMatches(path string, matches []*model.Match) {
	for _, m := range matches {
		f.Logger.Printf("%v:%v:  %v", path, m.LineNumber, m.Content)
	}
}

func (f *FmtOutputHandler) AddPathError(path string, e error) {
	f.ErrorLogger.Printf("%v %v\n", path, e.Error())
}
