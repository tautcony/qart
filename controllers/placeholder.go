package controllers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type PlaceHolder struct {
	Width  int
	Height int
	Random string
	Title  string
}

var SvgTemplate = `<svg width="{{.Width}}" height="{{.Height}}" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {{.Width}} {{.Height}}" preserveAspectRatio="none">
	<defs>
		<style type="text/css">#holder_{{.Random}} text { fill:rgba(255,255,255,.75);font-weight:normal;font-family:Helvetica, monospace;font-size:{{fontsize .}}pt; padding: 10% } </style>
	</defs>
	<g id="holder_{{.Random}}">
		<rect width="{{.Width}}" height="{{.Height}}" fill="#777"></rect>
		<g>
			<text x="0" y="{{fontsize .}}">{{.Title}}</text>
		</g>
	</g>
</svg>`

func GetTemplate() (*template.Template, error) {
	fm := template.FuncMap{"fontsize": func(ph *PlaceHolder) int {
		if ph.Width/ph.Height > 10 {
			return ph.Height / 10
		}
		return ph.Width / 20
	}}
	tpl, err := template.New("svg").Funcs(fm).Parse(SvgTemplate)
	return tpl, err
}

func PlaceHolderHandler(c *gin.Context) {
	width := 0
	height := 0
	var err error
	size := c.Param("size")
	title := c.Param("title")
	
	if size != "" {
		seps := strings.Split(size, "x")
		if len(seps) == 2 {
			width, err = strconv.Atoi(seps[0])
			if err != nil {
				width = 0
			}
			height, err = strconv.Atoi(seps[1])
			if err != nil {
				height = 0
			}
		}
	}
	if width == 0 && height == 0 {
		width = 200
		height = 200
	}
	if title == "" {
		title = fmt.Sprintf("%vx%v", width, height)
	}
	placeHolder := &PlaceHolder{
		Width:  width,
		Height: height,
		Random: strconv.Itoa(int(rand.Int31())),
		Title:  title,
	}
	tpl, err := GetTemplate()
	if err != nil {
		panic(err)
	}
	var svg bytes.Buffer
	err = tpl.Execute(&svg, placeHolder)
	if err != nil {
		panic(err)
	}
	c.Data(http.StatusOK, "image/svg+xml", svg.Bytes())
}
