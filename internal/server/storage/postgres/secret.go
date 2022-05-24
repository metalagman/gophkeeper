package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
func (r *SecretRepository) Create(ctx context.Context, uid uuid.UUID, secret *model.Secret) (*model.Secret, error) {
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

func (r *SecretRepository) ReadByName(ctx context.Context, uid uuid.UUID, name string) (*model.Secret, error) {
	const SQL = `
		SELECT id, type, name, content
		FROM secrets
		WHERE user_id = $1 AND name = $2;
`
	m := &model.Secret{}

	err := r.db.QueryRowContext(ctx, SQL, uid.String(), name).Scan(&m.ID, &m.Type, &m.Name, &m.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("select: %w", err)
	}

	return m, nil
}

func (r *SecretRepository) DeleteByName(ctx context.Context, uid uuid.UUID, name string) error {
	const SQL = `
		DELETE
		FROM secrets
		WHERE user_id = $1 AND name = $2;
`
	_, err := r.db.ExecContext(ctx, SQL, uid.String(), name)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (r *SecretRepository) List(ctx context.Context, uid uuid.UUID) ([]*model.Secret, error) {
	const SQL = `
		SELECT
			id,
			type,
			name
		FROM secrets
		WHERE user_id = $1
		ORDER BY name
`
	rows, err := r.db.QueryContext(ctx, SQL, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("select: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	res := make([]*model.Secret, 0)

	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows next: %w", err)
		}
		m := &model.Secret{}
		if err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.Type,
			&m.Name,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		res = append(res, m)
	}

	return res, nil
}
