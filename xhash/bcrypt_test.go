package xhash_test

import (
	"testing"

	"github.com/hardiksachan/x/xhash"
	"github.com/hardiksachan/x/xtest"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := xtest.RandomString6()
	hashedPassword, err := xhash.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
}

func TestVerifyPassword(t *testing.T) {
	password := xtest.RandomString6()
	hashedPassword, err := xhash.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = xhash.ComparePassword(hashedPassword, password)
	require.NoError(t, err)
}
