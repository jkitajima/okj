package httphandler

import (
	"net/http"

	"okj/lib/otel"
	"okj/lib/responder"
	"okj/pkg/user"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

func (s *UserServer) handleUserSoftDeleteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := s.tracer.Start(r.Context(), "user_soft_delete_by_id")
		defer span.End()

		_, claims, err := jwtauth.FromContext(ctx)
		if err != nil {
			span.SetStatus(codes.Error, "user_soft_delete_by_id failed")
			span.RecordError(err)
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "Bearer token is malformatted.")
			return
		}

		sub, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			span.SetStatus(codes.Error, "user_soft_delete_by_id failed")
			span.RecordError(err)
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "Invalid UUID.")
			return
		}

		id := r.PathValue("userID")
		uuid, err := uuid.Parse(id)
		if err != nil {
			span.SetStatus(codes.Error, "user_soft_delete_by_id failed")
			span.RecordError(err)
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "User ID must be a valid UUID.")
			return
		}

		if sub != uuid {
			span.SetStatus(codes.Error, "user_soft_delete_by_id failed")
			span.RecordError(err)
			responder.RespondMetaMessage(w, r, http.StatusForbidden, "You are not allowed to request deletion of other user.")
			return
		}

		err = s.service.SoftDeleteByID(ctx, user.SoftDeleteByIDRequest{ID: uuid})
		if err != nil {
			span.SetStatus(codes.Error, "user_soft_delete_by_id failed")
			span.RecordError(err)
			switch err {
			case user.ErrNotFoundByID:
				responder.RespondMetaMessage(w, r, http.StatusNotFound, "Could not find any user with provided ID.")
			default:
				responder.RespondInternalError(w, r)
			}
			return
		}

		if err := responder.Respond(w, r, http.StatusNoContent, nil); err != nil {
			span.SetStatus(codes.Error, "user_soft_delete_by_id failed")
			span.RecordError(err)
			s.logger.ErrorContext(ctx, otel.FormatLog(Path, "soft_delete_by_id.go [handleUserSoftDeleteByID]: failed to encode response", err))
			responder.RespondInternalError(w, r)
			return
		}
	}
}
