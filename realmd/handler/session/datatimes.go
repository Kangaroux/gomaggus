package session

import (
	"log"
	"time"

	"github.com/kangaroux/gomaggus/realmd"
)

// https://gtker.com/wow_messages/docs/smsg_account_data_times.html#client-version-335
type dataTimesResponse struct {
	Time      uint32
	Activated bool

	// CacheMask is the size of Cache but it's represented by how many bits are set, not the actual value.
	CacheMask uint32
	Cache     []uint32
}

func DataTimesHandler(client *realmd.Client) error {
	resp := dataTimesResponse{
		Time:      uint32(time.Now().Unix()),
		Activated: true,
		CacheMask: 0xFF,              // 8 bits = 8 cache values
		Cache:     make([]uint32, 8), // Can leave it as all zeroes for now
	}
	if err := client.SendPacket(realmd.OpServerAccountDataTimes, &resp); err != nil {
		return err
	}

	log.Println("sent account data times")
	return nil
}
