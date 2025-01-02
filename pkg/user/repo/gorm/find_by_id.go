package gorm

import (
	"context"
	"fmt"

	"okj/lib/otel"
	"okj/pkg/user"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func (db *DB) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	span := trace.SpanFromContext(ctx)

	var model UserModel
	result := db.First(&model, "id = ?", id.String())
	if result.Error != nil {
		switch result.Error {
		case gorm.ErrRecordNotFound:
			return nil, user.ErrNotFoundByID
		default:
			span.AddEvent("db query failed")
			db.logger.WarnContext(ctx, otel.FormatLog(Path, "find_by_id.go [FindByID]: failed to query for user", result.Error))
			return nil, user.ErrInternal
		}
	}
	db.logger.InfoContext(ctx, otel.FormatLog(Path, fmt.Sprintf("find_by_id.go [FindByID]: found user with id %q", model.ID.String()), nil))
	span.AddEvent(fmt.Sprintf("db query returned user_id %q", id.String()))

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
