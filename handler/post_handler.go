package handler

import (
	"go-socmed/dto"
	errorhandler "go-socmed/errorHandler"
	"go-socmed/helper"
	"go-socmed/service"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type postHandler struct {
	Service service.PostService
}

func NewPostHandler(s service.PostService) *postHandler {
	return &postHandler{
		Service: s,
	}
}

func (h *postHandler) Create(c *gin.Context) {
	var post dto.PostRequest

	if err := c.ShouldBind(&post); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	// cek apakah mengirim file picture
	if post.Picture != nil {
		if err := os.MkdirAll("/public/picture", 0755); err != nil {
			errorhandler.HandleError(c, &errorhandler.InternalServerError{Message: err.Error()})
			return
		}

		// rename picture
		ext := filepath.Ext(post.Picture.Filename)
		newFileName := uuid.New().String() + ext

		// save image to directory
		dst := filepath.Join("public/picture", filepath.Base(newFileName))
		c.SaveUploadedFile(post.Picture, dst)

		post.Picture.Filename = newFileName
	}

	userId := 1
	post.UserID = userId

	if err := h.Service.Create(&post); err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "success post your tweet",
	})

	c.JSON(http.StatusCreated, res)

}
