
<p align="center">
    <img src="./assets/multigrep.gif" width="700"/>
</p>
<h1 align="center">
    Multigrep
</h1>

A dumb and platform independent command line tool made in Go that performs the equivalent of the following **grep** command, but without taking **ages** to complete.

```sh

grep -E -r -l "mypattern" "/my/path"

```

## Why

I always end up having to search for files containing a simple keyword inside a multitude of directories and subdirectories. Most of the times, grep does its job. Others it takes so long that I end up killing the process.

I made this tool in order to perform the same kind of lookup while taking advantage of **goroutines** in order to parallelize the research in files while the path is being traversed.

Grep still remains the best tool, but for specific needs **multigrep** may come handy too.

### Note

Please keep in mind that files will be read to memory while being examined, so exclude big files with the `-e` flag in order to avoid saturating your RAM.

## How it works

The program follows a simple producer/consumer pattern in which the main thread enqueues all the valid paths while other previously spawned goroutines (2 by default, but you can choose the number when running the command) process them. This is done by dequeuing a path, opening the respective file and searching for a match in its content.

The queue implementation in based on a thread-safe linked list and follows a standard FIFO policy. Despite **channels** are recommended (and much easier to use), it is really difficult to estimate a good buffer size in order to retain performances and ensure that it will be enough. This is the best choice I came up with but it's open to any possible improvement.

## Installation

Install **go** and run 

```sh

go get github.com/canta2899/multigrep

```

This will download and compile the program for your platform. To check if everything works, run 

```sh
multigrep --help
```


