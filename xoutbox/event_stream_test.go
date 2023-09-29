package xoutbox_test

import (
	"strings"

	"github.com/Logistics-Coordinators/dqf-v2-backend/internal/xoutbox"
	"github.com/Logistics-Coordinators/x/xerrors"
	"github.com/Logistics-Coordinators/x/xtest"
)

type testEventStream struct {
	messages map[string]*xoutbox.Message
}

func (es *testEventStream) Send(m *xoutbox.Message) error {
	if strings.HasPrefix(m.ID, "fail:") {
		return xerrors.E(xerrors.Message("failed to send message"))
	}

	if strings.HasPrefix(m.ID, "retry:") {
		// Fail 60% of the time
		if xtest.RandomInt(0, 100) < 60 {
			return xerrors.E(xerrors.Message("failed to send message"))
		}
	}

	es.messages[m.ID] = m
	return nil
}

func (es *testEventStream) isSent(id string) bool {
	_, ok := es.messages[id]
	return ok
}

func newTestEventStream() *testEventStream {
	return &testEventStream{
		messages: make(map[string]*xoutbox.Message),
	}
}
