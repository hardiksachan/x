package xhash

import (
	"github.com/Logistics-Coordinators/x/xerrors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	op := xerrors.Op("xhash.HashPassword")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", xerrors.E(op, xerrors.Internal, err)
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	op := xerrors.Op("xhash.ComparePassword")

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return xerrors.E(op, xerrors.Invalid, err)
	}
	return nil
}
