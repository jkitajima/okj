package httphandler

import (
	"net/http"
	"time"

	"okj/internal/user"
	"okj/pkg/responder"

	"github.com/google/uuid"
)

func (s *UserServer) handleUserFindByID() http.HandlerFunc {
	type response struct {
		Entity    string    `json:"entity"`
		ID        uuid.UUID `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("userID")
		uuid, err := uuid.Parse(id)
		if err != nil {
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "User ID must be a valid UUID.")
			return
		}

		findResponse, err := s.service.FindByID(r.Context(), user.FindByIDRequest{ID: uuid})
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
			ID:        findResponse.User.ID,
			FirstName: findResponse.User.FirstName,
			LastName:  *findResponse.User.LastName,
			Role:      string(findResponse.User.Role),
			CreatedAt: findResponse.User.CreatedAt,
			UpdatedAt: findResponse.User.UpdatedAt,
		}

		if err := responder.Respond(w, r, http.StatusOK, &responder.DataField{Data: resp}); err != nil {
			responder.RespondInternalError(w, r)
			return
		}
	}
}
