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

// https://gtker.com/wow_messages/docs/billingplanflags.html
type BillingFlag uint8

const (
	BillingNone          BillingFlag = 0x00
	BillingUnused        BillingFlag = 0x01
	BillingRecurringBill BillingFlag = 0x02
	BillingFreeTrial     BillingFlag = 0x04
	BillingIgr           BillingFlag = 0x08
	BillingUsage         BillingFlag = 0x10
	BillingTimeMixture   BillingFlag = 0x20
	BillingRestricted    BillingFlag = 0x40
	BillingEnableCais    BillingFlag = 0x80
)

// https://gtker.com/wow_messages/docs/cmsg_auth_session.html#client-version-335
type proofRequest struct {
	ClientBuild     uint32
	LoginServerId   uint32
	Username        string `binary:"zstring"`
	LoginServerType uint32
	ClientSeed      [4]byte
	RegionId        uint32
	BattlegroundId  uint32
	RealmId         uint32
	DOSResponse     uint64
	ClientProof     [20]byte
	AddonInfo       []byte
}

// https://gtker.com/wow_messages/docs/smsg_auth_response.html#client-version-335
type proofSuccess struct {
	ResponseCode  realmd.ResponseCode
	BillingTime   uint32
	BillingFlags  BillingFlag
	BillingRested uint32
	Expansion     realmd.Expansion
}

// Unused for now
// type proofWaitQueue struct {
// 	ResponseCode         realmd.ResponseCode
// 	HasFreeCharMigration bool
// }

func ProofHandler(svc *realmd.Service, client *realmd.Client, data *realmd.ClientPacket) error {
	log.Println("starting auth session")

	r := bytes.NewReader(data.Payload)

	p := proofRequest{}
	if err := binary.Read(r, binary.LittleEndian, &p.ClientBuild); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &p.LoginServerId); err != nil {
		return err
	}
	username, err := internal.ReadCString(r)
	if err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &p.LoginServerType); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &p.ClientSeed); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &p.RegionId); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &p.BattlegroundId); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &p.RealmId); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &p.DOSResponse); err != nil {
		return err
	}
	if _, err := r.Read(p.ClientProof[:]); err != nil {
		return err
	}
	addonInfoBuf := bytes.Buffer{}
	if _, err := r.WriteTo(&addonInfoBuf); err != nil {
		return err
	}
	p.AddonInfo = addonInfoBuf.Bytes()
	p.Username = username

	authenticated, err := authenticateClient(svc, client, &p)
	if err != nil {
		return err
	}
	client.Authenticated = authenticated

	if !client.Authenticated {
		// The client expects the authentication to be successful and the header to be encrypted.
		// If auth failed, we don't know how to encrypt the header, thus we can't send an error response.
		// Just drop the connection.
		return errors.New("client could not be authenticated")
	}

	// Header crypto can be initialized now that we know the session key
	headerCrypto := realmd.NewHeaderCrypto(client.Session.SessionKey())
	if err := headerCrypto.Init(); err != nil {
		return err
	}
	client.HeaderCrypto = headerCrypto

	resp := proofSuccess{
		ResponseCode:  realmd.RespCodeAuthOk,
		BillingTime:   0,
		BillingFlags:  BillingNone,
		BillingRested: 0,
		Expansion:     realmd.ExpansionWrath,
	}

	if err := client.SendPacket(realmd.OpServerAuthResponse, &resp); err != nil {
		return err
	}

	log.Println("sent auth response")

	return nil
}

// Returns whether the client is authenticated. Clients authenticate with authd and then send a proof
// to realmd to show that they know the session key. We verify the proof by fetching the session key
// from the DB, calculating the proof ourselves, and checking that it matches. If it does, the client
// is considered authenticated.
func authenticateClient(svc *realmd.Service, client *realmd.Client, p *proofRequest) (bool, error) {
	acct, err := svc.Accounts.Get(&model.AccountGetParams{Username: p.Username})
	if err != nil {
		return false, err
	} else if acct == nil {
		log.Printf("authenticateClient: no account with username %s exists", p.Username)
		return false, nil
	}

	realm, err := svc.Realms.Get(p.RealmId)
	if err != nil {
		return false, err
	} else if realm == nil {
		log.Printf("authenticateClient: no realm with id %d exists", p.RealmId)
		return false, nil
	}

	session, err := svc.Sessions.Get(acct.Id)
	if err != nil {
		return false, err
	} else if session == nil {
		log.Printf("authenticateClient: no session for username %s exists", acct.Username)
		return false, nil
	}

	if err := session.Decode(); err != nil {
		return false, err
	}

	proof := srp.CalculateWorldProof(p.Username, p.ClientSeed[:], client.ServerSeed, session.SessionKey())

	// Client sent correct proof?
	if !bytes.Equal(proof, p.ClientProof[:]) {
		log.Println("authenticateClient: invalid proof")
		return false, nil
	}

	client.Account = acct
	client.Realm = realm
	client.Session = session

	log.Println("client authenticated successfully")

	return true, nil
}
