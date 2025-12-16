package utils

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/tautcony/qart/internal"
)

func Read(path string) ([]byte, *internal.FileInfo, error) {
	p, err := filepath.Abs(path)
	log.Debug().Str("path", p).Msg("Read")
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
	log.Debug().Str("path", p).Msg("Write")
	if err != nil {
		panic(err)
	}
	return os.WriteFile(p, data, 0666)
}

func Remove(path string) error {
	p, err := filepath.Abs(path)
	log.Debug().Str("path", p).Msg("Remove")
	if err != nil {
		panic(err)
	}
	return os.Remove(p)
}
