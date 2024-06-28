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

type StorageTimesHandler struct {
	Client  *realmd.Client
	Service *realmd.Service
}

func (h *StorageTimesHandler) Handle() error {
	return h.Send(model.AllAccountStorage)
}

func (h *StorageTimesHandler) Send(mask uint8) error {
	// The client expects every storage time is sent, regardless of what the StorageMask is. When the
	// client is on the character select screen, any character storage times will always be zero.
	times := make([]uint32, model.AllStorageCount)

	// Load account storage times
	if mask&model.AllAccountStorage > 0 {
		storages, err := h.Service.AccountStorage.List(h.Client.Account.Id, mask)
		if err != nil {
			return err
		}

		for _, s := range storages {
			times[s.Type] = uint32(s.UpdatedAt.Unix())
		}
	}

	// Load character storage times (if they are logged in to the world)
	if h.Client.Character != nil && mask&model.AllCharacterStorage > 0 {
		storages, err := h.Service.CharacterStorage.List(h.Client.Character.Id, mask)
		if err != nil {
			return err
		}

		for _, s := range storages {
			times[s.Type] = uint32(s.UpdatedAt.Unix())
		}
	}

	resp := storageTimesResponse{
		Time:         uint32(time.Now().Unix()),
		Activated:    true,
		StorageMask:  uint32(mask),
		StorageTimes: times,
	}
	return h.Client.SendPacket(realmd.OpServerClientStorageTimes, &resp)
}
