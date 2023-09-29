package outbox_test

import (
	"strings"

	"github.com/Logistics-Coordinators/x/xerrors"
	"github.com/Logistics-Coordinators/x/xmessage"
	"github.com/Logistics-Coordinators/x/xtest"
)

type testEventStream struct {
	publishings map[string]*xmessage.Publishing
}

func (es *testEventStream) Send(p *xmessage.Publishing) error {
	if strings.HasPrefix(string(p.Topic), "fail:") {
		return xerrors.E(xerrors.Message("failed to send message"))
	}

	if strings.HasPrefix(string(p.Topic), "retry:") {
		// Fail 60% of the time
		if xtest.RandomInt(0, 100) < 60 {
			return xerrors.E(xerrors.Message("failed to send message"))
		}
	}

	es.publishings[p.Message.ID] = p
	return nil
}

func (es *testEventStream) isSent(id string) bool {
	_, ok := es.publishings[id]
	return ok
}

func newTestEventStream() *testEventStream {
	return &testEventStream{
		publishings: make(map[string]*xmessage.Publishing),
	}
}
