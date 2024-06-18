package realm

import (
	"log"

	"github.com/kangaroux/gomaggus/realmd"
	"github.com/mixcode/binarystruct"
)

type RealmSplitState uint32

const (
	SplitNormal    RealmSplitState = 0
	SplitConfirmed RealmSplitState = 1
	SplitPending   RealmSplitState = 2
)

// https://gtker.com/wow_messages/docs/cmsg_realm_split.html
type splitRequest struct {
	RealmId uint32
}

// https://gtker.com/wow_messages/docs/smsg_realm_split.html
type splitResponse struct {
	RealmId   uint32
	State     RealmSplitState
	SplitDate string `binary:"zstring"`
}

func SplitInfoHandler(client *realmd.Client, data *realmd.ClientPacket) error {
	req := splitRequest{}
	if _, err := binarystruct.Unmarshal(data.Payload, binarystruct.LittleEndian, &req); err != nil {
		return err
	}

	resp := splitResponse{
		RealmId:   req.RealmId,
		State:     SplitNormal,
		SplitDate: "01/01/01",
	}
	if err := client.SendPacket(realmd.OpServerRealmSplit, &resp); err != nil {
		return err
	}

	log.Println("sent realm split")
	return nil
}
