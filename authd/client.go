package authd

import (
	"io"

	"github.com/kangaroux/gomaggus/model"
)

type ClientState int

const (
	// Initial state, waiting for the client to send an auth or reconnect challenge
	StateAuthChallenge ClientState = iota

	// We've responded to the auth challenge and are waiting for the client's proof
	StateAuthProof

	// Waiting for client to send proof in order to reconnect
	StateReconnectProof

	// Client has authenticated successfully
	StateAuthenticated

	// Client failed to authenticate or an error occurred, and the connection should be closed
	StateInvalid
)

type Client struct {
	// Using ReadWriteCloser instead of net.Conn makes test mocks simpler
	Conn            io.ReadWriteCloser
	Username        string
	ReconnectData   []byte
	SessionKey      []byte
	ClientPublicKey []byte
	ServerPublicKey []byte
	PrivateKey      []byte
	State           ClientState
	Account         *model.Account
}
