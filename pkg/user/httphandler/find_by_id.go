package httphandler

import (
	"net/http"
	"time"

	"okj/lib/otel"
	"okj/lib/responder"
	"okj/pkg/user"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
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
		ctx, span := s.tracer.Start(r.Context(), "user_find_by_id")
		defer span.End()

		id := r.PathValue("userID")
		uuid, err := uuid.Parse(id)
		if err != nil {
			span.SetStatus(codes.Error, "user_find_by_id failed")
			span.RecordError(err)
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "User ID must be a valid UUID.")
			return
		}

		findResponse, err := s.service.FindByID(ctx, user.FindByIDRequest{ID: uuid})
		if err != nil {
			span.SetStatus(codes.Error, "user_find_by_id failed")
			span.RecordError(err)
			switch err {
			case user.ErrNotFoundByID:
				responder.RespondMetaMessage(w, r, http.StatusNotFound, "Could not find any user with provided ID.")
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
			span.SetStatus(codes.Error, "user_find_by_id failed")
			span.RecordError(err)
			s.logger.ErrorContext(ctx, otel.FormatLog(Path, "find_by_id.go [handleUserFindByID]: failed to encode response", err))
			responder.RespondInternalError(w, r)
			return
		}
	}
}
