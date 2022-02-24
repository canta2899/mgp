package main

import (
    "github.com/akamensky/argparse"
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
        Default: 2,
    })

    params.exclude = parser.StringList("e", "exclude", &argparse.Options{
        Help: "Path that will be excluded (including subpaths)",
        Required: false,
        Default: []string{},
    })

    if err := parser.Parse(os.Args); err != nil {
        return params, err
    }

    return params, nil
}
