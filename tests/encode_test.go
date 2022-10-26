package test

import (
	"github.com/tautcony/qart/internal/qr"
	"github.com/tautcony/qart/internal/utils"
	"image"
	"os"
	"path/filepath"
	"rsc.io/qr/coding"
	"testing"
)

var (
	version    coding.Version = 6
	size                      = 4
	imageFile  image.Image
	targetData [][]byte
)

func init() {
	f, err := os.Open(filepath.Join("..", "assets", "default.png"))
	if err != nil {
		panic(err)
	}
	defer func() {
		ferr := f.Close()
		if ferr != nil {
			panic(ferr)
		}
	}()
	imageFile, err = utils.GetImageThumbnail(f)
	if err != nil {
		imageFile = nil
		panic(err)
	}
	targetData = qr.MakeTarget(imageFile, 17+4*int(version)+size)
}

func BenchmarkMakeTarget(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = qr.MakeTarget(imageFile, 17+4*int(version)+size)
	}
}

func BenchmarkEncode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		qrImage := &qr.Image{
			Name:         "",
			Dx:           4,
			Dy:           4,
			URL:          "https://example.com",
			Version:      version,
			Mask:         2,
			Level:        coding.L,
			RandControl:  false,
			Dither:       false,
			OnlyDataBits: false,
			SaveControl:  false,
			Scale:        4,
			Seed:         -1366185600000,
			Rotation:     0,
			Size:         size,
			Target:       targetData,
		}

		if err := qrImage.Encode(); err != nil {
			panic(err)
		}
		_ = qrImage.Code.PNG()
	}
}
