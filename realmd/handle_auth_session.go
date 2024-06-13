package realmd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"

	"github.com/kangaroux/gomaggus/models"
	"github.com/kangaroux/gomaggus/srp"
)

func handleAuthSession(services *Services, client *Client, data []byte) error {
	log.Println("starting auth session")

	var err error
	r := bytes.NewReader(data[6:])

	// https://gtker.com/wow_messages/docs/cmsg_auth_session.html#client-version-335
	p := AuthSessionPacket{}
	if err = binary.Read(r, binary.LittleEndian, &p.ClientBuild); err != nil {
		return err
	}
	if err = binary.Read(r, binary.LittleEndian, &p.LoginServerId); err != nil {
		return err
	}
	if p.Username, err = readCString(r); err != nil {
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

	client.authenticated, err = authenticateClient(services, client, &p)
	if err != nil {
		return err
	}

	if !client.authenticated {
		// We can't return an error to the client due to the header encryption, just drop the connection
		return errors.New("client could not be authenticated")
	}

	inner := bytes.Buffer{}
	inner.WriteByte(byte(RespCodeAuthOk))
	inner.Write([]byte{0, 0, 0, 0})       // billing time
	inner.WriteByte(0x0)                  // billing flags
	inner.Write([]byte{0, 0, 0, 0})       // billing rested
	inner.WriteByte(byte(ExpansionWrath)) // exp

	// https://gtker.com/wow_messages/docs/smsg_auth_response.html#client-version-335
	resp := bytes.Buffer{}
	respHeader, err := makeServerHeader(OpServerAuthResponse, uint32(inner.Len()))
	if err != nil {
		return err
	}
	resp.Write(client.crypto.Encrypt(respHeader))
	resp.Write(inner.Bytes())

	if _, err := client.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent auth response")

	return nil
}

func authenticateClient(services *Services, client *Client, p *AuthSessionPacket) (bool, error) {
	var err error

	client.account, err = services.accounts.Get(&models.AccountGetParams{Username: p.Username})
	if err != nil {
		return false, err
	} else if client.account == nil {
		log.Printf("no account with username %s exists", p.Username)
		return false, nil
	}

	if client.realm, err = services.realms.Get(p.RealmId); err != nil {
		return false, err
	} else if client.realm == nil {
		log.Printf("no realm with id %d exists", p.RealmId)
		return false, nil
	}

	if client.session, err = services.sessions.Get(client.account.Id); err != nil {
		return false, err
	} else if client.session == nil {
		log.Printf("no session for username %s exists", client.account.Username)
		return false, nil
	}

	if err := client.session.Decode(); err != nil {
		return false, err
	}

	client.crypto = NewWrathHeaderCrypto(client.session.SessionKey())
	if err := client.crypto.Init(); err != nil {
		return false, err
	}

	proof := srp.CalculateWorldProof(p.Username, p.ClientSeed[:], client.serverSeed[:], client.session.SessionKey())

	if !bytes.Equal(proof, p.ClientProof[:]) {
		log.Println("proofs don't match")
		log.Printf("got:    %x\n", p.ClientProof)
		log.Printf("wanted: %x\n", proof)
		return false, nil
	}

	log.Println("client authenticated successfully")

	return true, nil
}
