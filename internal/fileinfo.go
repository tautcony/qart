package internal

import "time"

type FileInfo struct {
	Name    string // final path element
	ModTime time.Time
	Size    int64
	IsDir   bool
}
