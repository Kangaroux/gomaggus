package player

import (
	"encoding/binary"

	"github.com/kangaroux/gomaggus/realmd"
	"github.com/mixcode/binarystruct"
)

type playedTimeRequest struct {
	ShowInChat bool
}

type playedTimeResponse struct {
	TotalTime  uint32
	LevelTime  uint32
	ShowInChat bool
}

func PlayedTimeHandler(client *realmd.Client, data []byte) error {
	req := playedTimeRequest{}
	if _, err := binarystruct.Unmarshal(data, binary.LittleEndian, &req); err != nil {
		return err
	}

	resp := playedTimeResponse{
		TotalTime:  100, // TODO
		LevelTime:  10,  // TODO
		ShowInChat: req.ShowInChat,
	}
	return client.SendPacket(realmd.OpServerPlayedTime, &resp)
}
