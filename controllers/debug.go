package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/tautcony/qart/internal/qr/debug"
	"github.com/tautcony/qart/internal/utils"
	"net/http"
	"rsc.io/qr"
	"rsc.io/qr/coding"
	"strconv"
	"strings"
)

func DebugFrame(c *gin.Context) {
	version, _ := strconv.Atoi(c.Query("version"))
	scale, _ := strconv.Atoi(c.Query("scale"))
	level, _ := strconv.Atoi(c.Query("level"))
	dots := c.Query("dots")

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
	c.Data(http.StatusOK, "image/png", data)
}

func DebugMask(c *gin.Context) {
	version, _ := strconv.Atoi(c.Query("version"))
	scale, _ := strconv.Atoi(c.Query("scale"))
	level, _ := strconv.Atoi(c.Query("level"))
	mask, _ := strconv.Atoi(c.Query("mask"))

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
	c.Data(http.StatusOK, "image/png", data)
}

func DebugEncode(c *gin.Context) {
	version, _ := strconv.Atoi(c.Query("version"))
	scale, _ := strconv.Atoi(c.Query("scale"))
	level, _ := strconv.Atoi(c.Query("level"))
	mask, _ := strconv.Atoi(c.Query("mask"))
	content := coding.String(c.Query("content"))

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

	c.Data(http.StatusOK, "image/png", code.PNG())
}
