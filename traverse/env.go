package traverse

import (
  "regexp"
  "sync"
  "github.com/canta2899/mgp/output"
)

type Env struct {
  Wg           sync.WaitGroup
  Schan        chan bool
  StopWalk     *bool
  MatchContext bool
  Msg          output.OutputHandler
  Pattern      *regexp.Regexp
  StartPath    string
  Exclude      []string
  LimitBytes   int
}
