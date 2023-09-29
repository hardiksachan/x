package outbox_test

import (
	"context"
	"testing"
	"time"

	"github.com/Logistics-Coordinators/x/xevent/outbox"
	"github.com/stretchr/testify/require"
)

func TestPostgresPoller(t *testing.T) {
	r := newMockMessageRepository()

	maxLockAge := time.Second * 2

	poller := outbox.NewPostgresPoller(r, outbox.NewPollingPolicy(outbox.WithLock(time.Millisecond*100, maxLockAge)))

	msgChan, err := poller.GetUnsentMessages(context.TODO())
	require.NoError(t, err)

	t.Run("GetUnsentMessage must return message that are added to repository", func(t *testing.T) {
		m := newMessage()
		err = r.AddMessage(context.TODO(), *m)
		require.NoError(t, err)

		msg := <-msgChan

		require.Equal(t, m.ID, msg.ID)
	})

	t.Run("GetUnsentMessage must not return message that are locked", func(t *testing.T) {
		m := newMessage()
		err = r.AddMessage(context.TODO(), *m)
		require.NoError(t, err)

		msg := <-msgChan

		require.Equal(t, m.ID, msg.ID)

		// Wait for a second to allow poller to lock the message
		time.Sleep(time.Second)

		poller2 := outbox.NewPostgresPoller(r, outbox.NewPollingPolicy())
		msgChan2, err := poller2.GetUnsentMessages(context.TODO())
		require.NoError(t, err)

		select {
		case msg2 := <-msgChan2:
			t.Errorf("unexpected message: %v", msg2)
		default:
			// No message received, which is expected
		}
	})

	t.Run("SetAsProcessed must set message as processed", func(t *testing.T) {
		m := newMessage()
		err := r.AddMessage(context.TODO(), *m)
		require.NoError(t, err)

		err = poller.SetAsProcessed(context.Background(), m.ID)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			return r.isProcessed(m.ID)
		}, time.Second, time.Millisecond*100)
	})

	t.Run("Locks must be cleared after timeout", func(t *testing.T) {
		m := newMessage()
		err := r.AddMessage(context.TODO(), *m)
		require.NoError(t, err)

		// Wait for a second to allow poller to lock the message
		time.Sleep(time.Second)

		require.Eventually(t, func() bool {
			return !r.isLocked(m.ID)
		}, maxLockAge+time.Second, time.Millisecond*100)
	})
}
