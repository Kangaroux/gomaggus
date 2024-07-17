package time

import (
	"time"

	"github.com/kangaroux/gomaggus/realmd"
)

type unixTimeResponse struct {
	Time uint32
}

func UnixTimeHandler(client *realmd.Client) error {
	resp := unixTimeResponse{Time: uint32(time.Now().Unix())}
	return client.SendPacket(realmd.OpServerUITime, &resp)
}
