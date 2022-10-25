package cli

import (
	"errors"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const VERSION string = "v1.4.1"
const MEGABYTE int = 1048576
const PROG_NAME = "mgp"

var STD_EXC_DIRS = []string{".bzr", "CVS", ".git", ".hg", ".svn", ".idea", ".tox", "node_modules"}

type Flags struct {
	workers    int
	raw        bool
	icase      bool
	exclude    string
	include    string
	limitBytes int
	matchAll   bool
	showCtx    bool
}

func (f *Flags) GetExcluded() []string {
	if f.exclude == "" {
		return STD_EXC_DIRS
	}

	return append(STD_EXC_DIRS, strings.Split(f.exclude, ",")...)
}

func (f *Flags) GetIncluded() []string {
	inc := []string{}

	if f.include == "" {
		return inc
	}

	return append(inc, strings.Split(f.include, ",")...)
}

type Parameters struct {
	Flags
	startpath string
	pattern   string
}

func PrintBriefHelpAndExit() {
	log.Println(PROG_NAME, VERSION)
	log.Println("Usage:", PROG_NAME, "[options] pattern starting/path")
	log.Println("Run", PROG_NAME, "-h for more information")
	os.Exit(0)
}

func parseByteLimit(limit string) (int, error) {
  limitExp, err := regexp.Compile("([0-9]+)\\s*(b|kb|mb|gb)?")
  
  if err != nil {
    return 0, err
  }

  if limitExp.MatchString(limit) {
    sub := limitExp.FindStringSubmatch(limit)

    size := sub[1]
    uom := sub[2]
    
    finalSize, _ := strconv.Atoi(size) 

    // b is default

    switch uom {
    case "kb":
      finalSize *= 1024 
    case "mb":
      finalSize *= 1048576
    case "gb":
      finalSize *= 1073741824
    }

    return finalSize, nil
  }

  return 0, errors.New("unable to match limit string")
}

func ParseArgs(args []string) *Parameters {

	f := Flags{}
  limit := ""

	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	printVersion := false

	flag.IntVar(&f.workers, "w", 1000, "Defines the max number of routines running at the same time")
	flag.BoolVar(&printVersion, "v", false, "Prints current mgp version")
  flag.StringVar(&limit, "lim", "5gb", "File size limit (for example 2kb or 10gb)")
	flag.BoolVar(&f.icase, "i", false, "Performs case insensitive matching")
	flag.BoolVar(&f.raw, "raw", false, "Disable colored output")
	flag.StringVar(&f.exclude, "exc", "", "Excluded paths (specified as a comma separated list like \"path1,path2\")")
	flag.StringVar(&f.include, "inc", "", "Included paths (specified as a comma separated list like \"path1,path2\")")
	flag.BoolVar(&f.matchAll, "all", false, "Show every match for a file")
	flag.BoolVar(&f.showCtx, "ctx", false, "Print match context")

	flag.Parse()
  flag.CommandLine.Parse(args)

	posArgs := flag.Args()

	if printVersion || len(posArgs) < 2 {
		PrintBriefHelpAndExit()
	}

  lim, err := parseByteLimit(limit)

  f.limitBytes = lim

  if err != nil {
    log.Println("Error parsing byte limit, the option will be ignored")
  }


	return &Parameters{
		Flags:     f,
		pattern:   posArgs[0],
		startpath: posArgs[1],
	}
}
