package realm

import (
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/mixcode/binarystruct"
)

type realmSplitState uint32

const (
	splitNormal    realmSplitState = 0
	splitConfirmed realmSplitState = 1
	splitPending   realmSplitState = 2
)

// https://gtker.com/wow_messages/docs/cmsg_realm_split.html
type splitRequest struct {
	RealmId uint32
}

// https://gtker.com/wow_messages/docs/smsg_realm_split.html
type splitResponse struct {
	RealmId   uint32
	State     realmSplitState
	SplitDate string `binary:"zstring"`
}

func SplitInfoHandler(client *realmd.Client, data []byte) error {
	req := splitRequest{}
	if _, err := binarystruct.Unmarshal(data, binarystruct.LittleEndian, &req); err != nil {
		return err
	}

	resp := splitResponse{
		RealmId:   req.RealmId,
		State:     splitNormal,
		SplitDate: "01/01/01",
	}
	return client.SendPacket(realmd.OpServerRealmSplit, &resp)
}
