package routers

import (
	"github.com/astaxie/beego"
	"github.com/tautcony/qart/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/render", &controllers.RenderController{}),
		beego.NSRouter("/render/config", &controllers.RenderController{}, "get:Config"),
		beego.NSRouter("/render/upload", &controllers.UploadController{}),
		beego.NSRouter("/share", &controllers.ShareController{}, "post:CreateShare"),
		beego.NSRouter("/debug/frame", &controllers.DebugController{}, "get:Frame"),
		beego.NSRouter("/debug/mask", &controllers.DebugController{}, "get:Mask"),
		beego.NSRouter("/debug/encode", &controllers.DebugController{}, "get:Encode"),
	)
	beego.AddNamespace(ns)
}
