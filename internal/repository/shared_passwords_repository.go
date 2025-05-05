package repository

import "gorm.io/gorm"

type SharedPasswordRepository interface {
}

type sharedPasswordRepository struct {
	db gorm.DB
}

func NewSharedPasswordRepository(db gorm.DB) SharedPasswordRepository {
	return &sharedPasswordRepository{
		db: db,
	}
}
