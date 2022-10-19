package fspathwalk

import (
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
	return filepath.Walk(pt.StartPath, f)
}
