//go:generate mockgen -source=./interface.go -destination=./mock/storage.go -package=storagemock
package storage

import (
	"context"
	"github.com/google/uuid"
	"gophkeeper/internal/server/model"
)

type UserRepository interface {
	// Create a new model.User
	Create(ctx context.Context, m *model.User) (*model.User, error)
	// ReadByEmailAndPassword instance of model.User
	ReadByEmailAndPassword(ctx context.Context, name string, password string) (*model.User, error)
	// Read instance of model.User
	Read(ctx context.Context, id uuid.UUID) (*model.User, error)
}
