package player

import (
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/kangaroux/gomaggus/realmd/player"
	"github.com/mixcode/binarystruct"
)

type standStateRequest struct {
	State player.StandState
	_     [3]byte `binary:"Pad"` // Client encodes state as 4 bytes but it's only 1 byte
}

type standStateResponse struct {
	State player.StandState
}

func StandStateHandler(client *realmd.Client, data []byte) error {
	req := standStateRequest{}
	if _, err := binarystruct.Unmarshal(data, binarystruct.LittleEndian, &req); err != nil {
		return err
	}

	switch req.State {
	case player.StateStand, player.StateSit, player.StateSleep, player.StateKneel:
		// TODO: store player state. for now, just mirror it back
		resp := standStateResponse{State: req.State}
		return client.SendPacket(realmd.OpServerStandState, &resp)

	default:
		// Ignore everything else, these states can only be set by the server
		return nil
	}
}
