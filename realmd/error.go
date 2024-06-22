package realmd

type ErrKickClient struct {
	Reason string
}

func (e *ErrKickClient) Error() string {
	msg := "kicking client"

	if e.Reason != "" {
		msg += "(" + e.Reason + ")"
	}

	return msg
}
