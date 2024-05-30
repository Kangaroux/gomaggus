package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	srpv2 "github.com/kangaroux/go-realmd/srp_v2"
)

const (
	OP_LOGIN_CHALLENGE     byte = 0
	OP_LOGIN_PROOF         byte = 1
	OP_RECONNECT_CHALLENGE byte = 2
	OP_RECONNECT_PROOF     byte = 3
	OP_REALM_LIST          byte = 16

	WOW_SUCCESS              byte = 0
	WOW_FAIL_UNKNOWN_ACCOUNT byte = 4

	MOCK_USERNAME = "TEST"
	MOCK_PASSWORD = "PASSWORD"
)

type RealmType uint8

const (
	REALMTYPE_PVE   RealmType = 0
	REALMTYPE_PVP   RealmType = 1
	REALMTYPE_RP    RealmType = 6
	REALMTYPE_RPPVP RealmType = 8
)

type RealmFlag uint8

const (
	REALMFLAG_NONE          RealmFlag = 0
	REALMFLAG_INVALID       RealmFlag = 1 // Realm is greyed out and can't be selected
	REALMFLAG_OFFLINE       RealmFlag = 2 // Population: "Offline" and can't be selected
	REALMFLAG_SPECIFY_BUILD RealmFlag = 4 // Includes version in realm name
	REALMFLAG_UNKNOWN1      RealmFlag = 8
	REALMFLAG_UNKNOWN2      RealmFlag = 16
	REALMFLAG_NEW_PLAYERS   RealmFlag = 32  // Population: "New Players" in blue text
	REALMFLAG_NEW_SERVER    RealmFlag = 64  // Population: "New" in green text
	REALMFLAG_FULL          RealmFlag = 128 // Population: "Full" in red text
)

type RealmRegion uint8

const (
	REALMREGION_DEV           RealmRegion = 1
	REALMREGION_US            RealmRegion = 2
	REALMREGION_OCEANIC       RealmRegion = 3
	REALMREGION_LATIN_AMERICA RealmRegion = 4
	REALMREGION_TOURNAMENT    RealmRegion = 5
	REALMREGION_KOREA         RealmRegion = 6
	REALMREGION_TOURNAMENT2   RealmRegion = 7
	REALMREGION_ENGLISH       RealmRegion = 8
	REALMREGION_GERMAN        RealmRegion = 9
	REALMREGION_FRENCH        RealmRegion = 10
	REALMREGION_SPANISH       RealmRegion = 11
	REALMREGION_RUSSIAN       RealmRegion = 12
	REALMREGION_TOURNAMENT3   RealmRegion = 13
	REALMREGION_TAIWAN        RealmRegion = 14
	REALMREGION_TOURNAMENT4   RealmRegion = 15
	REALMREGION_CHINA         RealmRegion = 16
	REALMREGION_CN1           RealmRegion = 17
	REALMREGION_CN2           RealmRegion = 18
	REALMREGION_CN3           RealmRegion = 19
	REALMREGION_CN4           RealmRegion = 20
	REALMREGION_CN5           RealmRegion = 21
	REALMREGION_CN6           RealmRegion = 22
	REALMREGION_CN7           RealmRegion = 23
	REALMREGION_CN8           RealmRegion = 24
	REALMREGION_TOURNAMENT5   RealmRegion = 25
	REALMREGION_TEST          RealmRegion = 26
	REALMREGION_TOURNAMENT6   RealmRegion = 27
	REALMREGION_QA            RealmRegion = 28
	REALMREGION_CN9           RealmRegion = 29
	REALMREGION_TEST2         RealmRegion = 30
	REALMREGION_CN10          RealmRegion = 31
	REALMREGION_CTC           RealmRegion = 32
	REALMREGION_CNC           RealmRegion = 33
	REALMREGION_CN1_4         RealmRegion = 34
	REALMREGION_CN2_6_9       RealmRegion = 35
	REALMREGION_CN3_7         RealmRegion = 36
	REALMREGION_CN5_8         RealmRegion = 37
)

