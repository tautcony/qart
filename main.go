package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/tautcony/qart/routers"
)

func main() {
	prod := flag.Bool("prod", false, "run in production mode")
	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	if *prod {
		gin.SetMode(gin.ReleaseMode)
	}

	r := routers.SetupRouter()
	
	if err := r.Run(":" + *port); err != nil {
		panic(err)
	}
}
