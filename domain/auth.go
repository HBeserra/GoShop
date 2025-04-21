package domain

import "github.com/google/uuid"

type Session struct {
	UserID uuid.UUID
	Role   SystemRole
}

type SystemRole string

const (
	SystemRoleAdmin SystemRole = "admin"
	SystemRoleUser  SystemRole = "user"
)
