package model

import "regexp"

type Options struct {
	Running    chan bool
	StopWalk   chan bool
	MatchAll   bool
	Pattern    *regexp.Regexp
	Exclude    []string
	Include    []string
	LimitBytes int
}
