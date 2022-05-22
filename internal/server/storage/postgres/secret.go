package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgerrcode"
	pg "github.com/lib/pq"
	"gophkeeper/internal/server/model"
	"gophkeeper/internal/server/storage"
	"gophkeeper/pkg/apperr"
)

// storage.SecretRepository interface implementation
var _ storage.SecretRepository = (*SecretRepository)(nil)

type SecretRepository struct {
	db *sql.DB
}

func NewSecretRepository(db *sql.DB) (*SecretRepository, error) {
	s := &SecretRepository{
		db: db,
	}

	return s, nil
}

// Create implementation of interface storage.SecretRepository
func (r *SecretRepository) Create(ctx context.Context, secret *model.Secret) (*model.Secret, error) {
	const SQL = `
		INSERT INTO secrets (user_id, type, name, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id
`

	err := r.db.QueryRowContext(ctx, SQL, secret.UserID, secret.Type, secret.Name, secret.Content).Scan(&secret.ID)
	if err != nil {
		if pgErr, ok := err.(*pg.Error); ok {
			if pgerrcode.IsIntegrityConstraintViolation(string(pgErr.Code)) {
				return nil, apperr.ErrConflict
			}
		}

		return nil, fmt.Errorf("insert: %w", err)
	}

	return secret, nil
}
