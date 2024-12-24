package gorm

import (
	"context"
	"fmt"

	"okj/pkg/user"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func (db *DB) Insert(ctx context.Context, u *user.User) error {
	model := &UserModel{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	if u.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{
			Time:  *u.DeletedAt,
			Valid: true,
		}
	}

	result := db.Create(model)
	if result.Error != nil {
		db.logger.ErrorContext(ctx, fmt.Errorf("user: repo: gorm: %w", result.Error).Error())
		err := result.Error.(*pgconn.PgError)
		switch err.Code {
		case "23505":
			return user.ErrUserAlreadyExists
		default:
			return user.ErrInternal
		}
	}

	u.ID = model.ID
	u.Role = model.Role
	u.CreatedAt = model.CreatedAt
	u.UpdatedAt = model.UpdatedAt

	return nil
}
