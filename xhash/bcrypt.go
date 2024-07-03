// Package xhash provides hashing functions.
package xhash

import (
	"github.com/hardiksachan/x/xerrors"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns a hashed password.
func HashPassword(password string) (string, error) {
	op := xerrors.Op("xhash.HashPassword")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", xerrors.E(op, xerrors.Internal, err)
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a hashed password with a password.
func ComparePassword(hashedPassword, password string) error {
	op := xerrors.Op("xhash.ComparePassword")

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return xerrors.E(op, xerrors.Invalid, err)
	}
	return nil
}
