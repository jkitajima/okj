package gorm

import (
	"okj/internal/user"

	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewRepo(db *gorm.DB) user.Repoer {
	return &DB{db}
}
