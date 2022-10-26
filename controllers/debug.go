package controllers

import (
	"github.com/tautcony/qart/controllers/base"
	"github.com/tautcony/qart/internal/qr/debug"
	"github.com/tautcony/qart/internal/utils"
	"rsc.io/qr"
	"rsc.io/qr/coding"
	"strings"
)

type DebugController struct {
	base.QArtController
}

func (c *DebugController) Frame() {
	version, _ := c.GetInt("version")
	scale, _ := c.GetInt("scale")
	level, _ := c.GetInt("level")
	dots := c.GetString("dots")

	if version == 0 {
		version = 6
	}
	if scale == 0 {
		scale = 8
	}
	if level == 0 {
		level = 6
	}

	dots = strings.ToLower(dots)
	showDots := dots != "" && dots != "nil" && dots != "null" && dots != "false"

	frame := debug.MakeFrame("", 0, coding.Version(version), coding.Level(level), coding.Mask(0), scale, showDots)
	data := utils.PngEncode(frame)
	c.Ctx.Output.ContentType(".png")
	err := c.Ctx.Output.Body(data)
	if err != nil {
		panic(err)
	}
}

func (c *DebugController) Mask() {
	version, _ := c.GetInt("version")
	scale, _ := c.GetInt("scale")
	level, _ := c.GetInt("level")
	mask, _ := c.GetInt("mask")

	if version == 0 {
		version = 6
	}
	if scale == 0 {
		scale = 8
	}
	if level == 0 {
		level = 6
	}

	frame := debug.MakeMask("", 0, coding.Version(version), coding.Level(level), coding.Mask(mask), scale)
	data := utils.PngEncode(frame)
	c.Ctx.Output.ContentType(".png")
	err := c.Ctx.Output.Body(data)
	if err != nil {
		panic(err)
	}
}

func (c *DebugController) Encode() {
	version, _ := c.GetInt("version")
	scale, _ := c.GetInt("scale")
	level, _ := c.GetInt("level")
	mask, _ := c.GetInt("mask")
	content := coding.String(c.GetString("content"))

	if version == 0 {
		version = 6
	}
	if scale == 0 {
		scale = 8
	}
	if level == 0 {
		level = 6
	}

	p, err := coding.NewPlan(coding.Version(version), coding.Level(level), coding.Mask(mask))
	if err != nil {
		panic(err)
	}
	cc, err := p.Encode(content)
	if err != nil {
		panic(err)
	}

	code := &qr.Code{Bitmap: cc.Bitmap, Size: cc.Size, Stride: cc.Stride, Scale: 8}

	c.Ctx.Output.ContentType(".png")
	err = c.Ctx.Output.Body(code.PNG())
	if err != nil {
		panic(err)
	}
}
