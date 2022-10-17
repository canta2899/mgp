package cli

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"

	"github.com/canta2899/mgp/internal/output"
	"github.com/canta2899/mgp/internal/pathwalk"
	"github.com/canta2899/mgp/pkg/services"
)

func compileRegex(pattern string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		// (?i) is case insensitivie notation in go regexp
		pattern = "(?i)" + pattern
	}

	return regexp.Compile(pattern)
}

// Sets a SIGTERM handler in order to stop when Ctrl+C is pressed
func setSignalHandlers(config *services.Application) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch
		config.StopWalk <- true
	}()
}

// Runs mgp from command line interface
func RunApp() {
	params := ParseArgs()

	if _, err := os.Stat(params.startpath); os.IsNotExist(err) {
		log.Println("Path does not exists")
		os.Exit(1)
	}

	handler := output.NewFmtOutputHandler(!params.raw, params.showCtx)
	explorer := pathwalk.NewPathTraverser(params.startpath)

	pattern, err := compileRegex(params.pattern, params.icase)

	if err != nil {
		log.Fatalf(err.Error())
	}

	env := &services.Application{
		Wg:         sync.WaitGroup{},
		Running:    make(chan bool, params.workers),
		Msg:        handler,
		MatchAll:   params.matchAll,
		Pattern:    pattern,
		StopWalk:   make(chan bool),
		Explorer:   explorer,
		Exclude:    params.GetExcludedDirs(),
		LimitBytes: params.limitBytes,
	}

	setSignalHandlers(env)

	// begin path traversation in search of matches
	env.Run()
}
