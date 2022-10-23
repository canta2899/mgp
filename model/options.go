package model

import "regexp"

type Options struct {
	MatchAll   bool
	Pattern    *regexp.Regexp
	Exclude    []string
	Include    []string
	LimitBytes int
}
