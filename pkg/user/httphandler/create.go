package httphandler

import (
	"net/http"
	"time"

	"okj/lib/otel"
	"okj/lib/responder"
	"okj/pkg/user"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

func (s *UserServer) handleUserCreate() http.HandlerFunc {
	type request struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name"`
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

	contract := map[string]responder.Field{
		"FirstName": {
			Name:       "first_name",
			Validation: "Field is required and cannot be an empty string.",
		},
		"LastName": {
			Name:       "last_name",
			Validation: "Field value cannot be an empty string.",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := s.tracer.Start(r.Context(), "user_create")
		defer span.End()

		req, err := responder.Decode[request](r)
		if err != nil {
			span.SetStatus(codes.Error, "user_create failed")
			span.RecordError(err)
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "Request body is invalid.")
			return
		}

		if errors := responder.ValidateInput(s.inputValidator, req, contract); len(errors) > 0 {
			span.SetStatus(codes.Error, "user_create failed")
			span.RecordError(err)
			responder.RespondClientErrors(w, r, errors...)
			return
		}

		_, claims, err := jwtauth.FromContext(ctx)
		if err != nil {
			span.SetStatus(codes.Error, "user_create failed")
			span.RecordError(err)
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "Bearer token is malformatted.")
			return
		}

		uuid, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			span.SetStatus(codes.Error, "user_create failed")
			span.RecordError(err)
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "Invalid UUID.")
			return
		}

		createResponse, err := s.service.Create(ctx, user.CreateRequest{
			User: &user.User{
				ID:        uuid,
				FirstName: req.FirstName,
				LastName:  &req.LastName,
				Role:      user.Default,
			},
		})
		if err != nil {
			span.SetStatus(codes.Error, "user_create failed")
			span.RecordError(err)
			switch err {
			case user.ErrUserAlreadyExists:
				responder.RespondMetaMessage(w, r, http.StatusBadRequest, "There is already an user registered for this token subject.")
			case user.ErrInternal:
				fallthrough
			default:
				responder.RespondInternalError(w, r)
			}
			return
		}

		resp := response{
			Entity:    s.entity,
			ID:        createResponse.User.ID,
			FirstName: createResponse.User.FirstName,
			LastName:  createResponse.User.LastName,
			Role:      string(createResponse.User.Role),
			CreatedAt: createResponse.User.CreatedAt,
			UpdatedAt: createResponse.User.UpdatedAt,
		}

		if err := responder.Respond(w, r, http.StatusCreated, &responder.DataField{Data: resp}); err != nil {
			span.SetStatus(codes.Error, "user_create failed")
			span.RecordError(err)
			s.logger.ErrorContext(ctx, otel.FormatLog(Path, "create.go [handleUserCreate]: failed to encode response", err))
			responder.RespondInternalError(w, r)
			return
		}
	}
}
