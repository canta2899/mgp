package cli

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	app "github.com/canta2899/mgp/application"
	"github.com/canta2899/mgp/model"
	"github.com/canta2899/mgp/pkg/fspathwalk"
	"github.com/canta2899/mgp/pkg/output"
)

func compileRegex(pattern string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		// (?i) is case insensitivie notation in go regexp
		pattern = "(?i)" + pattern
	}

	return regexp.Compile(pattern)
}

// Sets a SIGTERM handler in order to stop when Ctrl+C is pressed
func setSignalHandlers(app *app.Application) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch
		app.Stop()
	}()
}

// Runs mgp from command line interface
func RunApp(args []string) {
	params := ParseArgs(args)

	if _, err := os.Stat(params.startpath); os.IsNotExist(err) {
		log.Println("Path does not exists")
		os.Exit(1)
	}

	outputHandler := output.NewFmtOutputHandler(!params.raw, params.showCtx)
	pathWalker := fspathwalk.NewFsPathWalk(params.startpath)

	pattern, err := compileRegex(params.pattern, params.icase)

	if err != nil {
		log.Fatalf(err.Error())
	}

	app := &app.Application{
		Msg:      outputHandler,
		Explorer: pathWalker,
		Options: &model.Options{
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
