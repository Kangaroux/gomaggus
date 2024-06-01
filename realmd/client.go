package realmd

import "net"

type ClientState int

const (
	// Initial state, waiting for the client to send an auth or reconnect challenge
	StateAuthChallenge = iota

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
	conn          net.Conn
	username      string
	reconnectData []byte
	sessionKey    []byte
	state         ClientState
}
