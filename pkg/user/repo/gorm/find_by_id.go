package gorm

import (
	"context"
	"fmt"

	"okj/pkg/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db *DB) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var model UserModel
	result := db.First(&model, "id = ?", id.String())
	if result.Error != nil {
		fmt.Println(result.Error)

		switch result.Error {
		case gorm.ErrRecordNotFound:
			return nil, user.ErrNotFoundByID
		default:
			return nil, user.ErrInternal
		}
	}

	user := user.User{
		ID:        model.ID,
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Role:      model.Role,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
	return &user, nil
}
