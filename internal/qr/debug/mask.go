package debug

import (
	"github.com/tautcony/qart/internal/utils"
	"image"
	"rsc.io/qr/coding"
)

var maskName = []string{
	"(x+y) % 2",
	"y % 2",
	"x % 3",
	"(x+y) % 3",
	"(y/2 + x/3) % 2",
	"xy%2 + xy%3",
	"(xy%2 + xy%3) % 2",
	"(xy%3 + (x+y)%2) % 2",
}

func MakeMask(font string, pt int, version coding.Version, level coding.Level, mask coding.Mask, scale int) image.Image {
	p, err := coding.NewPlan(version, level, mask)
	if err != nil {
		panic(err)
	}
	m := utils.MakeImage(maskName[mask], font, pt, len(p.Pixel), 0, scale, func(x, y int) uint32 {
		pix := p.Pixel[y][x]
		switch pix.Role() {
		case coding.Data, coding.Check:
			if pix&coding.Invert != 0 {
				return 0x000000ff
			}
		}
		return 0xffffffff
	})
	return m
}
