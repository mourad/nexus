package aat

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gammazero/nexus/auth"
	"github.com/gammazero/nexus/client"
	"github.com/gammazero/nexus/wamp"
)

func TestJoinRealmWithCRAuth(t *testing.T) {
	// Connect callee session.
	cli, err := connectClientNoJoin()
	if err != nil {
		t.Fatal("Failed to connect client:", err)
	}

	details := wamp.Dict{
		"username": "jdoe", "authmethods": []string{"testauth"}}
	authMap := map[string]client.AuthFunc{"testauth": testAuthFunc}

	details, err = cli.JoinRealm("nexus.test.auth", details, authMap)
	if err != nil {
		t.Fatal(err)
	}

	if wamp.OptionString(details, "authrole") != "user" {
		t.Fatal("missing or incorrect authrole")
	}
}

func TestJoinRealmWithCRAuthBad(t *testing.T) {
	// Connect callee session.
	cli, err := connectClientNoJoin()
	if err != nil {
		t.Fatal("Failed to connect client:", err)
	}

	details := wamp.Dict{
		"username": "malory", "authmethods": []string{"testauth"}}
	authMap := map[string]client.AuthFunc{"testauth": testAuthFunc}

	_, err = cli.JoinRealm("nexus.test.auth", details, authMap)
	if err == nil {
		t.Fatal("expected error with bad username")
	}
	if !strings.HasSuffix(err.Error(), "invalid signature") {
		t.Fatal("wrong error message:", err)
	}
}

func testAuthFunc(d wamp.Dict, c wamp.Dict) (string, wamp.Dict) {
	ch := wamp.OptionString(c, "challenge")
	return testCRSign(ch), wamp.Dict{}
}

type testCRAuthenticator struct{}

// pendingTestAuth implements the PendingCRAuth interface.
type pendingTestAuth struct {
	authID string
	secret string
	role   string
}

func (t *testCRAuthenticator) Challenge(details wamp.Dict) (auth.PendingCRAuth, error) {
	username := wamp.OptionString(details, "username")
	if username == "" {
		return nil, errors.New("no username given")
	}
	var secret string
	if username == "jdoe" {
		secret = testCRSign(username)
	}

	return &pendingTestAuth{
		authID: username,
		role:   "user",
		secret: secret,
	}, nil
}

func testCRSign(uname string) string {
	return uname + "123xyz"
}

// Return the test challenge message.
func (p *pendingTestAuth) Msg() *wamp.Challenge {
	return &wamp.Challenge{
		AuthMethod: "testauth",
		Extra:      wamp.Dict{"challenge": p.authID},
	}
}

func (p *pendingTestAuth) Timeout() time.Duration { return time.Second }

func (p *pendingTestAuth) Authenticate(msg *wamp.Authenticate) (*wamp.Welcome, error) {

	if p.secret == "" || p.secret != msg.Signature {
		return nil, errors.New("invalid signature")
	}

	// Create welcome details containing auth info.
	details := wamp.Dict{
		"authid":     p.authID,
		"authmethod": "testauth",
		"authrole":   p.role,
	}

	return &wamp.Welcome{Details: details}, nil
}
