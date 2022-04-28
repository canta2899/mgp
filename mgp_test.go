package main

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func IsPathExpected(expected []string, current string) bool {
	currentAbs, _ := filepath.Abs(filepath.Clean(current))

	for _, entry := range expected {
		expectedAbs, _ := filepath.Abs(filepath.Clean(entry))

		if expectedAbs == currentAbs {
			return true
		}
	}
	return false
}

func TestValidMatches(t *testing.T) {

	limit := 1048576 * 500

	var buf bytes.Buffer

	expectedOutputs := map[string][]string{
		"level1": {"./test/data.txt"},
		"level2": {"./test/level2/data.txt"},
		"level3": {"./test/level3/data.txt"},
	}

	for key, expected := range expectedOutputs {

		Run(&buf, 16, false, false, "./test", key, []string{".bzr", "CVS", ".git", ".hg", ".svn", ".idea", ".tox"}, limit)

		for {
			line, err := buf.ReadBytes('\n')

			obtained := strings.TrimSuffix(string(line), "\n")

			if err == io.EOF {
				break
			}

			if err != nil {
				t.Error("Error while reading command output")
			}

			if !IsPathExpected(expected, obtained) {
				t.Errorf("%v is not expected. Output should contain %v\n", string(line), expected)
			}

		}

	}
}
