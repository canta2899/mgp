package main

import (
	"fmt"
	"log"
	"regexp"

	app "github.com/canta2899/mgp/application"
	"github.com/canta2899/mgp/internal/mockups"
	"github.com/canta2899/mgp/model"
	"github.com/canta2899/mgp/pkg/fspathwalk"
)

// Runs mgp using a custom handler and iterates over the output
func main() {

	// the test output handler just saves entries inside a slice
	handler := mockups.NewTestOutputHandler()

	// walking inside current path
	explorer := fspathwalk.NewFsPathWalk(".")

	// compiles case insensitive "level" regexp
	pattern, err := regexp.Compile("(?i)level")

	if err != nil {
		log.Fatalf(err.Error())
	}

	// sames as cli program defaults
	exclude := []string{".bzr", "CVS", ".git", ".hg", ".svn", ".idea", ".tox", "node_modules"}
	// if empty includes all
	include := []string{}

	// configures application
	mgp := &app.Application{
		Msg:      handler,
		Explorer: explorer,
		Options: &model.Options{
			MatchAll:   false,
			Pattern:    pattern,
			Exclude:    exclude,
			Include:    include,
			LimitBytes: 5000 * 1048576, // 5 gb
		},
	}

	log.Println("Searching for entries...")

	mgp.Run(100) // run with 100 workers

	// and then accesses the results
	for _, entry := range handler.Matches {
		fmt.Println("Found entry: ", entry)
	}
}
