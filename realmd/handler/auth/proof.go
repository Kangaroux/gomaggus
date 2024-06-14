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

func ProofHandler(svc *realmd.Service, client *realmd.Client, data []byte) error {
	log.Println("starting auth session")

	var err error
	r := bytes.NewReader(data[6:])

	p := proofRequest{}
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

	client.Authenticated, err = authenticateClient(svc, client, &p)
	if err != nil {
		return err
	}

	if !client.Authenticated {
		// We can't return an error to the client due to the header encryption, just drop the connection
		return errors.New("client could not be authenticated")
	}

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
	var err error

	client.Account, err = svc.Accounts.Get(&model.AccountGetParams{Username: p.Username})
	if err != nil {
		return false, err
	} else if client.Account == nil {
		log.Printf("authenticateClient: no account with username %s exists", p.Username)
		return false, nil
	}

	if client.Realm, err = svc.Realms.Get(p.RealmId); err != nil {
		return false, err
	} else if client.Realm == nil {
		log.Printf("authenticateClient: no realm with id %d exists", p.RealmId)
		return false, nil
	}

	if client.Session, err = svc.Sessions.Get(client.Account.Id); err != nil {
		return false, err
	} else if client.Session == nil {
		log.Printf("authenticateClient: no session for username %s exists", client.Account.Username)
		return false, nil
	}

	if err := client.Session.Decode(); err != nil {
		return false, err
	}

	// Initialize the header encryption using the session key
	client.Crypto = realmd.NewHeaderCrypto(client.Session.SessionKey())
	if err := client.Crypto.Init(); err != nil {
		return false, err
	}

	proof := srp.CalculateWorldProof(p.Username, p.ClientSeed[:], client.ServerSeed, client.Session.SessionKey())

	// Check to see if the proof sent by the client is correct
	if !bytes.Equal(proof, p.ClientProof[:]) {
		log.Println("authenticateClient: invalid proof")
		return false, nil
	}

	log.Println("client authenticated successfully")

	return true, nil
}
