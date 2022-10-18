package fspathwalk

import (
	"os"
	"path/filepath"
)

type FsPathWalk struct {
	StartPath string
}

func NewFsPathWalk(sp string) *FsPathWalk {
	return &FsPathWalk{
		StartPath: sp,
	}
}

func (pt *FsPathWalk) Walk(f func(pathname string, info os.FileInfo, err error) error) error {
	return filepath.Walk(pt.StartPath, f)
}
