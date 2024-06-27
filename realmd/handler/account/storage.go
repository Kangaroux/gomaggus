package account

import (
	"encoding/binary"

	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/mixcode/binarystruct"
)

const (
	// TrinityCore uses this as the max uncompressed size
	MaxUncompressedStorageSize = 1 << 16
)

// https://gtker.com/wow_messages/docs/cmsg_update_account_data.html#client-version-335
type putStorageRequest struct {
	Type             uint32
	Time             uint32
	UncompressedSize uint32
	Data             []byte
}

// https://gtker.com/wow_messages/docs/smsg_update_account_data_complete.html
type putStorageResponse struct {
	Type    uint32
	Unknown uint32
}

type StoragePutHandler struct {
	Client  *realmd.Client
	Service *realmd.Service
	req     putStorageRequest
}

func (h *StoragePutHandler) Handle(data []byte) error {
	n, err := binarystruct.Unmarshal(data, binary.LittleEndian, &h.req)
	if err != nil {
		return err
	}

	h.req.Data = data[n:]

	if h.req.UncompressedSize > MaxUncompressedStorageSize {
		h.Client.Log.Warn().Uint32("size", h.req.UncompressedSize).Msg("storage data is too large")
		return &realmd.ErrKickClient{Reason: "storage too big"}
	}

	if t := model.AccountStorageType(h.req.Type); t.IsAAccountStorageType() {
		if err := h.putAccountStorage(t); err != nil {
			return err
		}
	} else if t := model.CharacterStorageType(h.req.Type); t.IsACharacterStorageType() {
		if err := h.putCharacterStorage(t); err != nil {
			return err
		}
	} else {
		h.Client.Log.Warn().Uint32("storageType", h.req.Type).Msg("tried updating unknown storage type")
		return &realmd.ErrKickClient{Reason: "invalid storage type"}
	}

	return h.reply()
}

func (h *StoragePutHandler) putAccountStorage(t model.AccountStorageType) error {
	// TODO: auth middleware
	if h.Client.Account == nil {
		return &realmd.ErrKickClient{Reason: "not authenticated"}
	}

	obj := model.AccountStorage{
		AccountId:        h.Client.Account.Id,
		Type:             t,
		Data:             h.req.Data,
		UncompressedSize: int(h.req.UncompressedSize),
	}

	created, err := h.Service.AccountStorage.UpdateOrCreate(&obj)
	if err != nil {
		return err
	}

	h.Client.Log.Trace().
		Int("len", len(h.req.Data)).
		Str("type", t.String()).
		Bool("created", created).
		Msg("updated account storage")

	return nil
}

func (h *StoragePutHandler) putCharacterStorage(t model.CharacterStorageType) error {
	// TODO: auth middleware
	if h.Client.Character == nil {
		return &realmd.ErrKickClient{Reason: "not playing"}
	}

	obj := model.CharacterStorage{
		CharacterId:      h.Client.Character.Id,
		Type:             t,
		Data:             h.req.Data,
		UncompressedSize: int(h.req.UncompressedSize),
	}

	created, err := h.Service.CharacterStorage.UpdateOrCreate(&obj)
	if err != nil {
		return err
	}

	h.Client.Log.Trace().
		Int("len", len(h.req.Data)).
		Str("type", t.String()).
		Str("char", h.Client.Character.String()).
		Bool("created", created).
		Msg("updated character storage")

	return nil
}

func (h *StoragePutHandler) reply() error {
	resp := putStorageResponse{
		Type:    h.req.Type,
		Unknown: 0, // TrinityCore hardcodes this as zero
	}
	return h.Client.SendPacket(realmd.OpServerPutStorageOK, &resp)
}

// https://gtker.com/wow_messages/docs/cmsg_request_account_data.html
type getStorageRequest struct {
	Type uint32
}

// https://gtker.com/wow_messages/docs/smsg_update_account_data.html
type getStorageResponse struct {
	Type             uint32
	UncompressedSize uint32
	Data             []byte
}

type StorageGetHandler struct {
	Client  *realmd.Client
	Service *realmd.Service
	req     getStorageRequest
}

func (h *StorageGetHandler) Handle(data []byte) error {
	if _, err := binarystruct.Unmarshal(data, binary.LittleEndian, &h.req); err != nil {
		return err
	}

	var err error
	var storageData []byte
	var uncompressedSize int

	if t := model.AccountStorageType(h.req.Type); t.IsAAccountStorageType() {
		uncompressedSize, storageData, err = h.getAccountStorage(t)
		if err != nil {
			return err
		}
	} else if t := model.CharacterStorageType(h.req.Type); t.IsACharacterStorageType() {
		uncompressedSize, storageData, err = h.getCharacterStorage(t)
		if err != nil {
			return err
		}
	} else {
		h.Client.Log.Warn().Uint32("storageType", h.req.Type).Msg("tried fetching unknown storage type")
		return &realmd.ErrKickClient{Reason: "invalid storage type"}
	}

	return h.reply(uncompressedSize, storageData)
}

func (h *StorageGetHandler) getAccountStorage(t model.AccountStorageType) (int, []byte, error) {
	// TODO: auth middleware
	if h.Client.Account == nil {
		return 0, nil, &realmd.ErrKickClient{Reason: "not authenticated"}
	}

	storage, err := h.Service.AccountStorage.Get(h.Client.Account.Id, t)
	if err != nil {
		return 0, nil, err
	}

	return storage.UncompressedSize, storage.Data, nil
}
func (h *StorageGetHandler) getCharacterStorage(t model.CharacterStorageType) (int, []byte, error) {
	// TODO: auth middleware
	if h.Client.Character == nil {
		return 0, nil, &realmd.ErrKickClient{Reason: "not playing"}
	}

	storage, err := h.Service.CharacterStorage.Get(h.Client.Character.Id, t)
	if err != nil {
		return 0, nil, err
	}

	return storage.UncompressedSize, storage.Data, nil
}

func (h *StorageGetHandler) reply(size int, data []byte) error {
	resp := getStorageResponse{
		Type:             h.req.Type,
		UncompressedSize: uint32(size),
		Data:             data,
	}
	return h.Client.SendPacket(realmd.OpServerGetStorage, &resp)
}
