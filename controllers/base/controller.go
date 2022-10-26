package base

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/beego/i18n"
	"github.com/tautcony/qart/models/response"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"path"
	"strings"
	"time"
)

var (
	AppVer string
)

type langType struct {
	Lang, Name string
}

var langTags []language.Tag // Languages are supported.

type QArtController struct {
	beego.Controller
	i18n.Locale
}

func (c *QArtController) JSON(data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}

func (c *QArtController) Success(data interface{}, code int) {
	r := &response.BaseResponse{
		Success: true,
		Code:    code,
		Data:    data,
	}
	c.JSON(r)
}

func (c *QArtController) Fail(data interface{}, code int, message string) {
	r := &response.BaseResponse{
		Code:    code,
		Data:    data,
		Message: message,
	}
	c.JSON(r)
}

// Prepare implemented Prepare method for baseRouter.
func (c *QArtController) Prepare() {
	// Setting properties.
	c.Data["AppVer"] = AppVer

	c.Data["PageStartTime"] = time.Now()

	// Redirect to make URL clean.
	if c.matchLang() {
		i := strings.Index(c.Ctx.Request.RequestURI, "?")
		c.Redirect(c.Ctx.Request.RequestURI[:i], 302)
		return
	}
}

func (c *QArtController) matchLang() bool {
	requireRedirect := false
	var matcher = language.NewMatcher(langTags)
	urlLang := c.Input().Get("lang")
	cookieLang := c.Ctx.GetCookie("lang")
	accept := c.Ctx.Request.Header.Get("Accept-Language")

	requireRedirect = urlLang != "" // language from url trigger a redirect

	curLang, _ := language.MatchStrings(matcher, urlLang, cookieLang, accept)

	// Save language information in cookies.
	if cookieLang == "" || cookieLang != curLang.String() {
		c.Ctx.SetCookie("lang", curLang.String(), 1<<31-1, "/")
	}

	restLangs := make([]*langType, 0, len(langTags)-1)
	for _, v := range langTags {
		if curLang != v {
			restLangs = append(restLangs, &langType{
				Lang: v.String(),
				Name: display.Self.Name(v),
			})
		}
	}

	// Set language properties.
	c.Lang = curLang.String()
	c.Data["Lang"] = curLang.String()
	c.Data["CurLang"] = langType{
		Lang: curLang.String(),
		Name: display.Self.Name(curLang),
	}
	c.Data["RestLangs"] = restLangs

	return requireRedirect
}

func initLocales() {
	// Initialized language type list.
	var availableLangs []string
	langConfig := beego.AppConfig.String("lang::available_lang")
	err := json.Unmarshal([]byte(langConfig), &availableLangs)
	if err != nil {
		logs.Error("Language config invalid", langConfig)
		return
	}

	langTags = make([]language.Tag, 0, len(availableLangs))
	for _, name := range availableLangs {
		l := language.Make(name)
		langTags = append(langTags, l)
	}

	for _, tag := range langTags {
		logs.Info("Loading language: %v[%v]", display.Self.Name(tag), tag.String())
		if err := i18n.SetMessage(tag.String(), path.Join("conf", "locale", fmt.Sprintf("locale_%v.ini", tag.String()))); err != nil {
			logs.Error("Fail to set message file: " + err.Error())
			return
		}
	}
}

func init() {
	initLocales()
}
