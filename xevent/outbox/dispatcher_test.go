package outbox_test

import (
	"context"
	"testing"
	"time"

	"github.com/Logistics-Coordinators/x/xevent/outbox"
	"github.com/Logistics-Coordinators/x/xretry"
	"github.com/Logistics-Coordinators/x/xtest"
	"github.com/stretchr/testify/require"
)

func newMessage() *outbox.Message {
	return &outbox.Message{
		ID:      xtest.RandomString6(),
		Topic:   xtest.RandomString6(),
		Type:    xtest.RandomString6(),
		Payload: []byte(xtest.RandomString6()),
	}
}

func newFailableMessage() *outbox.Message {
	return &outbox.Message{
		ID:      "fail:" + xtest.RandomString6(),
		Topic:   xtest.RandomString6(),
		Type:    xtest.RandomString6(),
		Payload: []byte(xtest.RandomString6()),
	}
}

func newRetriableMessage() *outbox.Message {
	return &outbox.Message{
		ID:      "retry:" + xtest.RandomString6(),
		Topic:   xtest.RandomString6(),
		Type:    xtest.RandomString6(),
		Payload: []byte(xtest.RandomString6()),
	}
}

func TestStart(t *testing.T) {
	ds := newTestDataStore()
	es := newTestEventStream()
	r := xretry.NewRetrier(xretry.NewRetryPolicy(xretry.WithNoRetries()))

	o := outbox.New(ds, es, r)

	failedMessages := o.FailedMessages()

	err := o.Start(context.Background())
	require.NoError(t, err)

	t.Run("when messages are present in datastore, they are sent to event stream", func(t *testing.T) {
		m := newMessage()
		ds.AddMessage(m)

		require.Eventually(t, func() bool {
			return es.isSent(m.ID)
		}, time.Second, time.Millisecond*100)
	})

	t.Run("when messages are sent to event stream, they are marked as processed", func(t *testing.T) {
		m := newMessage()
		ds.AddMessage(m)

		require.Eventually(t, func() bool {
			return ds.isProcessed(m.ID)
		}, time.Second, time.Millisecond*100)
	})

	t.Run("when messages fail to be sent to event stream, they are marked as processed", func(t *testing.T) {
		m := newFailableMessage()
		ds.AddMessage(m)

		fm := <-failedMessages

		require.Equal(t, m, fm.Msg)

		require.Eventually(t, func() bool {
			return ds.isProcessed(m.ID)
		}, time.Second, time.Millisecond*100)
	})
}

func TestStartWithRetry(t *testing.T) {
	ds := newTestDataStore()
	es := newTestEventStream()
	r := xretry.NewRetrier(xretry.NewRetryPolicy(xretry.WithImmediateRetries(10)))

	o := outbox.New(ds, es, r)

	err := o.Start(context.Background())
	require.NoError(t, err)

	t.Run("when messages fail to be sent to event stream, they are retried", func(t *testing.T) {
		m := newRetriableMessage()
		ds.AddMessage(m)

		require.Eventually(t, func() bool {
			return es.isSent(m.ID)
		}, time.Second, time.Millisecond*100)

		require.Eventually(t, func() bool {
			return ds.isProcessed(m.ID)
		}, time.Second, time.Millisecond*100)
	})
}
