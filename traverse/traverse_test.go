package traverse

import (
	"errors"
	"path/filepath"
	"regexp"
	"sync"
	"testing"

	"github.com/canta2899/mgp/model"
	"github.com/canta2899/mgp/output"
)

func compileRegex(pattern string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		pattern = "(?i)" + pattern
	}

	if r, err := regexp.Compile(pattern); err == nil {
		return r, nil
	}

	return nil, errors.New("unable to compile regex pattern")
}

func IsPathExpected(expected []string, current string) bool {
	// todo should abs path be used for safety?
	// currentAbs, _ := filepath.Abs(filepath.Clean(current))

	for _, entry := range expected {
		// expectedAbs, _ := filepath.Abs(filepath.Clean(entry))

		if entry == current {
			return true
		}
	}
	return false
}

func TestValidMatches(t *testing.T) {

	limit := 1048576 * 500

	expectedOutputs := map[string][]string{
		"level1": {filepath.Join("testdata", "data.txt")},
		"level2": {filepath.Join("testdata", "level2", "data.txt")},
		"level3": {filepath.Join("testdata", "level3", "data.txt")},
	}

	for key, expected := range expectedOutputs {

		pattern, err := compileRegex(key, false)
		handler := output.NewTestOutputHandler()

		if err != nil {
			t.Fatal("Error compiling regexp")
		}

		stopWalk := false

		config := &model.Config{
			Wg:         sync.WaitGroup{},
			Schan:      make(chan bool, 16),
			Msg:        handler,
			Pattern:    pattern,
			StopWalk:   &stopWalk,
			StartPath:  "./testdata",
			Exclude:    []string{".bzr", "CVS", ".git", ".hg", ".svn", ".idea", ".tox"},
			LimitBytes: limit,
		}

		TraversePath(config)

		for _, value := range handler.Matches {
			if !IsPathExpected(expected, value) {
				t.Fatalf("%v is not expected. Output should contain %v\n", value, expected)
			}
		}

	}
}
