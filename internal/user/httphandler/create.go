package httphandler

import (
	"net/http"
	"time"

	"okj/internal/user"
	"okj/pkg/responder"

	"github.com/google/uuid"
)

func (s *UserServer) handleUserCreate() http.HandlerFunc {
	type request struct {
		FirstName string `json:"first_name" validate:"required"`
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

	contract := map[string]responder.Field{
		"FirstName": {
			Name:       "first_name",
			Validation: "Field is required and cannot be an empty string.",
		},
		"LastName": {
			Name:       "last_name",
			Validation: "Field value cannot be an empty string.",
		},
		"Role": {
			Name:       "role",
			Validation: "Field value must be either 'default' or 'admin'.",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := responder.Decode[request](r)
		if err != nil {
			responder.RespondMetaMessage(w, r, http.StatusBadRequest, "Request body is invalid.")
			return
		}

		if errors := responder.ValidateInput(s.inputValidator, req, contract); len(errors) > 0 {
			responder.RespondClientErrors(w, r, errors...)
			return
		}

		createResponse, err := s.service.Create(r.Context(), user.CreateRequest{
			User: &user.User{
				FirstName: req.FirstName,
				LastName:  &req.LastName,
				Role:      user.Role(req.Role),
			},
		})
		if err != nil {
			responder.RespondInternalError(w, r)
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
			responder.RespondInternalError(w, r)
			return
		}
	}
}
