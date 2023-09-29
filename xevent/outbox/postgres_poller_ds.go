package outbox

import (
	"context"
	"time"

	"github.com/Logistics-Coordinators/x/xerrors"
	"github.com/Logistics-Coordinators/x/xlog"
	"github.com/google/uuid"
)

const (
	defaultPollingInterval = 1 * time.Second
	defaultLockInterval    = 5 * time.Second
	defaultMaxLockAge      = 120 * time.Second
	defaultMaxRetries      = 3
	defaultRetryInterval   = 30 * time.Second
)

// MessageRepository is the interface that will need to be implemented by the consumer
type MessageRepository interface {
	GetUnsentMessage(ctx context.Context, instanceID string, maxRetries int) (*Message, error)
	SetAsProcessed(ctx context.Context, id string) error
	MarkForRetry(ctx context.Context, id string, retryAt time.Time) error
	ClearLocks(ctx context.Context, instanceID string, obtainedBefore time.Time) error
}

// PollingPolicy is used to configure the polling interval
type PollingPolicy struct {
	pollingInterval time.Duration
	lockingInterval time.Duration
	maxLockAge      time.Duration
	maxRetries      int
	retryInterval   time.Duration
}

// PollingOption is used to configure the polling interval
type PollingOption func(*PollingPolicy)

// WithPollingInterval sets the polling interval
func WithPollingInterval(interval time.Duration) PollingOption {
	return func(p *PollingPolicy) {
		p.pollingInterval = interval
	}
}

// WithRetries sets the max number of retries and retry interval
func WithRetries(maxRetries int, retryInterval time.Duration) PollingOption {
	return func(p *PollingPolicy) {
		p.maxRetries = maxRetries
		p.retryInterval = retryInterval
	}
}

// WithLock sets the locking interval
func WithLock(interval time.Duration, age time.Duration) PollingOption {
	return func(p *PollingPolicy) {
		p.lockingInterval = interval
		p.maxLockAge = age
	}
}

// NewPollingPolicy creates a new PollingPolicy
func NewPollingPolicy(opts ...PollingOption) *PollingPolicy {
	p := PollingPolicy{
		pollingInterval: defaultPollingInterval,
		lockingInterval: defaultLockInterval,
		maxLockAge:      defaultMaxLockAge,
		maxRetries:      defaultMaxRetries,
		retryInterval:   defaultRetryInterval,
	}
	for _, opt := range opts {
		opt(&p)
	}
	return &p
}

// PostgresPoller implements outbox.DataStore that polls the database for unsent messages
type PostgresPoller struct {
	r          MessageRepository
	p          *PollingPolicy
	instanceID string
}

// NewPostgresPoller creates a new PostgresPoller
func NewPostgresPoller(r MessageRepository, p *PollingPolicy) *PostgresPoller {
	return &PostgresPoller{
		instanceID: uuid.NewString(),
		r:          r,
		p:          p,
	}
}

func (p *PostgresPoller) startPolling(ctx context.Context, messages chan<- *Message) {
	op := xerrors.Op("outbox.PostgresPoller.startPolling")

	for {
		time.Sleep(p.p.pollingInterval)

		msg, err := p.r.GetUnsentMessage(ctx, p.instanceID, p.p.maxRetries)
		if err == nil {
			messages <- msg
			continue
		}

		xlog.Infof("%s: error getting unsent messages: %+v", op, err)
	}
}

func (p *PostgresPoller) clearLocks(ctx context.Context) {
	op := xerrors.Op("outbox.PostgresPoller.clearLocks")

	for {
		time.Sleep(p.p.lockingInterval)

		err := p.r.ClearLocks(ctx, p.instanceID, time.Now().Add(-p.p.maxLockAge))
		if err != nil {
			xlog.Infof("%s, error clearing locks: %+v", op, err)
		}
	}
}

// GetUnsentMessages will return all unsent messages
func (p *PostgresPoller) GetUnsentMessages(ctx context.Context) (<-chan *Message, error) {
	messages := make(chan *Message)

	go p.startPolling(ctx, messages)
	go p.clearLocks(ctx)

	return messages, nil
}

// SetAsProcessed will set the message as processed
func (p *PostgresPoller) SetAsProcessed(ctx context.Context, id string) error {
	op := xerrors.Op("outbox.PostgresPoller.SetAsProcessed")

	err := p.r.SetAsProcessed(ctx, id)
	if err != nil {
		return xerrors.E(op, err)
	}

	return nil
}

// RetryMessage will set failed message to be retried
func (p *PostgresPoller) RetryMessage(ctx context.Context, id string) error {
	op := xerrors.Op("outbox.PostgresPoller.RetryMessage")

	err := p.r.MarkForRetry(ctx, id, time.Now().Add(p.p.retryInterval))
	if err != nil {
		return xerrors.E(op, err)
	}

	return nil
}
