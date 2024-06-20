package auth

import (
	"bytes"
	"errors"
	"log"

	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/kangaroux/gomaggus/srp"
	"github.com/mixcode/binarystruct"
)

// https://gtker.com/wow_messages/docs/billingplanflags.html
type billingFlag uint8

const (
	billingNone          billingFlag = 0x00
	billingUnused        billingFlag = 0x01
	billingRecurringBill billingFlag = 0x02
	billingFreeTrial     billingFlag = 0x04
	billingIgr           billingFlag = 0x08
	billingUsage         billingFlag = 0x10
	billingTimeMixture   billingFlag = 0x20
	billingRestricted    billingFlag = 0x40
	billingEnableCais    billingFlag = 0x80
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
	BillingFlags  billingFlag
	BillingRested uint32
	Expansion     realmd.Expansion
}

// TODO
// type proofWaitQueue struct {
// 	ResponseCode         realmd.ResponseCode
// 	HasFreeCharMigration bool
// }

func ProofHandler(svc *realmd.Service, client *realmd.Client, data []byte) error {
	req := proofRequest{}
	if _, err := binarystruct.Unmarshal(data, binarystruct.LittleEndian, &req); err != nil {
		return err
	}

	authenticated, err := authenticateClient(svc, client, &req)
	if err != nil {
		return err
	}

	if !authenticated {
		// The client expects the proof response to be successful (since it just authenticated with authd).
		// If the authentication failed, no response is returned and the connection is closed.
		return errors.New("client could not be authenticated")
	}

	client.Authenticated = true

	headerCrypto := realmd.NewHeaderCrypto(client.Session.SessionKey())
	if err := headerCrypto.Init(); err != nil {
		return err
	}
	client.HeaderCrypto = headerCrypto

	resp := proofSuccess{
		ResponseCode:  realmd.RespCodeAuthOk,
		BillingTime:   0,
		BillingFlags:  billingNone,
		BillingRested: 0,
		Expansion:     realmd.ExpansionWrath,
	}
	return client.SendPacket(realmd.OpServerAuthResponse, &resp)
}

// authenticateClient reports whether the client's proof is valid.
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

	// Fetch the last known session from the DB. This is set by authd when the client logs in.
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

	// Did the client authenticate with authd first?
	if !bytes.Equal(proof, p.ClientProof[:]) {
		log.Println("authenticateClient: invalid proof")
		return false, nil
	}

	client.Account = acct
	client.Realm = realm
	client.Session = session

	log.Println("Authenticated", client.Account, "on", client.Realm)

	return true, nil
}
