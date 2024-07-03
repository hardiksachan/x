// Package xretry provides a retry mechanism
package xretry

import (
	"time"

	"github.com/hardiksachan/x/xerrors"
)

// RetryPolicy is the retry policy
type RetryPolicy struct {
	immediateRetries   int
	retriesWithBackoff int
	delay              time.Duration
	backoffFactor      float64
}

// RetryPolicyOption is the option for the retry policy
type RetryPolicyOption func(*RetryPolicy)

// WithImmediateRetries sets the immediate retries
func WithImmediateRetries(retries int) RetryPolicyOption {
	return func(p *RetryPolicy) {
		p.immediateRetries = retries
	}
}

// WithRetriesWithBackoff sets the retries with backoff
func WithRetriesWithBackoff(retries int, delay time.Duration, backoffFactor float64) RetryPolicyOption {
	return func(p *RetryPolicy) {
		p.retriesWithBackoff = retries
		p.delay = delay
		p.backoffFactor = backoffFactor
	}
}

// WithNoRetries sets the retries to 0
func WithNoRetries() RetryPolicyOption {
	return func(p *RetryPolicy) {
		p.immediateRetries = 0
		p.retriesWithBackoff = 0
	}
}

// NewRetryPolicy creates a new RetryPolicy
func NewRetryPolicy(opts ...RetryPolicyOption) RetryPolicy {
	p := RetryPolicy{
		immediateRetries:   0,
		retriesWithBackoff: 0,
		delay:              0,
		backoffFactor:      0,
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}

// Retrier is the interface that wraps the Retry method
type Retrier struct {
	p RetryPolicy
}

// NewRetrier creates a new Retrier
func NewRetrier(p RetryPolicy) *Retrier {
	return &Retrier{
		p: p,
	}
}

// Retry will retry the given function
func (r *Retrier) Retry(f func() error) error {
	op := xerrors.Op("outbox.Retrier.Retry")
	err := immediatelyRetry(f, r.p.immediateRetries)
	if err != nil {
		err = retryWithBackoff(f, r.p.retriesWithBackoff, r.p.delay, r.p.backoffFactor)
		if err != nil {
			return xerrors.E(op, err)
		}
	}

	return nil
}

func immediatelyRetry(f func() error, retriesLeft int) error {
	err := f()
	if err == nil {
		return nil
	}

	if retriesLeft == 0 {
		return err
	}

	return immediatelyRetry(f, retriesLeft-1)
}

func retryWithBackoff(f func() error, retriesLeft int, delay time.Duration, backoff float64) error {
	err := f()
	if err == nil {
		return nil
	}

	if retriesLeft == 0 {
		return err
	}

	time.Sleep(delay)
	return retryWithBackoff(f, retriesLeft-1, time.Duration(float64(delay)*backoff), backoff)
}
