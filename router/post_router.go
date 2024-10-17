package router

import (
	"go-socmed/config"
	"go-socmed/handler"
	"go-socmed/repository"
	"go-socmed/service"

	"github.com/gin-gonic/gin"
)

func PostRouter(api *gin.RouterGroup) {
	postRepository := repository.NewPostRepository(config.DB)
	postService := service.NewPostService(postRepository)
	postHandler := handler.NewPostHandler(postService)

	r := api.Group("/tweets")
	r.POST("/", postHandler.Create)
}
