package main

import (
	"fmt"
	"go-socmed/config"
	"go-socmed/router"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadConfig()
	config.LoadDB()

	r := gin.Default()

	api := r.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.AuthRouter(api)

	r.Run(fmt.Sprintf("localhost:%v", config.ENV.PORT)) // listen and serve on 0.0.0.0:8080
}
