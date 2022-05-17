package apperr

import (
	"errors"
	"fmt"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal")
	ErrNotFound     = fmt.Errorf("not found: %w", ErrInvalidInput)
	ErrConflict     = fmt.Errorf("conflict: %w", ErrInvalidInput)
	ErrSoftConflict = fmt.Errorf("soft conflict: %w", ErrInvalidInput)
)
