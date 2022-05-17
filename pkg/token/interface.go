//go:generate mockgen -source=./interface.go -destination=./mock/token.go -package=tokenmock
package token

import "time"

type Identity interface {
	Identity() string
}

type Manager interface {
	// Issue a new token for a given Identity with exp time
	Issue(id Identity, exp time.Duration) (string, error)
	// Decode provided token to the Identity
	Decode(tk string) (Identity, error)
	// Validate if a provided token valid for a target Identity
	Validate(token string, target Identity) error
}