var (
	MOCK_SALT        []byte
	MOCK_VERIFIER    []byte
	MOCK_PRIVATE_KEY []byte
	MOCK_PUBLIC_KEY  []byte

	MOCK_REALMS = []Realm{
		{
			Type:            REALMTYPE_PVE,
			Locked:          false,
			Flags:           REALMFLAG_NONE,
			Name:            "Test Realm\x00",
			Host:            "localhost:8085\x00",
			Population:      0.01,
			NumCharsOnRealm: 0,
			Region:          REALMREGION_US,
			Id:              0,
			// Version:         RealmVersion{Major: 4, Minor: 3, Patch: 6, Build: 12340},
		},
		// {
		// 	Type:            REALMTYPE_PVP,
		// 	Locked:          false,
		// 	Flags:           REALMFLAG_NONE,
		// 	Name:            "Test Realm1\x00",
		// 	Host:            "localhost:8085\x00",
		// 	Population:      0,
		// 	NumCharsOnRealm: 0,
		// 	Region:          REALMREGION_US,
		// 	Id:              1,
		// 	// Version:         RealmVersion{Major: 4, Minor: 3, Patch: 6, Build: 12340},
		// },
	}
)

var (
	// Maps usernames to session keys. This tracks the last known session key for a particular
	// user and is used for reconnecting.
	SESSION_KEY_MAP = make(map[string][]byte)
)

func init() {
	MOCK_SALT = make([]byte, 32)
	if _, err := rand.Read(MOCK_SALT); err != nil {
		log.Fatalf("error generating salt: %v\n", err)
	}

	MOCK_VERIFIER = srpv2.CalculateVerifier(MOCK_USERNAME, MOCK_PASSWORD, MOCK_SALT)
	MOCK_PRIVATE_KEY = srpv2.NewPrivateKey()
	MOCK_PUBLIC_KEY = srpv2.CalculateServerPublicKey(MOCK_VERIFIER, MOCK_PRIVATE_KEY)
}

func main() {
	listener, err := net.Listen("tcp", ":3724")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Print("listening on port 3724")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleClient(conn)
	}
}

type Client struct {
	conn          net.Conn
	username      string
	reconnectData []byte
	sessionKey    []byte
}

func handleClient(c net.Conn) {
	buf := make([]byte, 4096)
	client := &Client{conn: c, reconnectData: make([]byte, 16)}

	log.Printf("client connected from %v\n", c.RemoteAddr().String())

	for {
		n, err := c.Read(buf)

		if err == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		} else if err != nil {
			log.Printf("error reading from client: %v\n", err)
			c.Close()
			return
		}

		log.Printf("read %d bytes\n", n)

		if err := handlePacket(client, buf[:n]); err != nil {
			log.Printf("error handling packet: %v\n", err)
			c.Close()
			return
		}
	}
}

// https://gtker.com/wow_messages/docs/cmd_auth_logon_challenge_client.html
type LoginChallengePacket struct {
	Opcode         byte // 0x0 or 0x2 if reconnecting
	Error          byte // unused?
	Size           uint16
	GameName       [4]byte
	Version        [3]byte
	Build          uint16
	OSArch         [4]byte
	OS             [4]byte
	Locale         [4]byte
	TimezoneBias   uint32
	IP             [4]byte
	UsernameLength uint8

	// The username is a variable size and needs to be read manually
	// Username string
}

// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_client.html#protocol-version-8
type LoginProofPacket struct {
	Opcode           byte // 0x1
	ClientPublicKey  [32]byte
	ClientProof      [20]byte
	CRCHash          [20]byte // unused
	NumTelemetryKeys uint8    // unused
}

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_proof_client.html
type ReconnectProofPacket struct {
	Opcode         byte // 0x3
	ProofData      [16]byte
	ClientProof    [20]byte
	ClientChecksum [20]byte // unused
	KeyCount       byte     // unused
}

type RealmVersion struct {
	Major uint8
	Minor uint8
	Patch uint8
	Build uint16
}

