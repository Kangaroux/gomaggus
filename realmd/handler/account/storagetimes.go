package account

import (
	"time"

	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
)

// https://gtker.com/wow_messages/docs/smsg_account_data_times.html#client-version-335
type storageTimesResponse struct {
	Time         uint32
	Activated    bool
	StorageMask  uint32
	StorageTimes []uint32
}

func StorageTimesHandler(svc *realmd.Service, client *realmd.Client) error {
	// The client seems to expect all 8 storage types to be included in the times, except that it
	// only accesses the ones for the account (we can't know what char it is because they are still
	// on the char select screen).
	times := make([]uint32, model.AllStorageCount)
	storages, err := svc.AccountStorage.List(client.Account.Id, model.AllAccountStorage)
	if err != nil {
		return err
	}

	for _, s := range storages {
		times[s.Type] = uint32(s.UpdatedAt.Unix())
	}

	resp := storageTimesResponse{
		Time:         uint32(time.Now().Unix()),
		Activated:    true,
		StorageMask:  model.AllAccountStorage,
		StorageTimes: times,
	}
	return client.SendPacket(realmd.OpServerAccountStorageTimes, &resp)
}
