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

type SecretRepository interface {
	// Create a new model.Secret
	Create(ctx context.Context, uid uuid.UUID, m *model.Secret) (*model.Secret, error)
	// ReadByName specified secret
	ReadByName(ctx context.Context, uid uuid.UUID, name string) (*model.Secret, error)
	// DeleteByName specified secret if available
	DeleteByName(ctx context.Context, uid uuid.UUID, name string) error
	// List all secrets of specified user
	List(ctx context.Context, uid uuid.UUID) ([]*model.Secret, error)
}
