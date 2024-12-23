package httphandler

import (
	"net/http"
	"time"

	"okj/internal/user"
	"okj/pkg/responder"

	"github.com/google/uuid"
)

func (s *UserServer) handleUserUpdateByID() http.HandlerFunc {
	type request struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role" validate:"omitempty,oneof=default admin"`
	}

	type response struct {
		Entity    string    `json:"entity"`
		ID        uuid.UUID `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  *string   `json:"last_name"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := responder.Decode[request](r)
		if err != nil {
			responder.RespondInternalError(w, r)
			return
		}

		id := r.PathValue("userID")
		uuid, err := uuid.Parse(id)
		if err != nil {
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "User ID must be a valid UUID.")
			return
		}

		updateResponse, err := s.service.UpdateByID(r.Context(), user.UpdateByIDRequest{
			ID: uuid,
			User: &user.User{
				FirstName: req.FirstName,
				LastName:  &req.LastName,
				Role:      user.Role(req.Role),
			},
		})
		if err != nil {
			switch err {
			case user.ErrNotFoundByID:
				responder.RespondMetaMessage(w, r, http.StatusBadRequest, "Could not find any user with provided ID.")
			default:
				responder.RespondInternalError(w, r)
			}
			return
		}

		resp := response{
			Entity:    s.entity,
			ID:        updateResponse.User.ID,
			FirstName: updateResponse.User.FirstName,
			LastName:  updateResponse.User.LastName,
			Role:      string(updateResponse.User.Role),
			CreatedAt: updateResponse.User.CreatedAt,
			UpdatedAt: updateResponse.User.UpdatedAt,
		}

		if err := responder.Respond(w, r, http.StatusOK, &responder.DataField{Data: resp}); err != nil {
			responder.RespondInternalError(w, r)
			return
		}
	}
}
