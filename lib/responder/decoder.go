package responder

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func Decode[T any](r *http.Request) (v T, err error) {
	if err = json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("responder: decode json: %w", err)
	}
	return v, nil
}

type Field struct {
	Name       string
	Validation string
}

func ValidateInput(validtr *validator.Validate, req any, contract map[string]Field) []ErrorObject {
	errors := make([]ErrorObject, 0, len(contract))

	err := validtr.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			structField := err.StructField()
			t := contract[structField].Name
			d := contract[structField].Validation
			errors = append(errors, ErrorObject{
				Title:  t,
				Detail: &d,
			})
		}
	}
	return errors
}
