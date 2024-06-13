package authd

import (
	"fmt"
)

type ErrPacketUnreadBytes struct {
	What  string
	Count int
}

func (e *ErrPacketUnreadBytes) Error() string {
	return fmt.Sprintf(
		"%s: parsed packet but there's still %d unread bytes",
		e.What,
		e.Count,
	)
}
