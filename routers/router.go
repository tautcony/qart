// @APIVersion 1.0.0
// @Title QArt API
// @Description TO BE FILLED
package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/tautcony/qart/controllers"
	"github.com/tautcony/qart/middleware"
	"html/template"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Setup custom template functions
	r.SetFuncMap(template.FuncMap{
		"i18n": func(lang, key string, args ...interface{}) string {
			return middleware.Tr(lang, key, args...)
		},
	})

	// Setup middleware
	r.Use(middleware.I18n())
	r.Use(middleware.Session())

	// Static files
	r.Static("/static", "./static")
	
	// HTML templates
	r.LoadHTMLGlob("views/*.tpl")

	// Main routes
	r.GET("/", controllers.Index)
	r.GET("/image/placeholder/:size", controllers.PlaceHolderHandler)
	r.GET("/image/placeholder/:size/:title", controllers.PlaceHolderHandler)
	r.GET("/share", controllers.GetShare)
	r.GET("/share/:sha", controllers.GetShare)

	// API v1 routes
	v1 := r.Group("/v1")
	{
		v1.POST("/render", controllers.Render)
		v1.GET("/render/config", controllers.RenderConfig)
		v1.POST("/render/upload", controllers.Upload)
		v1.POST("/share", controllers.CreateShare)
		
		// Debug routes
		debug := v1.Group("/debug")
		{
			debug.GET("/frame", controllers.DebugFrame)
			debug.GET("/mask", controllers.DebugMask)
			debug.GET("/encode", controllers.DebugEncode)
		}
	}

	return r
}
