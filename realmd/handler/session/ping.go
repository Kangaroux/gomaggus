package session

import (
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/mixcode/binarystruct"
)

// https://gtker.com/wow_messages/docs/cmsg_ping.html#client-version-19-client-version-110-client-version-111-client-version-112-client-version-2-client-version-3
type pingRequest struct {
	SequenceId    uint32
	RoundTripTime uint32 // zero if server hasn't responded?
}

// https://gtker.com/wow_messages/docs/smsg_pong.html
type pingResponse struct {
	SequenceId uint32
}

func PingHandler(client *realmd.Client, data []byte) error {
	req := pingRequest{}
	if _, err := binarystruct.Unmarshal(data, binarystruct.LittleEndian, &req); err != nil {
		return err
	}

	resp := pingResponse{SequenceId: req.SequenceId}
	return client.SendPacket(realmd.OpServerPong, &resp)
}
