package main

import "os"

// Extends os.FileInfo in order to provide the path too
type FileInfo struct {
	os.FileInfo
	Path string
}
