package model

import (
  "regexp"
  "sync"
)

type Config struct {
  Wg           sync.WaitGroup
  Schan        chan bool
  StopWalk     *bool
  MatchContext bool
  Msg          OutputHandler
  Pattern      *regexp.Regexp
  StartPath    string
  Exclude      []string
  LimitBytes   int
}
