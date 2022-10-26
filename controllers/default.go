package controllers

import (
	"github.com/tautcony/qart/controllers/base"
)

type MainController struct {
	base.QArtController
}

func (c *MainController) Get() {
	c.TplName = "index.tpl"
}
