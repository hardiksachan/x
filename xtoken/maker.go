// Package xtoken provides a token maker for creating and verifying tokens.
package xtoken

import (
	"time"
)

// Maker is the interface for creating and verifying tokens.
type Maker interface {
	CreateToken(embedding map[string]interface{}, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
