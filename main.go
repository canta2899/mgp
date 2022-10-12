package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
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

func setSignalHandlers(stopWalk *bool) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch
		*stopWalk = true
	}()
}

func main() {

	params := ParseArgs()
	// handler := NewMessageHandler(!params.nocolor, os.Stdout)
  handler := NewFmtOutputHandler(true)
	pattern, err := compileRegex(params.pattern, params.icase)

	if err != nil {
		log.Fatalf(err.Error())
	}

	stopWalk := false

	setSignalHandlers(&stopWalk)

	env := &Env{
		wg:         sync.WaitGroup{},
		sChan:      make(chan bool, params.workers),
		msg:        handler,
		pattern:    pattern,
		stopWalk:   &stopWalk,
		startpath:  params.startpath,
		exclude:    params.GetExcludedDirs(),
		limitBytes: params.limitBytes,
	}

	env.Run()
}
