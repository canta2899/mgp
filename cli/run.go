package cli

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"

	"github.com/canta2899/mgp/model"
	"github.com/canta2899/mgp/output"
	"github.com/canta2899/mgp/traverse"
)

func compileRegex(pattern string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		// (?i) is case insensitivie notation in go regexp
		pattern = "(?i)" + pattern
	}

	return regexp.Compile(pattern)
}

// Sets a SIGTERM handler in order to stop when Ctrl+C is pressed
func setSignalHandlers(config *model.Config) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch
		*config.StopWalk = true
	}()
}

// Runs mgp from command line interface
func RunApp() {
	// parses cli input
	params := ParseArgs()

	// cli handler for output relatex stuff
	handler := output.NewFmtOutputHandler(!params.raw, params.showCtx)

	pattern, err := compileRegex(params.pattern, params.icase)

	if err != nil {
		log.Fatalf(err.Error())
	}

	stopWalk := false

	// the config object will be used by the core application
	// in order to get handlers, flags, etc.
	env := &model.Config{
		Wg:         sync.WaitGroup{},
		Schan:      make(chan bool, params.workers),
		Msg:        handler,
		MatchAll:   params.matchAll,
		Pattern:    pattern,
		StopWalk:   &stopWalk,
		StartPath:  params.startpath,
		Exclude:    params.GetExcludedDirs(),
		LimitBytes: params.limitBytes,
	}

	setSignalHandlers(env)

	// begin path traversation in search of matches
	traverse.TraversePath(env)
}
