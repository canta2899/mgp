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

func (pt *FsPathWalk) Walk(f filepath.WalkFunc) error {
	callback := (func(pathname string, info os.FileInfo, err error) error)(f)
	return filepath.Walk(pt.StartPath, callback)
}
