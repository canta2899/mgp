package model

import "os"

type PathExplorer interface {
	Walk(func(pathname string, info os.FileInfo, err error) error) error
}
