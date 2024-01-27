package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


func main() {
	gin.SetMode(gin.ReleaseMode)

	Server := gin.Default()

	Server.GET("/",func(ctx *gin.Context){
		ctx.String(http.StatusOK,"Hello world")
	})

	Server.Run(":8080")
}