// https://gtker.com/wow_messages/docs/realm.html#protocol-version-8
type Realm struct {
	Type   RealmType
	Locked bool
	Flags  RealmFlag
	Name   string // C-style NUL terminated, e.g. "Test Realm\x00"
	Host   string // C-style NUL terminated, e.g. "localhost:8085\x00"

	// A percentage of how full the server is with active sessions. Mangos has the upper limit of this
	// value as 2.0 for some reason. The game client only seems to interpret this value on an absolute
	// scale if there is only one realm. It seems like when there are multiple realms, it compares the
	// pop relatively, i.e. whatever realm has the highest pop is now the upper limit. Suffice to say,
	// it's not important whether this value is accurate.
	Population      float32
	NumCharsOnRealm uint8 // Number of characters for the logged in account
	Region          RealmRegion
	Id              uint8
	Version         RealmVersion // included only if REALMFLAG_SPECIFY_BUILD flag is set
}

func handlePacket(c *Client, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("error: packet is empty")
	}

	switch data[0] {
	case OP_LOGIN_CHALLENGE:
		log.Println("Starting login challenge")
		p := LoginChallengePacket{}
		reader := bytes.NewReader(data)
		if err := binary.Read(reader, binary.LittleEndian, &p); err != nil {
			return err
		}
		usernameBytes := make([]byte, p.UsernameLength)
		if _, err := reader.Read(usernameBytes); err != nil {
			return err
		}
		c.username = strings.ToUpper(string(usernameBytes))
		log.Printf("client trying to login as '%s'\n", c.username)

		// https://gtker.com/wow_messages/docs/cmd_auth_logon_challenge_server.html#protocol-version-8
		resp := &bytes.Buffer{}
		resp.WriteByte(OP_LOGIN_CHALLENGE)
		resp.WriteByte(0) // protocol version

		if c.username == MOCK_USERNAME {
			resp.WriteByte(WOW_SUCCESS)
			resp.Write(MOCK_PUBLIC_KEY)
			resp.WriteByte(1)  // generator size (1 byte)
			resp.WriteByte(7)  // generator
			resp.WriteByte(32) // large prime size (32 bytes)
			resp.Write(srpv2.Reverse(srpv2.LargeSafePrime))
			resp.Write(MOCK_SALT)
			resp.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) // crc hash
			resp.WriteByte(0)
		} else {
			resp.WriteByte(WOW_FAIL_UNKNOWN_ACCOUNT)
		}

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("Replied to login challenge")
		return nil
	case OP_LOGIN_PROOF:
		log.Println("Starting login proof")
		p := LoginProofPacket{}
		reader := bytes.NewReader(data)
		if err := binary.Read(reader, binary.LittleEndian, &p); err != nil {
			return err
		}

		clientPublicKey := p.ClientPublicKey[:]
		clientProof := p.ClientProof[:]

		c.sessionKey = srpv2.CalculateServerSessionKey(
			clientPublicKey, MOCK_PUBLIC_KEY, MOCK_PRIVATE_KEY, MOCK_VERIFIER)
		calculatedClientProof := srpv2.CalculateClientProof(
			MOCK_USERNAME, MOCK_SALT, clientPublicKey, MOCK_PUBLIC_KEY, c.sessionKey,
		)

		// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_server.html#protocol-version-8
		resp := &bytes.Buffer{}
		resp.WriteByte(OP_LOGIN_PROOF)

		if !bytes.Equal(calculatedClientProof, clientProof) {
			resp.WriteByte(WOW_FAIL_UNKNOWN_ACCOUNT)
			resp.Write([]byte{0, 0}) // padding
		} else {
			resp.WriteByte(WOW_SUCCESS)
			resp.Write(srpv2.CalculateServerProof(clientPublicKey, clientProof, c.sessionKey))
			resp.Write([]byte{0, 0, 0, 0}) // Account flag
			resp.Write([]byte{0, 0, 0, 0}) // Hardware survey ID
			resp.Write([]byte{0, 0})       // Unknown

			// Save the session key in case the client needs to reconnect later
			SESSION_KEY_MAP[c.username] = c.sessionKey
		}

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("Replied to login proof")
		return nil
	case OP_RECONNECT_CHALLENGE:
		log.Println("Starting reconnect challenge")
		p := LoginChallengePacket{}
		reader := bytes.NewReader(data)
		if err := binary.Read(reader, binary.LittleEndian, &p); err != nil {
			return err
		}
		usernameBytes := make([]byte, p.UsernameLength)
		if _, err := reader.Read(usernameBytes); err != nil {
			return err
		}
		c.username = strings.ToUpper(string(usernameBytes))
		log.Printf("client trying to login as '%s'\n", c.username)

		// Generate random data that will be used for the reconnect proof
		if _, err := rand.Read(c.reconnectData); err != nil {
			return err
		}

		sessionKey, hasSessionKey := SESSION_KEY_MAP[c.username]

		// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_challenge_server.html#protocol-version-8
		resp := &bytes.Buffer{}
		resp.WriteByte(OP_RECONNECT_CHALLENGE)

		if c.username == MOCK_USERNAME && hasSessionKey {
			resp.WriteByte(WOW_SUCCESS)
			resp.Write(c.reconnectData)
			resp.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) // checksum salt

			c.sessionKey = sessionKey
		} else {
			resp.WriteByte(WOW_FAIL_UNKNOWN_ACCOUNT)
		}

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("Replied to reconnect challenge")
		return nil
	case OP_RECONNECT_PROOF:
		log.Println("Starting reconnect proof")
		p := ReconnectProofPacket{}
		reader := bytes.NewReader(data)
		if err := binary.Read(reader, binary.LittleEndian, &p); err != nil {
			return err
		}

		serverProof := srpv2.CalculateReconnectProof(c.username, p.ProofData[:], c.reconnectData, c.sessionKey)

		log.Printf("computed recon proof: %x\n", serverProof)
		log.Printf("client proof: %x\n", p.ClientProof)

		// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_server.html#protocol-version-8
		resp := &bytes.Buffer{}
		resp.WriteByte(OP_RECONNECT_PROOF)
		if !bytes.Equal(serverProof, p.ClientProof[:]) {
			resp.WriteByte(WOW_FAIL_UNKNOWN_ACCOUNT)
			resp.Write([]byte{0, 0}) // padding
		} else {
			resp.WriteByte(WOW_SUCCESS)
			resp.Write(serverProof)
			resp.Write([]byte{0, 0, 0, 0}) // Account flag
			resp.Write([]byte{0, 0, 0, 0}) // Hardware survey ID
			resp.Write([]byte{0, 0})       // Unknown
		}

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("Replied to reconnect proof")
		return nil
	case OP_REALM_LIST:
		// https://gtker.com/wow_messages/docs/cmd_realm_list_server.html#protocol-version-8
		resp := &bytes.Buffer{}
		resp.WriteByte(OP_REALM_LIST)

		realmList := &bytes.Buffer{}
		realmList.Write([]byte{0, 0, 0, 0}) // header padding
		binary.Write(realmList, binary.LittleEndian, uint16(len(MOCK_REALMS)))
		for _, r := range MOCK_REALMS {
			realmList.WriteByte(byte(r.Type))

			if r.Locked {
				realmList.WriteByte(1)
			} else {
				realmList.WriteByte(0)
			}

			realmList.WriteByte(byte(r.Flags))
			realmList.WriteString(r.Name)
			realmList.WriteString(r.Host)
			binary.Write(realmList, binary.LittleEndian, r.Population)
			realmList.WriteByte(byte(r.NumCharsOnRealm))
			realmList.WriteByte(byte(r.Region))
			realmList.WriteByte(byte(r.Id))

			if r.Flags&REALMFLAG_SPECIFY_BUILD > 0 {
				realmList.WriteByte(byte(r.Version.Major))
				realmList.WriteByte(byte(r.Version.Minor))
				realmList.WriteByte(byte(r.Version.Patch))
				binary.Write(realmList, binary.LittleEndian, r.Version.Build)
			}
		}
		realmList.Write([]byte{0, 0}) // footer padding

		// Write size of realm list payload
		binary.Write(resp, binary.LittleEndian, uint16(realmList.Len()))
		// Concat to main payload
		realmList.WriteTo(resp)

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("Replied to realm list")
		return nil
	default:
		return fmt.Errorf("error: unknown opcode (%v)", data[0])
	}
}
