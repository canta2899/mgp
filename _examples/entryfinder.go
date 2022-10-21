package main

import (
	"fmt"
	"log"
	"regexp"
	"sync"

	app "github.com/canta2899/mgp/application"
	"github.com/canta2899/mgp/internal/fspathwalk"
	"github.com/canta2899/mgp/internal/output"
	"github.com/canta2899/mgp/model"
)

// Runs mgp using a custom handler and iterates over the output
func main() {

	// the test output handler just saves entries inside a slice
	handler := output.NewTestOutputHandler()

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
		Wg:       sync.WaitGroup{},
		Msg:      handler,
		Explorer: explorer,
		Options: &model.Options{
			Running:    make(chan bool, 1000),
			StopWalk:   make(chan bool),
			MatchAll:   false,
			Pattern:    pattern,
			Exclude:    exclude,
			Include:    include,
			LimitBytes: 5000 * 1048576, // 5 gb
		},
	}

	log.Println("Searching for entries...")

	mgp.Run()

	// and then accesses the results
	for _, entry := range handler.Matches {
		fmt.Println("Found entry: ", entry)
	}
}
