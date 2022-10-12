package main

import (
	"path/filepath"
	"sync"
	"testing"
)

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
    handler := NewTestOutputHandler()

		if err != nil {
			t.Fatal("Error compiling regexp")
		}

		stopWalk := false

		env := &Env{
			wg:         sync.WaitGroup{},
			sChan:      make(chan bool, 16),
			msg:        handler,
			pattern:    pattern,
			stopWalk:   &stopWalk,
			startpath:  "./testdata",
			exclude:    []string{".bzr", "CVS", ".git", ".hg", ".svn", ".idea", ".tox"},
			limitBytes: limit,
		}

		env.Run()

    for _, value := range handler.Matches {
      if !IsPathExpected(expected, value) {
        t.Fatalf("%v is not expected. Output should contain %v\n", value, expected)
      }
    }

	}
}
