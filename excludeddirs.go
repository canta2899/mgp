package main

import (
    "fmt"
)

type ExcludedDirs []string

func (ed *ExcludedDirs) String() string {
    return fmt.Sprintln(*ed)
}

func (ed *ExcludedDirs) Set(s string) error {
    *ed = append(*ed, s)
    return nil
}

