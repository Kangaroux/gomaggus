package account

import (
	"log"
	"time"

	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
)

// https://gtker.com/wow_messages/docs/smsg_account_data_times.html#client-version-335
type storageTimesResponse struct {
	Time         uint32
	Activated    bool
	StorageMask  model.StorageMask `binary:"uint32"` // encoded as 4 bytes
	StorageTimes []uint32
}

func StorageTimesHandler(svc *realmd.Service, client *realmd.Client) error {
	times := make([]uint32, 8)
	storages, err := svc.AccountStorage.List(client.Account.Id, model.StorageMaskAll)
	if err != nil {
		return err
	}

	for _, s := range storages {
		// s.Type maps nicely to the order the client is expecting.
		times[s.Type] = uint32(s.UpdatedAt.Unix())
	}

	resp := storageTimesResponse{
		Time:         uint32(time.Now().Unix()),
		Activated:    true,
		StorageMask:  model.StorageMaskAll,
		StorageTimes: times,
	}
	if err := client.SendPacket(realmd.OpServerAccountStorageTimes, &resp); err != nil {
		return err
	}

	log.Println("sent account storage times")
	return nil
}
