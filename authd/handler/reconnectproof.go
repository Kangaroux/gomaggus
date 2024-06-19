package handler

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"time"

	"github.com/kangaroux/gomaggus/authd"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/srp"
	"github.com/mixcode/binarystruct"
)

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_proof_client.html
type reconnectProofRequest struct {
	Opcode         authd.Opcode // OpReconnectProof
	ProofData      [srp.ProofDataSize]byte
	ClientProof    [srp.ProofSize]byte
	ClientChecksum [20]byte
	KeyCount       uint8
}

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_proof_server.html#protocol-version-8
type reconnectProofResponse struct {
	Opcode    authd.Opcode
	ErrorCode authd.RespCode
	_         [2]byte // padding
}

type ReconnectProof struct {
	Client   *authd.Client
	Sessions model.SessionService
	request  *reconnectProofRequest
}

func (h *ReconnectProof) Handle() error {
	if h.Client.State != authd.StateAuthProof {
		return &ErrWrongState{
			Handler:  "RealmList",
			Expected: authd.StateAuthProof,
			Actual:   h.Client.State,
		}
	}

	log.Println("Starting reconnect proof")

	authenticated := false

	if h.Client.Account != nil {
		session, err := h.Sessions.Get(h.Client.Account.Id)
		if err != nil {
			return err
		}

		if session != nil {
			if err := session.Decode(); err != nil {
				return err
			}
			h.Client.SessionKey = session.SessionKey()

			serverProof := srp.CalculateReconnectProof(h.Client.Username, h.request.ProofData[:], h.Client.ReconnectData, h.Client.SessionKey)
			authenticated = bytes.Equal(serverProof, h.request.ClientProof[:])
		}
	}

	resp := reconnectProofResponse{Opcode: authd.OpcodeReconnectProof}

	if !authenticated {
		resp.ErrorCode = authd.UnknownAccount
	} else {
		resp.ErrorCode = authd.Success
	}

	respBuf := bytes.Buffer{}
	binary.Write(&respBuf, binary.BigEndian, &resp)

	if _, err := h.Client.Conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to reconnect proof")

	if authenticated {
		session := model.Session{
			AccountId:     h.Client.Account.Id,
			SessionKeyHex: hex.EncodeToString(h.Client.SessionKey),
			Connected:     1,
			ConnectedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		}
		if _, err := h.Sessions.UpdateOrCreate(&session); err != nil {
			return err
		}
		h.Client.State = authd.StateAuthenticated
	} else {
		h.Client.State = authd.StateInvalid
	}

	return nil
}

// Read reads the packet data and parses it as a reconnect proof request. If data is too small then
// Read returns ErrPacketReadEOF.
func (h *ReconnectProof) Read(data []byte) (int, error) {
	n, err := binarystruct.Unmarshal(data, binary.LittleEndian, &h.request)

	if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
		return 0, ErrPacketReadEOF
	} else if err != nil {
		return 0, err
	}

	return n, nil
}
