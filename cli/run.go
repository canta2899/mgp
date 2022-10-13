package cli

import (
  "errors"
  "log"
  "os"
  "os/signal"
  "regexp"
  "sync"
  "syscall"
  env "github.com/canta2899/mgp/env"
  "github.com/canta2899/mgp/output"
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

func RunApp() {
  params := ParseArgs()
  handler := output.NewFmtOutputHandler(true)
  pattern, err := compileRegex(params.pattern, params.icase)

  if err != nil {
    log.Fatalf(err.Error())
  }

  stopWalk := false

  setSignalHandlers(&stopWalk)

  env := &env.Env{
    Wg:           sync.WaitGroup{},
    Schan:        make(chan bool, params.workers),
    Msg:          handler,
    MatchContext: params.matchContext,
    Pattern:      pattern,
    StopWalk:     &stopWalk,
    StartPath:    params.startpath,
    Exclude:      params.GetExcludedDirs(),
    LimitBytes:   params.limitBytes,
  }

  env.Run()
}
