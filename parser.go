package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

const VERSION string = "v1.3.0"
const MEGABYTE int = 1048576
const PROG_NAME = "mgp"

var STD_EXC_DIRS = []string{".bzr", "CVS", ".git", ".hg", ".svn", ".idea", ".tox"}

type Flags struct {
	workers    int
	nocolor    bool
	icase      bool
	exclude    string
	limitBytes int
}

func (f *Flags) GetExcludedDirs() []string {
	if f.exclude == "" {
		return STD_EXC_DIRS
	}

	return append(STD_EXC_DIRS, strings.Split(f.exclude, ",")...)
}

type Parameters struct {
	Flags
	startpath string
	pattern   string
}

func PrintVersionAndExit() {
	log.Println(PROG_NAME, VERSION)
	os.Exit(0)
}

func PrintBriefHelpAndExit() {
	log.Println("Usage:", PROG_NAME, "[options] pattern starting/path")
	log.Println("Run", PROG_NAME, "-h for more information")
	os.Exit(0)
}

func ParseArgs() *Parameters {

	f := Flags{}

	printVersion := false

	flag.IntVar(&f.workers, "w", 16, "Defines the number of workers")
	flag.BoolVar(&printVersion, "v", false, "Prints current mgp version")
	flag.IntVar(&f.limitBytes, "l", 100, "File size limit")
	flag.BoolVar(&f.icase, "i", false, "Performs case insensitive matching")
	flag.BoolVar(&f.nocolor, "c", false, "Disable colored output")
	flag.StringVar(&f.exclude, "e", "", "Excluded paths (specified as a comma separated list like \"path1,path2\")")

	f.limitBytes = f.limitBytes * MEGABYTE

	flag.Parse()

	posArgs := flag.Args()
	if len(posArgs) < 2 {
		PrintBriefHelpAndExit()
	}

	if printVersion {
		PrintVersionAndExit()
	}

	return &Parameters{
		Flags:     f,
		pattern:   posArgs[0],
		startpath: posArgs[1],
	}
}
