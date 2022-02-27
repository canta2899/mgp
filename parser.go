package main

import (
    "github.com/akamensky/argparse"
    "runtime"
    "os"
)

type Parameters struct {
    startpath *string
    pattern   *string
    workers   *int
    exclude   *[]string
}

func ParseArgs() (*Parameters, error) {

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

    if err := parser.Parse(os.Args); err != nil {
        return params, err
    }

    return params, nil
}
