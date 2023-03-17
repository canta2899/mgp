package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

func compileRegex(pattern string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		pattern = "(?i)" + pattern
	}
	return regexp.Compile(pattern)
}

// Sets a SIGTERM handler in order to stop when Ctrl+C is pressed
func setSignalHandlers(app *Finder) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch
		app.Stop()
	}()
}

// Runs mgp from command line interface
func main() {
  args := os.Args[1:]
	params := ParseArgs(args)

	if _, err := os.Stat(params.startpath); os.IsNotExist(err) {
		log.Println("Path does not exists")
		os.Exit(1)
	}

	outputHandler := NewFmtOutputHandler(!params.raw, params.showCtx)
	pathWalker := NewPathWalk(params.startpath)

	pattern, err := compileRegex(params.pattern, params.icase)

	if err != nil {
		log.Fatalf(err.Error())
	}

	app := &Finder{
		Msg:      outputHandler,
		Explorer: pathWalker,
		Options: &Options{
			MatchAll:   params.matchAll,
			Pattern:    pattern,
			Exclude:    params.GetExcluded(),
			Include:    params.GetIncluded(),
			LimitBytes: params.limitBytes,
		},
	}

	setSignalHandlers(app)

	// begin path traversation in search of matches
	app.Run(params.workers)
}
