package main

import "os"

type FileInfo struct {
    os.FileInfo
    Path string
}

