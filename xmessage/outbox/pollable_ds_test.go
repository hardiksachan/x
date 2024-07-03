package outbox_test

import (
	"context"
	"testing"
	"time"

	"github.com/hardiksachan/x/xmessage/outbox"
	"github.com/stretchr/testify/require"
)

func TestGetUnsentMessage(t *testing.T) {
	r := newMockMessageRepository()

	maxLockAge := time.Second * 2

	poller := outbox.NewPollableDataSource(r, outbox.NewPollingPolicy(outbox.WithLock(time.Millisecond*100, maxLockAge)))

	msgChan, err := poller.GetUnsentPublishings(context.TODO())
	require.NoError(t, err)

	t.Run("GetUnsentMessage must return message that are added to repository", func(t *testing.T) {
		m := newPublishing()
		err = r.AddMessage(context.TODO(), *m)
		require.NoError(t, err)

		msg := <-msgChan

		require.Equal(t, m.Message.ID, msg.Message.ID)
	})

	t.Run("GetUnsentMessage must not return message that are locked", func(t *testing.T) {
		m := newPublishing()
		err = r.AddMessage(context.TODO(), *m)
		require.NoError(t, err)

		msg := <-msgChan

		require.Equal(t, m.Message.ID, msg.Message.ID)

		// Wait for a second to allow poller to lock the message
		time.Sleep(time.Second)

		poller2 := outbox.NewPollableDataSource(r, outbox.NewPollingPolicy())
		msgChan2, err := poller2.GetUnsentPublishings(context.TODO())
		require.NoError(t, err)

		select {
		case msg2 := <-msgChan2:
			t.Errorf("unexpected message: %v", msg2)
		default:
			// No message received, which is expected
		}
	})
}

func TestSetAsProcessed(t *testing.T) {
	r := newMockMessageRepository()

	maxLockAge := time.Second * 2

	poller := outbox.NewPollableDataSource(r, outbox.NewPollingPolicy(outbox.WithLock(time.Millisecond*100, maxLockAge)))

	t.Run("SetAsProcessed must set message as processed", func(t *testing.T) {
		m := newPublishing()
		err := r.AddMessage(context.TODO(), *m)
		require.NoError(t, err)

		err = poller.SetAsProcessed(context.Background(), m.Message.ID)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			return r.isProcessed(m.Message.ID)
		}, time.Second, time.Millisecond*100)
	})
}

func TestClearLocks(t *testing.T) {
	r := newMockMessageRepository()

	maxLockAge := time.Second * 2

	poller := outbox.NewPollableDataSource(r, outbox.NewPollingPolicy(outbox.WithLock(time.Millisecond*100, maxLockAge)))

	msgChan, err := poller.GetUnsentPublishings(context.TODO())
	require.NoError(t, err)

	t.Run("Locks must be cleared after timeout", func(t *testing.T) {
		m := newPublishing()
		err := r.AddMessage(context.TODO(), *m)
		require.NoError(t, err)

		<-msgChan

		// Wait for a second to allow poller to lock the message
		time.Sleep(time.Second)

		require.Eventually(t, func() bool {
			return !r.isLocked(m.Message.ID)
		}, maxLockAge+time.Second, time.Millisecond*100)
	})
}
