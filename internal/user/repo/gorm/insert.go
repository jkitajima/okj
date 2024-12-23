package gorm

import (
	"context"
	"fmt"

	"okj/internal/user"

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
		fmt.Println(result.Error)
		return user.ErrInternal
	}

	u.ID = model.ID
	u.Role = model.Role
	u.CreatedAt = model.CreatedAt
	u.UpdatedAt = model.UpdatedAt

	return nil
}
