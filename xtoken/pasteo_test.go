package xtoken_test

import (
	"testing"
	"time"

	"github.com/hardiksachan/x/xtest"
	"github.com/hardiksachan/x/xtoken"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := xtoken.NewPasetoMaker(xtest.RandomString(32))
	require.NoError(t, err)

	email := xtest.RandomEmailString()
	id := xtest.RandomString(32)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(id, email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)
	require.Equal(t, id, payload.UserID)
	require.Equal(t, email, payload.Email)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.TokenID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := xtoken.NewPasetoMaker(xtest.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(xtest.RandomString(32), xtest.RandomEmailString(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}
