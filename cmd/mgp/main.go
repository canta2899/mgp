package main

import (
	"os"
	"github.com/canta2899/mgp/cli"
)

func main() {
  cli.RunApp(os.Args[1:])
}
