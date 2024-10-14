package repository

import (
	"fmt"
	"go-socmed/entity"

	"gorm.io/gorm"
)

type AuthRepository interface {
	EmailExist(email string) bool
	Register(reqUser *entity.User) error
}

type authRepository struct {
	DB *gorm.DB
}

// constructor
func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{
		DB: db,
	}
}

func (r *authRepository) Register(reqUser *entity.User) error {
	err := r.DB.Create(&reqUser).Error

	return err
}

func (r *authRepository) EmailExist(email string) bool {
	var user *entity.User
	err := r.DB.First(&user, "email=?", email).Error

	fmt.Printf("Error Email Exist : \n %v", err)

	return err == nil
}
