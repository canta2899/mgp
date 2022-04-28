package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {

	outputFile, err := os.Create("./output_test.txt")

	if err != nil {
		t.Error("Unable to create output test file")
	}

	wd, _ := os.Getwd()
	abs, _ := filepath.Abs("../logo-ls/")

	log.Println("Current path is", wd)

	limit := 1048576 * 500

	Run(outputFile, 16, false, false, abs, "Copy", []string{".bzr", "CVS", ".git", ".hg", ".svn", ".idea", ".tox"}, limit)
}
