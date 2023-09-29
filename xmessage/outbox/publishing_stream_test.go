package outbox_test

import (
	"context"
	"strings"

	"github.com/Logistics-Coordinators/x/xerrors"
	"github.com/Logistics-Coordinators/x/xmessage"
	"github.com/Logistics-Coordinators/x/xtest"
)

type mockPublishingStream struct {
	publishings map[string]*xmessage.Publishing
}

func (ps *mockPublishingStream) Send(_ context.Context, p *xmessage.Publishing) error {
	if strings.HasPrefix(string(p.Topic), "fail:") {
		return xerrors.E(xerrors.Message("failed to send message"))
	}

	if strings.HasPrefix(string(p.Topic), "retry:") {
		// Fail 60% of the time
		if xtest.RandomInt(0, 100) < 60 {
			return xerrors.E(xerrors.Message("failed to send message"))
		}
	}

	ps.publishings[p.Message.ID] = p
	return nil
}

func (ps *mockPublishingStream) isSent(id string) bool {
	_, ok := ps.publishings[id]
	return ok
}

func newMockPublishingStream() *mockPublishingStream {
	return &mockPublishingStream{
		publishings: make(map[string]*xmessage.Publishing),
	}
}
