package xtoken

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/hardiksachan/x/xerrors"
	"github.com/o1egl/paseto"
)

const (
	// ErrExpiredToken is returned when a token has expired.
	ErrExpiredToken = xerrors.Message("Your session has expired")
)

// PasetoMaker is a Paseto token maker.
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker returns a new PasetoMaker.
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	op := xerrors.Op("xtoken.NewPasetoMaker")

	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, xerrors.E(
			op,
			xerrors.Invalid,
			fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize),
		)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// CreateToken creates a new token.
func (maker *PasetoMaker) CreateToken(userID, email string, duration time.Duration) (string, *Payload, error) {
	op := xerrors.Op("xtoken.PasetoMaker.CreateToken")

	payload, err := NewPayload(userID, email, duration)
	if err != nil {
		return "", payload, xerrors.E(op, xerrors.Internal, err)
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", payload, xerrors.E(op, xerrors.Internal, err)
	}
	return token, payload, nil
}

// VerifyToken verifies a token.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	op := xerrors.Op("xtoken.PasetoMaker.VerifyToken")

	//nolint:exhaustruct
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, xerrors.E(op, xerrors.Invalid, err, ErrExpiredToken)
	}

	err = payload.Valid()
	if err != nil {
		return nil, xerrors.E(op, xerrors.Invalid, err)
	}

	return payload, nil
}
