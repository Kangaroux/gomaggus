package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"log"

	"github.com/kangaroux/gomaggus/authd"
	"github.com/kangaroux/gomaggus/model"
	"github.com/mixcode/binarystruct"
)

const (
	ReconnectDataLen = 16
)

type reconnectChallengeRequest = loginChallengeRequest

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_challenge_server.html#protocol-version-8
type reconnectChallengeResponse struct {
	Opcode        authd.Opcode // OpReconnectChallenge
	ErrorCode     authd.RespCode
	ReconnectData [ReconnectDataLen]byte
	ChecksumSalt  [16]byte
}

type ReconnectChallenge struct {
	Client   *authd.Client
	Accounts model.AccountService
}

func (h *ReconnectChallenge) Handle(data []byte) error {
	if h.Client.State != authd.StateAuthChallenge {
		return &ErrWrongState{
			Handler:  "RealmList",
			Expected: authd.StateAuthChallenge,
			Actual:   h.Client.State,
		}
	}

	log.Println("Starting reconnect challenge")

	req := reconnectChallengeRequest{}
	if _, err := binarystruct.Unmarshal(data, binarystruct.LittleEndian, &req); err != nil {
		return err
	}

	log.Printf("client trying to reconnect as '%s'\n", req.Username)

	acct, err := h.Accounts.Get(&model.AccountGetParams{Username: req.Username})
	if err != nil {
		return err
	}

	// Generate random data that will be used for the reconnect proof
	if _, err := rand.Read(h.Client.ReconnectData); err != nil {
		return err
	}

	resp := reconnectChallengeResponse{
		Opcode: authd.OpcodeReconnectChallenge,

		// Always return success to prevent a bad actor from mining usernames.
		ErrorCode:    authd.Success,
		ChecksumSalt: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	copy(resp.ReconnectData[:], h.Client.ReconnectData)

	respBuf := bytes.Buffer{}
	binary.Write(&respBuf, binary.BigEndian, &resp)

	if _, err := h.Client.Conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to reconnect challenge")

	if acct != nil {
		h.Client.Account = acct
		h.Client.Username = req.Username
	}

	h.Client.State = authd.StateReconnectProof

	return nil
}
