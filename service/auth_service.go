package service

import (
	"go-socmed/dto"
	"go-socmed/entity"
	errorhandler "go-socmed/errorHandler"
	"go-socmed/helper"
	"go-socmed/repository"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
}

type authService struct {
	repository repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) *authService {
	return &authService{
		repository: r,
	}
}

func (s *authService) Register(req *dto.RegisterRequest) error {

	// cek apakah email sudah terdaftar
	if EmailExist := s.repository.EmailExist(req.Email); EmailExist {
		return &errorhandler.BadRequestError{Message: "email already registered!"}
	}

	// cek apakah password dan konfirmasi password sama
	if req.Password != req.PasswordConfirmation {
		return &errorhandler.BadRequestError{Message: "Password not match!"}
	}

	passwordHash, err := helper.HashPassword(req.Password)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: passwordHash,
		Gender:   req.Gender,
	}

	if err := s.repository.Register(&user); err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return nil
}
