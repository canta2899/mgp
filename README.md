![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Version](https://img.shields.io/github/v/release/canta2899/multigrep?display_name=tag&label=version&style=for-the-badge)
![WrittenInGo](https://img.shields.io/badge/Written%20in%20Go-lightblue?style=for-the-badge&logo=go&color=111111)

<p align="center">
    <img src="./assets/multigrep.gif" width="700"/>
</p>
<h1 align="center">
    MGP
</h1>

A dumb and platform independent command line tool made in Go that performs the equivalent of the following **grep** command, but without taking **ages** to complete.

```sh

grep -E -r -l "mypattern" "/my/path"

```

## Why

I always end up having to search for files containing a simple keyword inside a multitude of directories and subdirectories. Most of the times, grep does its job. Others it takes so long that I end up killing the process.

I made this tool in order to perform the same kind of lookup while taking advantage of **goroutines** in order to parallelize the research in files while the path is being traversed.

Grep still remains the best tool, but for specific needs **MGP** may come handy too.

## How it works

The implementation follows a simple producer/consumer pattern in which a single goroutine traverses the given directory recursively adding all the valid paths to a queue. Meanwhile, a series of parallel goroutines (whose number is proportional to the amount of logical CPUs available) dequeues each path concurrently and searches for a match between its content and the pattern provided. The queue implementation in based on go's buffered channels.

## Usage

Two parameters are required

- The **pattern** that needs to be matched
- The starting **path** for the recursive research

These can be specified as positional arguments like in grep. Moreover, additional flags can be specified before the pattern and the starting path. These allow to: 

- Exclude specific path or directories `-e "path1,path2,path3"` 
- Specify a size limit (in Megabytes) in order to exclude big files `-l 800`
- Specify a number of workers in order to change the degree of parallelism `-w 16`
- Disable the colored output `-c`
- Perform case insensitive matching `-i`

### Examples

Here's an example that searches for the word *Panda* recursively starting from the current directory and ignoring directories named *not-me* at any level.

```sh
mgp -e "not-me" Panda . 
```

Here's, instead, an example that searches for the word *Node* and the word *node* recursively starting from the */home/user/* path and specifically ignoring the */home/user/.local/bin* directory and directories named *.git* at any level.

```sh
mgp -e ".git,/home/user.local/bin" "[Nn]ode" /home/user/ 
```

<p align="center">
    <h6 align="center">Pretty easy isn't it?</h6>
</p>


Running `multigrep -h` or `multigrep --help` will prompt a usage guide too.

## Installation

### Binaries

Precompiled binaries are available in the **Releases** section of this repository. Once downloaded (let's say, for example, I've downloaded the *mgp-v1.1.0-darwin-amd64.tar.gz* archive), one can run

```sh
tar -xzf mgp-v1.0.0-darwin-amd64.tar.gz
```

This will extract the **executable** and a text file containting the **license**. You can, then, place the binary file in your path (or symlink it). Running `mgp -v` should, then, prompt a message stating the current version.

### Source code

You can also download **MGP** as a Go module. You'll have to install the Go distribution for your system and then run

```sh
go install github.com/canta2899/mgp@latest
```

This will download and compile the program for your platform.



