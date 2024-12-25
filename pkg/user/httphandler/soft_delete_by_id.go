package httphandler

import (
	"net/http"

	"okj/lib/otel"
	"okj/lib/responder"
	"okj/pkg/user"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

func (s *UserServer) handleUserSoftDeleteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "Bearer token is malformatted.")
			return
		}

		sub, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "Invalid UUID.")
			return
		}

		id := r.PathValue("userID")
		uuid, err := uuid.Parse(id)
		if err != nil {
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "User ID must be a valid UUID.")
			return
		}

		if sub != uuid {
			responder.RespondMetaMessage(w, r, http.StatusForbidden, "You are not allowed to request deletion of other user.")
			return
		}

		err = s.service.SoftDeleteByID(
			r.Context(),
			user.SoftDeleteByIDRequest{ID: uuid},
		)
		if err != nil {
			switch err {
			case user.ErrNotFoundByID:
				responder.RespondMetaMessage(w, r, http.StatusNotFound, "Could not find any user with provided ID.")
			default:
				responder.RespondInternalError(w, r)
			}
			return
		}

		if err := responder.Respond(w, r, http.StatusNoContent, nil); err != nil {
			s.logger.WarnContext(r.Context(), otel.FormatLog(Path, "soft_delete_by_id.go [handleUserSoftDeleteByID]: failed to encode response", err))
			responder.RespondInternalError(w, r)
			return
		}
	}
}
