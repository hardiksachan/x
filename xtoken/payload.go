package xtoken

import (
	"time"

	"github.com/Logistics-Coordinators/x/xerrors"
	"github.com/google/uuid"
)

// Payload is the payload of a token.
type Payload struct {
	TokenID   string                 `json:"token_id"`
	IssuedAt  time.Time              `json:"issued_at"`
	ExpiredAt time.Time              `json:"expired_at"`
	Embedding map[string]interface{} `json:"embedding"`
}

// NewPayload returns a new Payload.
func NewPayload(embedding map[string]interface{}, duration time.Duration) (*Payload, error) {
	op := xerrors.Op("xtoken.NewPayload")

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, xerrors.E(op, xerrors.Internal, err)
	}

	payload := &Payload{
		TokenID:   tokenID.String(),
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
		Embedding: embedding,
	}
	return payload, nil
}

// Valid returns an error if a payload is invalid.
func (payload *Payload) Valid() error {
	op := xerrors.Op("xtoken.Payload.Valid")
	if time.Now().After(payload.ExpiredAt) {
		return xerrors.E(op, xerrors.Expired, xerrors.Message("Your session has expired, please login again."))
	}
	return nil
}
