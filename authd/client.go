package authd

import (
	"io"

	"github.com/kangaroux/gomaggus/model"
)

type ClientState int

const (
	// Initial state. Waiting for the client to send an auth or reconnect challenge.
	StateAuthChallenge ClientState = iota

	// Waiting for client to send auth proof.
	StateAuthProof

	// Waiting for client to send reconnect proof.
	StateReconnectProof

	// Client is fully authenticated.
	StateAuthenticated

	// Client failed to authenticate or an error occurred.
	StateInvalid
)

type Client struct {
	// Using ReadWriteCloser instead of net.Conn results in cleaner test mocks.
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
