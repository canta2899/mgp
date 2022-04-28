package main

import (
	"os"
)

func main() {

	params := ParseArgs()

	Run(
		os.Stdout, params.workers,
		params.icase, !params.nocolor,
		params.startpath, params.pattern,
		params.GetExcludedDirs(), params.limitBytes,
	)
}
