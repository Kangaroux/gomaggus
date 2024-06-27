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
	times := make([]uint32, model.AccountStorageCount)
	storages, err := svc.AccountStorage.List(client.Account.Id, model.AllAccountStorage)
	if err != nil {
		return err
	}

	for i, s := range storages {
		times[i] = uint32(s.UpdatedAt.Unix())
	}

	resp := storageTimesResponse{
		Time:         uint32(time.Now().Unix()),
		Activated:    true,
		StorageMask:  model.AllAccountStorage,
		StorageTimes: times,
	}
	return client.SendPacket(realmd.OpServerAccountStorageTimes, &resp)
}
