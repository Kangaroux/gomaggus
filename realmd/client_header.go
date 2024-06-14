package realmd

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type ClientHeader struct {
	Size   uint16
	Opcode ClientOpcode
}

func ParseClientHeader(c *Client, data []byte) (*ClientHeader, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("parseHeader: payload should be at least 6 bytes but it's only %d", len(data))
	}

	headerData := data[:6]

	if c.Authenticated {
		if c.Crypto == nil {
			return nil, errors.New("parseHeader: client is authenticated but client.crypto is nil")
		}

		headerData = c.Crypto.Decrypt(headerData)
	}

	h := &ClientHeader{
		Size:   binary.BigEndian.Uint16(headerData[:2]),
		Opcode: ClientOpcode(binary.LittleEndian.Uint32(headerData[2:6])),
	}

	return h, nil
}
