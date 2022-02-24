package main

import (
    "fmt"
    "path/filepath"
)

type ExcludedDirs []string

func (ed *ExcludedDirs) String() string {
    return fmt.Sprintln(*ed)
}

func (ed *ExcludedDirs) Set(s string) error {
    fullpath, err := filepath.Abs(s)
    if err != nil {
        return err
    }
    *ed = append(*ed, fullpath)
    return nil
}

