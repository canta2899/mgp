package main

import (
	"bytes"
	"io"
	"log"
	"strings"
	"sync"
	"testing"
)

func IsPathExpected(expected []string, current string) bool {
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

	var buf bytes.Buffer

	expectedOutputs := map[string][]string{
		"level1": {"testdata/data.txt"},
		"level2": {"testdata/level2/data.txt"},
		"level3": {"testdata/level3/data.txt"},
	}

	for key, expected := range expectedOutputs {

		pattern, err := compileRegex(key, false)

		if err != nil {
			t.Fatal("Error compiling regexp")
		}

		env := &env{
			wg:         sync.WaitGroup{},
			sChan:      make(chan bool, 16),
			msg:        NewMessageHandler(false, &buf),
			pattern:    pattern,
			startpath:  "./testdata",
			exclude:    []string{".bzr", "CVS", ".git", ".hg", ".svn", ".idea", ".tox"},
			limitBytes: limit,
		}

		env.Run()

		log.Println("here")

		for {
			line, err := buf.ReadBytes('\n')

			if err == io.EOF {
				break
			}

			obtained := strings.TrimSuffix(string(line), "\n")

			if err != nil {
				t.Error("Error while reading command output")
			}

			if !IsPathExpected(expected, obtained) {
				t.Errorf("%v is not expected. Output should contain %v\n", obtained, expected)
			}

		}
	}
}
