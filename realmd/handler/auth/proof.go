package auth

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/kangaroux/gomaggus/srp"
)

// https://gtker.com/wow_messages/docs/cmsg_auth_session.html#client-version-335
type authRequest struct {
	ClientBuild     uint32
	LoginServerId   uint32
	Username        string
	LoginServerType uint32
	ClientSeed      [4]byte
	RegionId        uint32
	BattlegroundId  uint32
	RealmId         uint32
	DOSResponse     uint64
	ClientProof     [20]byte
	AddonInfo       []byte
}

func ProofHandler(services *realmd.Service, client *realmd.Client, data []byte) error {
	log.Println("starting auth session")

	var err error
	r := bytes.NewReader(data[6:])

	// https://gtker.com/wow_messages/docs/cmsg_auth_session.html#client-version-335
	p := authRequest{}
	if err = binary.Read(r, binary.LittleEndian, &p.ClientBuild); err != nil {
		return err
	}
	if err = binary.Read(r, binary.LittleEndian, &p.LoginServerId); err != nil {
		return err
	}
	if p.Username, err = internal.ReadCString(r); err != nil {
		return err
	}
	if err = binary.Read(r, binary.LittleEndian, &p.LoginServerType); err != nil {
		return err
	}
	if err = binary.Read(r, binary.BigEndian, &p.ClientSeed); err != nil {
		return err
	}
	if err = binary.Read(r, binary.LittleEndian, &p.RegionId); err != nil {
		return err
	}
	if err = binary.Read(r, binary.LittleEndian, &p.BattlegroundId); err != nil {
		return err
	}
	if err = binary.Read(r, binary.LittleEndian, &p.RealmId); err != nil {
		return err
	}
	if err = binary.Read(r, binary.LittleEndian, &p.DOSResponse); err != nil {
		return err
	}
	if _, err = r.Read(p.ClientProof[:]); err != nil {
		return err
	}
	addonInfoBuf := bytes.Buffer{}
	if _, err = r.WriteTo(&addonInfoBuf); err != nil {
		return err
	}
	p.AddonInfo = addonInfoBuf.Bytes()

	client.Authenticated, err = authenticateClient(services, client, &p)
	if err != nil {
		return err
	}

	if !client.Authenticated {
		// We can't return an error to the client due to the header encryption, just drop the connection
		return errors.New("client could not be authenticated")
	}

	inner := bytes.Buffer{}
	inner.WriteByte(byte(realmd.RespCodeAuthOk))
	inner.Write([]byte{0, 0, 0, 0})              // billing time
	inner.WriteByte(0x0)                         // billing flags
	inner.Write([]byte{0, 0, 0, 0})              // billing rested
	inner.WriteByte(byte(realmd.ExpansionWrath)) // exp

	// https://gtker.com/wow_messages/docs/smsg_auth_response.html#client-version-335
	resp := bytes.Buffer{}
	respHeader, err := realmd.BuildHeader(realmd.OpServerAuthResponse, uint32(inner.Len()))
	if err != nil {
		return err
	}
	resp.Write(client.Crypto.Encrypt(respHeader))
	resp.Write(inner.Bytes())

	if _, err := client.Conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent auth response")

	return nil
}

func authenticateClient(svc *realmd.Service, client *realmd.Client, p *authRequest) (bool, error) {
	var err error

	client.Account, err = svc.Accounts.Get(&model.AccountGetParams{Username: p.Username})
	if err != nil {
		return false, err
	} else if client.Account == nil {
		log.Printf("no account with username %s exists", p.Username)
		return false, nil
	}

	if client.Realm, err = svc.Realms.Get(p.RealmId); err != nil {
		return false, err
	} else if client.Realm == nil {
		log.Printf("no realm with id %d exists", p.RealmId)
		return false, nil
	}

	if client.Session, err = svc.Sessions.Get(client.Account.Id); err != nil {
		return false, err
	} else if client.Session == nil {
		log.Printf("no session for username %s exists", client.Account.Username)
		return false, nil
	}

	if err := client.Session.Decode(); err != nil {
		return false, err
	}

	client.Crypto = realmd.NewHeaderCrypto(client.Session.SessionKey())
	if err := client.Crypto.Init(); err != nil {
		return false, err
	}

	proof := srp.CalculateWorldProof(p.Username, p.ClientSeed[:], client.ServerSeed[:], client.Session.SessionKey())

	if !bytes.Equal(proof, p.ClientProof[:]) {
		log.Println("proofs don't match")
		log.Printf("got:    %x\n", p.ClientProof)
		log.Printf("wanted: %x\n", proof)
		return false, nil
	}

	log.Println("client authenticated successfully")

	return true, nil
}
