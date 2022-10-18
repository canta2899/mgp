package output

import (
	"github.com/canta2899/mgp/pkg/model"
)

type TestOutputHandler struct {
	Matches []string
	Errors  []string
}

// Mockup output handler that allows to collect paths for
// testing purposes
func NewTestOutputHandler() *TestOutputHandler {
	return &TestOutputHandler{
		Matches: []string{},
		Errors:  []string{},
	}
}

func (f *TestOutputHandler) AddMatches(path string, matches []*model.Match) {
	f.Matches = append(f.Matches, path)
}

func (f *TestOutputHandler) AddPathError(path string, e error) {
	f.Errors = append(f.Errors, path)
}
