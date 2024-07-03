package outbox_test

import (
	"context"
	"sync"
	"time"

	"github.com/hardiksachan/x/xmessage"
)

type publishingWithStatus struct {
	publishing *xmessage.Publishing
	processed  bool
}

type testDataStore struct {
	publishings map[string]*publishingWithStatus
	sync.RWMutex
}

func (ds *testDataStore) GetUnsentPublishings(_ context.Context) (<-chan *xmessage.Publishing, error) {
	pubChan := make(chan *xmessage.Publishing)

	go func() {
		for {
			ds.RLock()
			for p := range ds.publishings {
				if !ds.publishings[p].processed {
					pubChan <- ds.publishings[p].publishing
				}
			}
			ds.RUnlock()

			time.Sleep(time.Millisecond * 100)
		}
	}()

	return pubChan, nil
}

func (ds *testDataStore) SetAsProcessed(_ context.Context, id string) error {
	ds.Lock()
	defer ds.Unlock()

	ds.publishings[id].processed = true
	return nil
}

func (ds *testDataStore) isProcessed(id string) bool {
	ds.RLock()
	defer ds.RUnlock()

	return ds.publishings[id].processed
}

func (ds *testDataStore) AddPublishing(p *xmessage.Publishing) {
	ds.publishings[p.Message.ID] = &publishingWithStatus{
		publishing: p,
		processed:  false,
	}
}

func newTestDataStore() *testDataStore {
	return &testDataStore{
		publishings: make(map[string]*publishingWithStatus),
		RWMutex:     sync.RWMutex{},
	}
}
