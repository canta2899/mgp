package main

import (
	"errors"
	"log"
	"os"
	"regexp"
	"sync"
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

func main() {

	params := ParseArgs()
	handler := NewMessageHandler(!params.nocolor, os.Stdout)
	pattern, err := compileRegex(params.pattern, params.icase)

	if err != nil {
		log.Fatalf(err.Error())
	}

	env := &env{
		wg:      sync.WaitGroup{},
		sChan:   make(chan bool),
		msg:     handler,
		params:  params,
		pattern: pattern,
	}

	env.Run()
}
