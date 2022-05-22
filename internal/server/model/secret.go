package model

import (
	"github.com/google/uuid"
)

type Secret struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Name    string
	Type    string
	Content []byte
}
