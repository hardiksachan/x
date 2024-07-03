package outbox_test

import (
	"context"
	"testing"
	"time"

	"github.com/hardiksachan/x/xmessage"
	"github.com/hardiksachan/x/xmessage/outbox"
	"github.com/hardiksachan/x/xretry"
	"github.com/hardiksachan/x/xtest"
	"github.com/stretchr/testify/require"
)

func newMessage() *xmessage.Message {
	return &xmessage.Message{
		ID:      xtest.RandomString6(),
		Type:    xtest.RandomString6(),
		Payload: []byte(xtest.RandomString6()),
	}
}

func newPublishing() *xmessage.Publishing {
	return &xmessage.Publishing{
		Message: newMessage(),
		Topic:   xmessage.Topic(xtest.RandomString6()),
	}
}

func newFailablePublishing() *xmessage.Publishing {
	return &xmessage.Publishing{
		Message: newMessage(),
		Topic:   "fail:" + xmessage.Topic(xtest.RandomString6()),
	}
}

func newRetriablePublishing() *xmessage.Publishing {
	return &xmessage.Publishing{
		Message: newMessage(),
		Topic:   "retry:" + xmessage.Topic(xtest.RandomString6()),
	}
}

func TestStart(t *testing.T) {
	ds := newTestDataStore()
	es := newMockPublishingStream()
	r := xretry.NewRetrier(xretry.NewRetryPolicy(xretry.WithNoRetries()))

	o := outbox.New(ds, es, r)

	fp := o.FailedPublishings()

	err := o.Start(context.Background())
	require.NoError(t, err)

	t.Run("when publishings are present in datastore, they are sent to event stream", func(t *testing.T) {
		p := newPublishing()
		ds.AddPublishing(p)

		require.Eventually(t, func() bool {
			return es.isSent(p.Message.ID)
		}, time.Second, time.Millisecond*100)
	})

	t.Run("when publishings are sent to event stream, they are marked as processed", func(t *testing.T) {
		p := newPublishing()
		ds.AddPublishing(p)

		require.Eventually(t, func() bool {
			return ds.isProcessed(p.Message.ID)
		}, time.Second, time.Millisecond*100)
	})

	t.Run("when publishings fail to be sent to event stream, they are marked as processed", func(t *testing.T) {
		p := newFailablePublishing()
		ds.AddPublishing(p)

		fm := <-fp

		require.Equal(t, p, fm.Publishing)

		require.Eventually(t, func() bool {
			return ds.isProcessed(p.Message.ID)
		}, time.Second, time.Millisecond*100)
	})
}

func TestStartWithRetry(t *testing.T) {
	ds := newTestDataStore()
	es := newMockPublishingStream()
	r := xretry.NewRetrier(xretry.NewRetryPolicy(xretry.WithImmediateRetries(10)))

	o := outbox.New(ds, es, r)

	err := o.Start(context.Background())
	require.NoError(t, err)

	t.Run("when publishings fail to be sent to event stream, they are retried", func(t *testing.T) {
		p := newRetriablePublishing()
		ds.AddPublishing(p)

		require.Eventually(t, func() bool {
			return es.isSent(p.Message.ID)
		}, time.Second, time.Millisecond*100)

		require.Eventually(t, func() bool {
			return ds.isProcessed(p.Message.ID)
		}, time.Second, time.Millisecond*100)
	})
}
