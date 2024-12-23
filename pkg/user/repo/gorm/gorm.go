package gorm

import (
	"log/slog"

	"okj/pkg/user"

	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
	logger *slog.Logger
}

func NewRepo(db *gorm.DB, logger *slog.Logger) user.Repoer {
	return &DB{db, logger}
}
