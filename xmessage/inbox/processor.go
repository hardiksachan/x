package inbox

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

// HandleFunc is a function that handles a message
type HandleFunc func(ctx context.Context, message *xmessage.Message) error

// Processor is a idempotent inbox
type Processor struct {
	r      Repository
	types  []string
	id     string
	handle HandleFunc

	pollingInterval time.Duration
	lockingInterval time.Duration
	maxLockAge      time.Duration
	maxRetries      int
	retryInterval   time.Duration
}

// ProcessorOption is a function that configures a Processor
type ProcessorOption func(*Processor)

// WithPollingInterval sets the polling interval
func WithPollingInterval(interval time.Duration) ProcessorOption {
	return func(i *Processor) {
		i.pollingInterval = interval
	}
}

// WithLocking sets locking interval and age
func WithLocking(interval time.Duration, age time.Duration) ProcessorOption {
	return func(i *Processor) {
		i.lockingInterval = interval
		i.maxLockAge = age
	}
}

// WithRetries sets the parameters for retry
func WithRetries(maxRetries int, retryInterval time.Duration) ProcessorOption {
	return func(i *Processor) {
		i.maxRetries = maxRetries
		i.retryInterval = retryInterval
	}
}

// NewProcessor creates a new inbox.Processor
func NewProcessor(r Repository, types []string, handler HandleFunc, opts ...ProcessorOption) Processor {
	i := Processor{
		r:      r,
		types:  types,
		id:     uuid.NewString(),
		handle: handler,

		pollingInterval: defaultPollingInterval,
		lockingInterval: defaultLockInterval,
		maxLockAge:      defaultMaxLockAge,
		maxRetries:      defaultMaxRetries,
		retryInterval:   defaultRetryInterval,
	}

	for _, opt := range opts {
		opt(&i)
	}

	return i
}

// Start will start the processor
func (p *Processor) Start(ctx context.Context) {
	go p.processMessages(ctx)
	go p.clearLocks(ctx)
}

func (p *Processor) processMessages(ctx context.Context) {
	for {
		time.Sleep(p.pollingInterval)

		message, err := p.r.GetUnprocessedMessage(ctx, p.id, p.maxRetries, p.types)
		if err != nil {
			continue
		}

		err = p.handle(ctx, message)
		if err != nil {
			xlog.Infof("error handling message %s: %+v", message.Type, err)
			_ = p.r.MarkForRetry(ctx, message.ID, time.Now().Add(p.retryInterval))
			continue
		}

		_ = p.r.SetAsProcessed(ctx, message.ID)
	}
}

func (p *Processor) clearLocks(ctx context.Context) {
	op := xerrors.Op("inbox.Processor.clearLocks")

	for {
		time.Sleep(p.lockingInterval)

		err := p.r.ClearLocks(ctx, p.id, time.Now().Add(-p.maxLockAge))
		if err != nil {
			xlog.Infof("%s, error clearing locks: %+v", op, err)
		}
	}
}
