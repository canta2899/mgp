package pathwalk

import (
	"os"
	"path/filepath"
)

type PathTraverser struct {
	StartPath string
}

func NewPathTraverser(sp string) *PathTraverser {
	return &PathTraverser{
		StartPath: sp,
	}
}

func (pt *PathTraverser) Walk(f func(pathname string, info os.FileInfo, err error) error) error {
	return filepath.Walk(pt.StartPath, f)
}
