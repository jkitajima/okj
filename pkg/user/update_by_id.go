package user

import (
	"context"

	"github.com/google/uuid"
)

type UpdateByIDRequest struct {
	ID   uuid.UUID
	User *User
}

type UpdateByIDResponse struct {
	User *User
}

func (s *Service) UpdateByID(ctx context.Context, req UpdateByIDRequest) (UpdateByIDResponse, error) {
	// Check if user exists first
	findResult, err := s.FindByID(ctx, FindByIDRequest{req.ID})
	if err != nil {
		return UpdateByIDResponse{nil}, err
	}

	findResult.User.FirstName = req.User.FirstName
	findResult.User.LastName = req.User.LastName
	findResult.User.Role = req.User.Role

	err = s.Repo.UpdateByID(ctx, req.ID, findResult.User)
	if err != nil {
		return UpdateByIDResponse{nil}, err
	}
	return UpdateByIDResponse(findResult), nil
}
