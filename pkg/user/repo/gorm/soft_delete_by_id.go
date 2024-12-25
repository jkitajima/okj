package gorm

import (
	"context"

	"okj/lib/otel"
	"okj/pkg/user"

	"github.com/google/uuid"
)

func (db *DB) SoftDeleteByID(ctx context.Context, id uuid.UUID) error {
	model := UserModel{ID: id}
	result := db.Delete(&model)
	if result.Error != nil {
		db.logger.WarnContext(ctx, otel.FormatLog(Path, "soft_delete_by_id.go [SoftDeleteByID]: failed to soft delete user", result.Error))
		return user.ErrInternal
	}

	return nil
}
