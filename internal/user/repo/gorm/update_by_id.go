package gorm

import (
	"context"
	"fmt"

	"okj/internal/user"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (db *DB) UpdateByID(ctx context.Context, id uuid.UUID, u *user.User) error {
	model := UserModel{ID: id}
	result := db.Model(&model).Clauses(clause.Returning{}).Updates(UserModel{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
	})
	if result.Error != nil {
		fmt.Println(result.Error)
		return user.ErrInternal
	}

	u.ID = model.ID
	u.FirstName = model.FirstName
	u.LastName = model.LastName
	u.CreatedAt = model.CreatedAt
	u.Role = model.Role
	u.UpdatedAt = model.UpdatedAt

	return nil
}
