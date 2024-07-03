package outbox_test

import (
	"context"
	"time"

	"github.com/hardiksachan/x/xerrors"
	"github.com/hardiksachan/x/xmessage"
)

type memMessage struct {
	message     xmessage.Publishing
	attempts    int
	nextRetryAt time.Time
	lockedBy    string
	lockedAt    time.Time
}

type mockMessageRepository struct {
	messages map[string]*memMessage
}

func newMockMessageRepository() *mockMessageRepository {
	return &mockMessageRepository{
		messages: make(map[string]*memMessage),
	}
}

func (m *mockMessageRepository) AddMessage(_ context.Context, msg xmessage.Publishing) error {
	op := xerrors.Op("outbox.mockMessageRepository.AddMessage")

	if _, ok := m.messages[msg.Message.ID]; ok {
		return xerrors.E(op, xerrors.Message("message already exists"))
	}

	m.messages[msg.Message.ID] = &memMessage{
		message:     msg,
		attempts:    0,
		nextRetryAt: time.Now(),
		lockedBy:    "",
		lockedAt:    time.Time{},
	}

	return nil
}

func (m *mockMessageRepository) GetUnsentPublishing(_ context.Context, instanceID string, maxRetries int) (*xmessage.Publishing, error) {
	op := xerrors.Op("outbox.mockMessageRepository.GetUnsentMessage")

	for _, msg := range m.messages {
		if msg.lockedBy == "" && msg.attempts < maxRetries && msg.nextRetryAt.Before(time.Now()) {
			msg.lockedBy = instanceID
			msg.lockedAt = time.Now()
			return &msg.message, nil
		}
	}

	return nil, xerrors.E(op, xerrors.Message("no unsent messages"))
}

func (m *mockMessageRepository) SetAsProcessed(_ context.Context, id string) error {
	op := xerrors.Op("outbox.mockMessageRepository.SetAsProcessed")

	msg, ok := m.messages[id]
	if !ok {
		return xerrors.E(op, xerrors.Message("message not found"))
	}

	msg.lockedBy = ""
	msg.lockedAt = time.Time{}
	msg.attempts++

	return nil
}

func (m *mockMessageRepository) MarkForRetry(_ context.Context, id string, retryAt time.Time) error {
	op := xerrors.Op("outbox.mockMessageRepository.MarkForRetry")

	msg, ok := m.messages[id]
	if !ok {
		return xerrors.E(op, xerrors.Message("message not found"))
	}

	msg.nextRetryAt = retryAt

	return nil
}

func (m *mockMessageRepository) ClearLocks(_ context.Context, instanceID string, maxLockAge time.Time) error {
	for _, msg := range m.messages {
		if msg.lockedBy == instanceID && msg.lockedAt.Before(maxLockAge) {
			msg.lockedBy = ""
			msg.lockedAt = time.Time{}
		}
	}

	return nil
}

func (m *mockMessageRepository) isProcessed(id string) bool {
	msg, ok := m.messages[id]
	if !ok {
		return false
	}

	return msg.attempts > 0
}

func (m *mockMessageRepository) isLocked(id string) bool {
	msg, ok := m.messages[id]
	if !ok {
		return false
	}

	return msg.lockedBy != "" && msg.lockedAt != time.Time{}
}
