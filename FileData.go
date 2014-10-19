package main

import (
	"os"
)

type FileData struct {
	os.FileInfo
	path string
}

func (fd *FileData) Path() string {
	return fd.path
}
