package xoutbox_test

import (
	"context"
	"sync"
	"time"

	"github.com/Logistics-Coordinators/x/xoutbox"
)

type messageWithStatus struct {
	message   *xoutbox.Message
	processed bool
}

type testDataStore struct {
	messages map[string]*messageWithStatus
	sync.RWMutex
}

func (ds *testDataStore) GetUnsentMessages(_ context.Context) (<-chan *xoutbox.Message, error) {
	msgChan := make(chan *xoutbox.Message)

	go func() {
		for {
			ds.RLock()
			for m := range ds.messages {
				if !ds.messages[m].processed {
					msgChan <- ds.messages[m].message
				}
			}
			ds.RUnlock()

			time.Sleep(time.Millisecond * 100)
		}
	}()

	return msgChan, nil
}

func (ds *testDataStore) SetAsProcessed(_ context.Context, id string) error {
	ds.Lock()
	defer ds.Unlock()

	ds.messages[id].processed = true
	return nil
}

func (ds *testDataStore) isProcessed(id string) bool {
	ds.RLock()
	defer ds.RUnlock()

	return ds.messages[id].processed
}

func (ds *testDataStore) AddMessage(m *xoutbox.Message) {
	ds.messages[m.ID] = &messageWithStatus{
		message:   m,
		processed: false,
	}
}

func newTestDataStore() *testDataStore {
	return &testDataStore{
		messages: make(map[string]*messageWithStatus),
		RWMutex:  sync.RWMutex{},
	}
}
