package repository

import "gorm.io/gorm"

type PasswordHistoryRepository interface {
}

type passwordHistoryRepository struct {
	db gorm.DB
}

func NewPasswordHistoryRepository(db gorm.DB) PasswordHistoryRepository {
	return &passwordHistoryRepository{
		db: db,
	}
}
