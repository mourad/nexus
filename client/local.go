package client

import (
	"time"

	"github.com/gammazero/nexus"
	"github.com/gammazero/nexus/stdlog"
	"github.com/gammazero/nexus/transport"
)

// NewLocalClient creates a new client directly connected to embedded router
//
// JoinRealm must be called before other client functions.
func NewLocalClient(router nexus.Router, responseTimeout time.Duration, logger stdlog.StdLog) (*Client, error) {
	localSide, routerSide := transport.LinkedPeers(logger)

	go func() {
		if err := router.Attach(routerSide); err != nil {
			logger.Print(err)
		}
	}()

	return NewClient(localSide, responseTimeout, logger), nil
}
