package utils

import (
	"github.com/tautcony/qart/internal"
	"log"
	"os"
	"path/filepath"
)

func Read(path string) ([]byte, *internal.FileInfo, error) {
	p, err := filepath.Abs(path)
	log.Printf("Read <- %v", p)
	if err != nil {
		panic(err)
	}
	dir, err := os.Stat(p)
	if err != nil {
		return nil, nil, err
	}
	fi := &internal.FileInfo{
		Name:    dir.Name(),
		ModTime: dir.ModTime(),
		Size:    dir.Size(),
		IsDir:   dir.IsDir(),
	}
	data, err := os.ReadFile(p)
	return data, fi, err
}

func Write(path string, data []byte) error {
	p, err := filepath.Abs(path)
	log.Printf("Write ->: %v", p)
	if err != nil {
		panic(err)
	}
	return os.WriteFile(p, data, 0666)
}

func Remove(path string) error {
	p, err := filepath.Abs(path)
	log.Printf("Remove x %v", p)
	if err != nil {
		panic(err)
	}
	return os.Remove(p)
}
