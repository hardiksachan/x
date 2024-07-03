package inbox

import (
	"context"
	"time"

	"github.com/hardiksachan/x/xmessage"
)

// Repository is the interface that provides a repository for messages
type Repository interface {
	SaveMessage(ctx context.Context, message *xmessage.Message) error
	GetUnprocessedMessage(ctx context.Context, instanceID string, maxRetries int, allowedTypes []string) (*xmessage.Message, error)
	SetAsProcessed(ctx context.Context, id string) error
	MarkForRetry(ctx context.Context, id string, retryAt time.Time) error
	ClearLocks(ctx context.Context, instanceID string, obtainedBefore time.Time) error
}
