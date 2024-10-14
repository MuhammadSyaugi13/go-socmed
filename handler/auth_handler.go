package handler

import (
	"go-socmed/dto"
	errorhandler "go-socmed/errorHandler"
	"go-socmed/helper"
	"go-socmed/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *authHandler {
	return &authHandler{
		service: s,
	}
}

func (h *authHandler) Register(c *gin.Context) {
	var register dto.RegisterRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := h.service.Register(&register); err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "Register Successfuly, please register!",
	})

	c.JSON(http.StatusCreated, res)
}

func (h *authHandler) Login(c *gin.Context) {

	var login dto.LoginRequest

	//  mengambil data request dan memasukan kedalam struct LoginRequest
	err := c.ShouldBindJSON(&login)
	if err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	result, err := h.service.Login(&login)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Login successfuly",
		Data:       result,
	})

	c.JSON(http.StatusOK, res)

}
