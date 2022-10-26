package utils

import (
	"bytes"
	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"path"
)

func GetImageThumbnail(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	// Convert image to 256x256 size.
	b := img.Bounds()
	const max = 256
	dx, dy := max, max
	if b.Dx() > b.Dy() {
		dy = b.Dy() * dx / b.Dx()
	} else {
		dx = b.Dx() * dy / b.Dy()
	}
	i128 := resize.Resize(uint(dx), uint(dy), img, resize.Bicubic)

	return i128, nil
}

func MakeImage(caption, font string, pt, size, border, scale int, f func(x, y int) uint32) *image.RGBA {
	d := (size + 2*border) * scale
	csize := 0
	if caption != "" {
		if pt == 0 {
			pt = 11
		}
		csize = pt * 2
	}
	c := image.NewRGBA(image.Rect(0, 0, d, d+csize))

	// white
	u := &image.Uniform{C: color.White}
	draw.Draw(c, c.Bounds(), u, image.Point{}, draw.Src)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			r := image.Rect((x+border)*scale, (y+border)*scale, (x+border+1)*scale, (y+border+1)*scale)
			rgba := f(x, y)
			u.C = color.RGBA{R: byte(rgba >> 24), G: byte(rgba >> 16), B: byte(rgba >> 8), A: byte(rgba)}
			draw.Draw(c, r, u, image.Point{}, draw.Src)
		}
	}

	if csize != 0 {
		if font == "" {
			font = "luxisr.ttf"
		}
		font = path.Join("assets", font)
		dat, _, err := Read(font)
		if err != nil {
			panic(err)
		}
		tfont, err := freetype.ParseFont(dat)
		if err != nil {
			panic(err)
		}
		ft := freetype.NewContext()
		ft.SetDst(c)
		ft.SetDPI(100)
		ft.SetFont(tfont)
		ft.SetFontSize(float64(pt))
		ft.SetSrc(image.NewUniform(color.Black))
		ft.SetClip(image.Rect(0, 0, 0, 0))
		wid, err := ft.DrawString(caption, freetype.Pt(0, 0))
		if err != nil {
			panic(err)
		}
		p := freetype.Pt(d, d+3*pt/2)
		p.X -= wid.X
		p.X /= 2
		ft.SetClip(c.Bounds())
		_, _ = ft.DrawString(caption, p)
	}

	return c
}

func PngEncode(c image.Image) []byte {
	var b bytes.Buffer
	err := png.Encode(&b, c)
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}
