package httphandler

import (
	"net/http"

	"okj/internal/user"
	"okj/pkg/responder"

	"github.com/google/uuid"
)

func (s *UserServer) handleUserSoftDeleteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("userID")
		uuid, err := uuid.Parse(id)
		if err != nil {
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "User ID must be a valid UUID.")
			return
		}

		err = s.service.SoftDeleteByID(
			r.Context(),
			user.SoftDeleteByIDRequest{ID: uuid},
		)
		if err != nil {
			switch err {
			case user.ErrNotFoundByID:
				responder.RespondMetaMessage(w, r, http.StatusBadRequest, "Could not find any user with provided ID.")
			default:
				responder.RespondInternalError(w, r)
			}
			return
		}

		if err := responder.Respond(w, r, http.StatusNoContent, nil); err != nil {
			responder.RespondInternalError(w, r)
			return
		}
	}
}
