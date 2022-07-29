package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"

	"github.com/fatih/color"
)

var wg sync.WaitGroup

type Entry struct {
	os.FileInfo
	Path string
}

type MessageType int64

const (
	MatchMessage MessageType = iota
	ErrorMatchMessage
	TextMessage
)

var coloredOutput bool = true

// Runes for emoji
const OK string = string('\u2713')
const KO string = string('\u00D7')

// Colors for printing
var green = color.New(color.FgHiGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

func printHandler(message string, messageType MessageType) {
	switch messageType {
	case MatchMessage:
		if coloredOutput {
			log.Printf("%v %v\n", green(OK), message)
		} else {
			log.Printf("%v\n", message)
		}
		return
	case ErrorMatchMessage:
		if coloredOutput {
			log.Printf("%v %v\n", red(KO), message)
		}
		return
	}

	log.Println(message)
}

// Routine performed by each worker
func handler(ch <-chan *Entry, closech <-chan bool, wg *sync.WaitGroup, r *regexp.Regexp) {
	defer wg.Done()

	for {

		select {
		case <-closech:
			return
		case e, more := <-ch:
			if !more {
				return
			}

			if !(*e).Mode().IsRegular() {
				continue // Skips
			}

			file, err := os.Open(e.Path)

			if err != nil {
				continue // Skips
			}

			bufread := bufio.NewReader(file)

			for {
				line, err := bufread.ReadBytes('\n')

				if err == io.EOF {
					break
				}

				if r.Match(line) {
					printHandler(e.Path, MatchMessage)
					break
				}
			}
			file.Close()
		default:
			continue
		}
	}
}

// Process path and enqueues if ok for match checking
func processPath(e *Entry, exc []string, limitMb int, r *regexp.Regexp) error {
	isdir := (*e).IsDir()

	for _, n := range exc {
		fullMatch, _ := filepath.Match(n, e.Path)
		baseMatch, _ := filepath.Match(n, filepath.Base(e.Path))
		if isdir && (fullMatch || baseMatch) {
			return filepath.SkipDir
		}
	}

	if !isdir && (*e).Size() < int64(limitMb) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !(*e).Mode().IsRegular() {
				return // Skips
			}

			file, err := os.Open(e.Path)

			if err != nil {
				return // Skips
			}

			bufread := bufio.NewReader(file)

			for {
				line, err := bufread.ReadBytes('\n')

				if err == io.EOF {
					break
				}

				if r.Match(line) {
					printHandler(e.Path, MatchMessage)
					break
				}
			}
			file.Close()
		}()
	}

	return nil
}

// Evaluates error for path and returns action to perform
func handlePathError(e *Entry, err error) error {

	if os.IsNotExist(err) {
		log.Fatal("Invalid path")
	}

	// Prints error line for current path
	printHandler(e.Path, ErrorMatchMessage)
	printHandler(err.Error(), TextMessage)

	if (*e).IsDir() {
		return filepath.SkipDir
	} else {
		return nil
	}
}

// Handler for sigterm (ctrl + c from cli)
func setSignalHandlers(closed chan bool, workers int, stopWalk *bool, wg *sync.WaitGroup) {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch

		*stopWalk = true

		for i := 0; i < workers; i++ {
			closed <- true
		}

	}()
}

func Run(out io.Writer, workers int,
	caseInsensitive bool, colors bool,
	startpath string, pattern string,
	exludedDirs []string, limitMb int) {

	// Configuring logger
	log.SetFlags(0)
	log.SetOutput(out)

	// Output with symbols and colors
	coloredOutput = colors

	// Regex compilation
	if caseInsensitive {
		pattern = "(?i)" + pattern
	}
	r, err := regexp.Compile(pattern)

	if err != nil {
		printHandler("Error in regex pattern", TextMessage)
		os.Exit(1)
	}

	// Tools for synchronization
	// var wg sync.WaitGroup
	stopWalk := false
	closeSignalChan := make(chan bool, workers)
	// ch := make(chan *Entry, 5000)
	setSignalHandlers(closeSignalChan, workers, &stopWalk, &wg)

	// Spawning routines
	// wg.Add(workers)
	// for i := 0; i < workers; i++ {
	// 	go handler(ch, closeSignalChan, &wg, r)
	// }

	// Traversing filepath
	filepath.Walk(startpath,

		func(pathname string, info os.FileInfo, err error) error {

			if stopWalk {
				// If the termination is requested, the path Walking
				// stops and the function returns with an error
				return errors.New("user requested termination")
			}

			e := &Entry{
				FileInfo: info,
				Path:     pathname,
			}

			// Checking permission and access errors
			if err != nil {
				return handlePathError(e, err)
			}

			// Processes path in search of matches with the given
			// pattern or the excluded directories
			return processPath(e, exludedDirs, limitMb, r)

		})

	// The channel is closed, this communicates that
	// no more values will be enqueued
	// close(ch)

	// Waits for goroutines to finish
	wg.Wait()

	// Ensures signal chan
	close(closeSignalChan)

	if stopWalk {
		printHandler("Ended by user", TextMessage)
	}
}
