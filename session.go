package nexus

import (
	"fmt"

	"github.com/gammazero/nexus/wamp"
)

// Session is an active WAMP session.
type Session struct {
	wamp.Peer
	ID      wamp.ID
	Details wamp.Dict

	stop chan wamp.URI
}

// String returns the session ID as a string.
func (s Session) String() string { return fmt.Sprintf("%d", s.ID) }

// HasRole returns true if the session supports the specified role.
func (s Session) HasRole(role string) bool {
	_, err := wamp.DictValue(s.Details, []string{"roles", role})
	return err == nil
}

// HasFeature returns true if the session has the specified feature for the
// specified role.
func (s Session) HasFeature(role, feature string) bool {
	b, _ := wamp.DictFlag(s.Details, []string{"roles", role, "features", feature})
	return b
}
