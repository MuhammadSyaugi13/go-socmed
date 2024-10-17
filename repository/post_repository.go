package repository

import (
	"go-socmed/entity"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *entity.Post) error
}

type postRepository struct {
	*gorm.DB
}

func NewPostRepository(db *gorm.DB) *postRepository {
	return &postRepository{
		DB: db,
	}
}

func (r *postRepository) Create(post *entity.Post) error {
	err := r.DB.Create(&post).Error
	return err
}
