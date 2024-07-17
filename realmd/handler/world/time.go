package world

import (
	"time"

	"github.com/kangaroux/gomaggus/realmd"
)

type uiTimeResponse struct {
	Time uint32
}

func UITimeHandler(client *realmd.Client) error {
	resp := uiTimeResponse{Time: uint32(time.Now().Unix())}
	return client.SendPacket(realmd.OpServerUITime, &resp)
}
