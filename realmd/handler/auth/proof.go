package auth

import (
	"bytes"

	srp "github.com/kangaroux/go-wow-srp6"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
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
		return &realmd.ErrKickClient{Reason: "auth failed"}
	}

	client.Authenticated = true

	if err := client.HeaderCrypto.Init(client.Session.SessionKey()); err != nil {
		return err
	}

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
		client.Log.Warn().Str("username", p.Username).Msg("username does not exist")
		return false, nil
	}

	realm, err := svc.Realms.Get(p.RealmId)
	if err != nil {
		return false, err
	} else if realm == nil {
		client.Log.Warn().Uint32("realm", p.RealmId).Msg("realm does not exist")
		return false, nil
	}

	// Fetch the last known session from the DB. This is set by authd when the client logs in.
	session, err := svc.Sessions.Get(acct.Id)
	if err != nil {
		return false, err
	} else if session == nil {
		client.Log.Warn().Str("username", p.Username).Msg("session does not exist")
		return false, nil
	}

	if err := session.Decode(); err != nil {
		return false, err
	}

	proof := srp.WorldProof(p.Username, p.ClientSeed[:], client.ServerSeed, session.SessionKey())

	// Did the client authenticate with authd first?
	if !bytes.Equal(proof, p.ClientProof[:]) {
		client.Log.Warn().Msg("invalid proof")
		return false, nil
	}

	client.Account = acct
	client.Realm = realm
	client.Session = session

	client.Log.Info().Str("realm", client.Realm.String()).Msg("client authenticated")

	return true, nil
}
