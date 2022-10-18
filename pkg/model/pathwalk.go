package model

import "os"

type PathWalk interface {
	Walk(func(pathname string, info os.FileInfo, err error) error) error
}
