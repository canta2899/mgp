package main

import (
  "regexp"
  "sync"
)

type Env struct {
  wg         sync.WaitGroup
  sChan      chan bool
  stopWalk   *bool
  msg        OutputHandler
  pattern    *regexp.Regexp
  startpath  string
  exclude    []string
  limitBytes int
}
