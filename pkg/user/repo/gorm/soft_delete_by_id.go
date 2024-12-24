package gorm

import (
	"context"
	"fmt"

	"okj/pkg/user"

	"github.com/google/uuid"
)

func (db *DB) SoftDeleteByID(ctx context.Context, id uuid.UUID) error {
	model := UserModel{ID: id}
	result := db.Delete(&model)
	if result.Error != nil {
		fmt.Println(result.Error)
		return user.ErrInternal
	}

	return nil
}
