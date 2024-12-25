package user

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Role string

func NewRole(role string) (Role, error) {
	switch role {
	case "default":
		return Default, nil
	case "admin":
		return Admin, nil
	case "":
		return Default, nil
	default:
		return "", ErrInvalidRole
	}
}

const (
	Default Role = "default"
	Admin   Role = "admin"
)

func (r *Role) Scan(src any) error {
	data, ok := src.(string)
	if !ok {
		return fmt.Errorf("user: role scanner received an invalid string from the db driver")
	}

	*r = Role(data)
	return nil
}

func (r *Role) Value() (driver.Value, error) {
	return string(*r), nil
}

type User struct {
	ID        uuid.UUID
	FirstName string
	LastName  *string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Service struct {
	Repo Repoer
}

type Repoer interface {
	Insert(context.Context, *User) error
	FindByID(context.Context, uuid.UUID) (*User, error)
	UpdateByID(context.Context, uuid.UUID, *User) error
	SoftDeleteByID(context.Context, uuid.UUID) error
}

var (
	ErrInternal          = errors.New("the user service encountered an unexpected condition that prevented it from fulfilling the request")
	ErrNotFoundByID      = errors.New("could not find any user with provided ID")
	ErrInvalidRole       = errors.New("role is not valid")
	ErrUserAlreadyExists = errors.New("an user already exists with the provided ID")
)
