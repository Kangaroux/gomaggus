package handler

import (
	"fmt"

	"github.com/kangaroux/gomaggus/authd"
)

type ErrPacketUnreadBytes struct {
	Handler     string
	UnreadCount int
}

func (e *ErrPacketUnreadBytes) Error() string {
	return fmt.Sprintf(
		"%s: parsed packet but there's still %d unread bytes",
		e.Handler, e.UnreadCount,
	)
}

type ErrWrongState struct {
	Handler  string
	Expected authd.ClientState
	Actual   authd.ClientState
}

func (e *ErrWrongState) Error() string {
	return fmt.Sprintf(
		"%s: client state does not match the required state (wanted %x, got %x)",
		e.Handler, e.Expected, e.Actual,
	)
}
