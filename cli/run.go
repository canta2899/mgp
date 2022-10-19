package cli

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"

	"github.com/canta2899/mgp/internal/fspathwalk"
	"github.com/canta2899/mgp/internal/output"
	app "github.com/canta2899/mgp/application"
	"github.com/canta2899/mgp/model"
)

func compileRegex(pattern string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		// (?i) is case insensitivie notation in go regexp
		pattern = "(?i)" + pattern
	}

	return regexp.Compile(pattern)
}

// Sets a SIGTERM handler in order to stop when Ctrl+C is pressed
func setSignalHandlers(config *app.Application) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch
		config.Options.StopWalk <- true
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
	explorer := fspathwalk.NewFsPathWalk(params.startpath)

	pattern, err := compileRegex(params.pattern, params.icase)

	if err != nil {
		log.Fatalf(err.Error())
	}

	env := &app.Application{
		Wg:         sync.WaitGroup{},
		Msg:        handler,
		Explorer:   explorer,
    Options:    &model.Options{
      Running:    make(chan bool, params.workers),
      StopWalk:   make(chan bool),
      MatchAll:   params.matchAll,
      Pattern:    pattern,
      Exclude:    params.GetExcluded(),
      Include:    params.GetIncluded(),
      LimitBytes: params.limitBytes,
    },
	}

	setSignalHandlers(env)

	// begin path traversation in search of matches
	env.Run()
}
