package model

import "path/filepath"

type PathWalk interface {
	Walk(f filepath.WalkFunc) error
}
