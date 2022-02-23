# Multigrep

A dumb CLI tool made in Go that performs the equivalent of 

```sh

grep -r -l "mypattern" /my/starting/path

```

withouth requiring **17** years to complete.

## Installation

Clone this repository and compile the `main.go` by running

```sh

go build . -o mgrep

```

Then, put the executable on your path or symlink it as follows

```sh

ln -s mgrep /usr/local/bin/mgrep

```

If everything works, running `mgrep` should prompt you a brief usage guide.
