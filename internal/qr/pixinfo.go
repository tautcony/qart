package qr

import (
	"rsc.io/qr/coding"
)

type Pixinfo struct {
	X        int
	Y        int
	Pix      coding.Pixel
	Targ     byte
	DTarg    int
	Contrast int
	HardZero bool
	Block    *BitBlock
	Bit      uint
}

func addDither(pixByOff []Pixinfo, pix coding.Pixel, err int) {
	if pix.Role() != coding.Data && pix.Role() != coding.Check {
		return
	}
	pinfo := &pixByOff[pix.Offset()]
	// log.Println("dither: add", pinfo.X, pinfo.Y, pinfo.DTarg, err)
	pinfo.DTarg += err
}
