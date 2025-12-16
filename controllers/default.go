package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	// Get language data from context (set by i18n middleware)
	lang, _ := c.Get("Lang")
	curLang, _ := c.Get("CurLang")
	restLangs, _ := c.Get("RestLangs")
	i18nFunc, _ := c.Get("i18n")

	c.HTML(http.StatusOK, "index.tpl", gin.H{
		"Lang":      lang,
		"CurLang":   curLang,
		"RestLangs": restLangs,
		"i18n":      i18nFunc,
	})
}
