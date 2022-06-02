package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var _ Manager = (*JWT)(nil)

var ErrInvalidToken = errors.New("invalid token")

type JWT struct {
	secretKey []byte
}

func NewJWT(secretKey string) (*JWT, error) {
	return &JWT{
		secretKey: []byte(secretKey),
	}, nil
}

type JWTClaims struct {
	jwt.StandardClaims
}

func (c *JWTClaims) Identity() string {
	return c.StandardClaims.Id
}

// Issue implementation of token.Manager
func (tm *JWT) Issue(id Identity, lifetime time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(lifetime)
	data := JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        id.Identity(),
			ExpiresAt: exp.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString(tm.secretKey)
}

// Decode implementation of token.Manager
func (tm *JWT) Decode(decode string) (Identity, error) {
	payload := &JWTClaims{}
	_, err := jwt.ParseWithClaims(decode, payload, tm.parseSecretGetter)
	if err != nil {
		return nil, fmt.Errorf("token parse: %w", err)
	}

	if payload.Valid() != nil {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

// Validate implementation of token.Manager
func (tm *JWT) Validate(token string, target Identity) error {
	id, err := tm.Decode(token)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	if target.Identity() != id.Identity() {
		return ErrInvalidToken
	}

	return nil
}

func (tm *JWT) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return tm.secretKey, nil
}
