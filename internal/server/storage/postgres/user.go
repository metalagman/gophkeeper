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

// storage.UserRepository interface implementation
var _ storage.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	s := &UserRepository{
		db: db,
	}

	return s, nil
}

// Create implementation of interface storage.UserRepository
func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	const SQL = `
		INSERT INTO users (email, password)
		VALUES ($1, crypt($2, gen_salt('bf')))
		RETURNING id
`

	err := r.db.QueryRowContext(ctx, SQL, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		if pgErr, ok := err.(*pg.Error); ok {
			if pgerrcode.IsIntegrityConstraintViolation(string(pgErr.Code)) {
				return nil, apperr.ErrConflict
			}
		}

		return nil, fmt.Errorf("insert: %w", err)
	}

	return user, nil
}

// Get implementation of interface storage.UserRepository
func (r *UserRepository) Read(ctx context.Context, id uuid.UUID) (*model.User, error) {
	const SQL = `
		SELECT id, email
		FROM users 
		WHERE id=$1
`
	user := &model.User{}

	err := r.db.QueryRowContext(ctx, SQL, id).Scan(&user.ID, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("select: %w", err)
	}

	return user, nil
}

func (r *UserRepository) ReadByEmailAndPassword(ctx context.Context, email string, password string) (*model.User, error) {
	const SQL = `
		SELECT id, email
		FROM users
		WHERE name = $1 
		AND password = crypt($2, password);
`
	user := &model.User{}

	err := r.db.QueryRowContext(ctx, SQL, email, password).Scan(&user.ID, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("select: %w", err)
	}

	return user, nil
}
