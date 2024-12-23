package user

import (
	"context"

	"github.com/google/uuid"
)

type SoftDeleteByIDRequest struct {
	ID uuid.UUID
}

func (s *Service) SoftDeleteByID(ctx context.Context, req SoftDeleteByIDRequest) error {
	// Check if the user exists first
	_, err := s.FindByID(ctx, FindByIDRequest(req))
	if err != nil {
		return err
	}

	err = s.Repo.SoftDeleteByID(ctx, req.ID)
	if err != nil {
		return err
	}
	return nil
}
