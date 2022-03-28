package main

import (
    "github.com/akamensky/argparse"
    "runtime"
    "os"
    "log"
)

const VERSION string = "v1.2.3"
const MEGABYTE int = 1048576

type Parameters struct {
    startpath *string
    pattern   *string
    workers   *int
    nocolor   *bool
    icase     *bool
    exclude   *[]string
    limitMb   *int
}

func ParseArgs() *Parameters {

    if len(os.Args) > 1 && os.Args[1] == "--version" {
        log.Println("Multigrep:", VERSION)
        os.Exit(0)
    }

    params := &Parameters{}

    parser := argparse.NewParser("multigrep", "A command line tool to search in files recursively")

    params.pattern = parser.String("m", "match", &argparse.Options{
        Help: "A regex pattern that requires to be matched",
        Required: true,
    })

    params.startpath = parser.String("p", "path", &argparse.Options{
        Help: "The path on which the recursive search starts",
        Required: true,
    })

    params.workers = parser.Int("w", "workers", &argparse.Options{
        Help: "Number of workers, in order to define a degree of parallelism",
        Required: false,
        Default: runtime.NumCPU(),
    })

    params.exclude = parser.StringList("e", "exclude", &argparse.Options{
        Help: "Excluded files or directories",
        Required: false,
        Default: []string{"*/.bzr","*/CVS","*/.git","*/.hg","*/.svn","*/.idea","*/.tox"},
    })

    params.nocolor = parser.Flag("c", "no-color", &argparse.Options{
        Help: "Unsets colored output",
        Required: false,
        Default: false,
    })

    params.icase = parser.Flag("i", "ignore-case", &argparse.Options{
        Help: "Case insensitive match",
        Required: false,
        Default: false,
    })

    params.limitMb = parser.Int("s", "size", &argparse.Options{
        Help: "Maximum size in Megabytes for files that will be scanned",
        Required: false,
        Default: 100*MEGABYTE,
    })

    if err := parser.Parse(os.Args); err != nil {
        log.Fatal(err.Error())
    }

    return params
}
