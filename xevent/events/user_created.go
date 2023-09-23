// Package events provides the events for the system
package events

import (
	"encoding/json"

	"github.com/Logistics-Coordinators/x/xerrors"
	"github.com/Logistics-Coordinators/x/xevent"
)

// TypeUserCreated is the identifier for events.UserCreated
const TypeUserCreated = "user/created"

// UserCreated is an event that is emitted when a user is created
type UserCreated struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Encode will encode the event
func (u *UserCreated) Encode() (xevent.Event, error) {
	op := xerrors.Op("events.UserCreated.Encode")

	b, err := json.Marshal(u)
	if err != nil {
		return xevent.Event{}, xerrors.E(op, err)
	}

	return xevent.Event{
		Type: TypeUserCreated,
		Data: b,
	}, nil
}

// Decode will decode the event
func (u *UserCreated) Decode(event xevent.Event) error {
	op := xerrors.Op("events.UserCreated.Decode")

	err := json.Unmarshal(event.Data, u)
	if err != nil {
		return xerrors.E(op, err)
	}

	return nil
}
