package qr

import (
	"github.com/nfnt/resize"
	"github.com/tautcony/qart/models/request"
	"image"
	"image/color"
	_ "image/png"
	"rsc.io/qr/coding"
)

func Draw(op *request.Operation, i image.Image) (*Image, error) {
	target := MakeTarget(i, 17+4*int(op.Version)+op.Size)

	qrImage := &Image{
		Name:         op.Image,
		Dx:           op.Dx,
		Dy:           op.Dy,
		URL:          op.URL,
		Version:      op.GetVersion(),
		Mask:         op.GetMask(),
		Level:        coding.L,
		RandControl:  op.RandControl,
		Dither:       op.Dither,
		OnlyDataBits: op.OnlyDataBits,
		SaveControl:  op.SaveControl,
		Scale:        op.GetScale(),
		Target:       target,
		Seed:         op.GetSeed(),
		Rotation:     op.GetRotation(),
		Size:         op.Size,
	}

	if err := qrImage.Encode(); err != nil {
		return nil, err
	}
	return qrImage, nil
}

func MakeTarget(i image.Image, max int) [][]byte {
	b := i.Bounds()
	dx, dy := max, max
	if b.Dx() > b.Dy() {
		dy = b.Dy() * dx / b.Dx()
	} else {
		dx = b.Dx() * dy / b.Dy()
	}
	thumbnail := resize.Resize(uint(dx), uint(dy), i, resize.Bilinear)

	b = thumbnail.Bounds()
	dx, dy = b.Dx(), b.Dy()
	target := make([][]byte, dy)
	arr := make([]byte, dx*dy)
	for y := 0; y < dy; y++ {
		target[y], arr = arr[:dx], arr[dx:]
		row := target[y]
		for x := 0; x < dx; x++ {
			p := thumbnail.At(x, y)
			luminance := color.Gray16Model.Convert(p).(color.Gray16)
			row[x] = byte(luminance.Y >> 8)
		}
	}
	return target
}
