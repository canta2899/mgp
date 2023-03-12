package main

import (
	"path/filepath"
)

type PathWalk struct {
	StartPath string
}

func NewPathWalk(sp string) *PathWalk {
	return &PathWalk{
		StartPath: sp,
	}
}

func (pt *PathWalk) Walk(f filepath.WalkFunc) error {
	return filepath.Walk(pt.StartPath, f)
}
