package handler

import (
	"fmt"

	"github.com/kangaroux/gomaggus/authd"
)

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
