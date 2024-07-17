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

type serverTimeResponse struct {
	Time               uint32
	NextDailyResetTime uint32
}

func ServerTimeHandler(client *realmd.Client) error {
	resp := serverTimeResponse{
		Time:               uint32(time.Now().Unix()),
		NextDailyResetTime: 0,
	}
	return client.SendPacket(realmd.OpServerTime, &resp)
}
