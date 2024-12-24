package gorm

import (
	"time"

	"okj/pkg/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	FirstName string    `gorm:"not null"`
	LastName  *string
	Role      user.Role      `gorm:"not null;type:user_role;default:'default'"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (*UserModel) TableName() string {
	return "User"
}
