package outbox

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hardiksachan/x/xerrors"
	"github.com/hardiksachan/x/xlog"
	"github.com/hardiksachan/x/xmessage"
)

const (
	defaultPollingInterval = 1 * time.Second
	defaultLockInterval    = 5 * time.Second
	defaultMaxLockAge      = 120 * time.Second
	defaultMaxRetries      = 3
	defaultRetryInterval   = 30 * time.Second
)

// PollableRepository is the interface that will need to be implemented by the consumer
type PollableRepository interface {
	GetUnsentPublishing(ctx context.Context, instanceID string, maxRetries int) (*xmessage.Publishing, error)
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

// PollableDataSource implements outbox.DataStore that polls the database for unsent messages
type PollableDataSource struct {
	r          PollableRepository
	p          *PollingPolicy
	instanceID string
}

// NewPollableDataSource creates a new PollableDataSource
func NewPollableDataSource(r PollableRepository, p *PollingPolicy) *PollableDataSource {
	return &PollableDataSource{
		instanceID: uuid.NewString(),
		r:          r,
		p:          p,
	}
}

func (p *PollableDataSource) startPolling(ctx context.Context, publishings chan<- *xmessage.Publishing) {
	op := xerrors.Op("outbox.PollableDataSource.startPolling")

	for {
		time.Sleep(p.p.pollingInterval)

		publishing, err := p.r.GetUnsentPublishing(ctx, p.instanceID, p.p.maxRetries)
		if err == nil {
			publishings <- publishing
			continue
		}

		if xerrors.ErrorCode(err) != xerrors.NotFound {
			xlog.Infof("%s: error getting unsent publishings: %+v", op, err)
		}
	}
}

func (p *PollableDataSource) clearLocks(ctx context.Context) {
	op := xerrors.Op("outbox.PollableDataSource.clearLocks")

	for {
		time.Sleep(p.p.lockingInterval)

		err := p.r.ClearLocks(ctx, p.instanceID, time.Now().Add(-p.p.maxLockAge))
		if err != nil {
			xlog.Infof("%s, error clearing locks: %+v", op, err)
		}
	}
}

// GetUnsentPublishings will return all unsent messages
func (p *PollableDataSource) GetUnsentPublishings(ctx context.Context) (<-chan *xmessage.Publishing, error) {
	messages := make(chan *xmessage.Publishing)

	go p.startPolling(ctx, messages)
	go p.clearLocks(ctx)

	return messages, nil
}

// SetAsProcessed will set the message as processed
func (p *PollableDataSource) SetAsProcessed(ctx context.Context, id string) error {
	op := xerrors.Op("outbox.PollableDataSource.SetAsProcessed")

	err := p.r.SetAsProcessed(ctx, id)
	if err != nil {
		return xerrors.E(op, err)
	}

	return nil
}

// RetryMessage will set failed message to be retried
func (p *PollableDataSource) RetryMessage(ctx context.Context, id string) error {
	op := xerrors.Op("outbox.PollableDataSource.RetryMessage")

	err := p.r.MarkForRetry(ctx, id, time.Now().Add(p.p.retryInterval))
	if err != nil {
		return xerrors.E(op, err)
	}

	return nil
}
