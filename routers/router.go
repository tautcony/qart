// @APIVersion 1.0.0
// @Title QArt API
// @Description TO BE FILLED
package routers

import (
	"github.com/astaxie/beego"
	"github.com/tautcony/qart/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/image/placeholder/:size/?:title", &controllers.PlaceHolderController{})
	beego.Router("/share/?:sha", &controllers.ShareController{})
}
