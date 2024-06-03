package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	mathrand "math/rand"
	"net"
	"strings"

	"github.com/kangaroux/gomaggus/internal/worldd"
)

// Opcodes sent by the server
const (
	OP_AUTH_CHALLENGE uint16 = 0x1EC
	OP_AUTH_RESPONSE  uint16 = 0x1EE
)

// Opcodes sent by the client
const (
	OP_AUTH_SESSION uint32 = 0x1ED
)

const (
	RespCodeResponseSuccess                                byte = 0x00
	RespCodeResponseFailure                                byte = 0x01
	RespCodeResponseCancelled                              byte = 0x02
	RespCodeResponseDisconnected                           byte = 0x03
	RespCodeResponseFailedToConnect                        byte = 0x04
	RespCodeResponseConnected                              byte = 0x05
	RespCodeResponseVersionMismatch                        byte = 0x06
	RespCodeCStatusConnecting                              byte = 0x07
	RespCodeCStatusNegotiatingSecurity                     byte = 0x08
	RespCodeCStatusNegotiationComplete                     byte = 0x09
	RespCodeCStatusNegotiationFailed                       byte = 0x0A
	RespCodeCStatusAuthenticating                          byte = 0x0B
	RespCodeAuthOk                                         byte = 0x0C
	RespCodeAuthFailed                                     byte = 0x0D
	RespCodeAuthReject                                     byte = 0x0E
	RespCodeAuthBadServerProof                             byte = 0x0F
	RespCodeAuthUnavailable                                byte = 0x10
	RespCodeAuthSystemError                                byte = 0x11
	RespCodeAuthBillingError                               byte = 0x12
	RespCodeAuthBillingExpired                             byte = 0x13
	RespCodeAuthVersionMismatch                            byte = 0x14
	RespCodeAuthUnknownAccount                             byte = 0x15
	RespCodeAuthIncorrectPassword                          byte = 0x16
	RespCodeAuthSessionExpired                             byte = 0x17
	RespCodeAuthServerShuttingDown                         byte = 0x18
	RespCodeAuthAlreadyLoggingIn                           byte = 0x19
	RespCodeAuthLoginServerNotFound                        byte = 0x1A
	RespCodeAuthWaitQueue                                  byte = 0x1B
	RespCodeAuthBanned                                     byte = 0x1C
	RespCodeAuthAlreadyOnline                              byte = 0x1D
	RespCodeAuthNoTime                                     byte = 0x1E
	RespCodeAuthDbBusy                                     byte = 0x1F
	RespCodeAuthSuspended                                  byte = 0x20
	RespCodeAuthParentalControl                            byte = 0x21
	RespCodeAuthLockedEnforced                             byte = 0x22
	RespCodeRealmListInProgress                            byte = 0x23
	RespCodeRealmListSuccess                               byte = 0x24
	RespCodeRealmListFailed                                byte = 0x25
	RespCodeRealmListInvalid                               byte = 0x26
	RespCodeRealmListRealmNotFound                         byte = 0x27
	RespCodeAccountCreateInProgress                        byte = 0x28
	RespCodeAccountCreateSuccess                           byte = 0x29
	RespCodeAccountCreateFailed                            byte = 0x2A
	RespCodeCharListRetrieving                             byte = 0x2B
	RespCodeCharListRetrieved                              byte = 0x2C
	RespCodeCharListFailed                                 byte = 0x2D
	RespCodeCharCreateInProgress                           byte = 0x2E
	RespCodeCharCreateSuccess                              byte = 0x2F
	RespCodeCharCreateError                                byte = 0x30
	RespCodeCharCreateFailed                               byte = 0x31
	RespCodeCharCreateNameInUse                            byte = 0x32
	RespCodeCharCreateDisabled                             byte = 0x33
	RespCodeCharCreatePvpTeamsViolation                    byte = 0x34
	RespCodeCharCreateServerLimit                          byte = 0x35
	RespCodeCharCreateAccountLimit                         byte = 0x36
	RespCodeCharCreateServerQueue                          byte = 0x37
	RespCodeCharCreateOnlyExisting                         byte = 0x38
	RespCodeCharCreateExpansion                            byte = 0x39
	RespCodeCharCreateExpansionClass                       byte = 0x3A
	RespCodeCharCreateLevelRequirement                     byte = 0x3B
	RespCodeCharCreateUniqueClassLimit                     byte = 0x3C
	RespCodeCharCreateCharacterInGuild                     byte = 0x3D
	RespCodeCharCreateRestrictedRaceclass                  byte = 0x3E
	RespCodeCharCreateCharacterChooseRace                  byte = 0x3F
	RespCodeCharCreateCharacterArenaLeader                 byte = 0x40
	RespCodeCharCreateCharacterDeleteMail                  byte = 0x41
	RespCodeCharCreateCharacterSwapFaction                 byte = 0x42
	RespCodeCharCreateCharacterRaceOnly                    byte = 0x43
	RespCodeCharCreateCharacterGoldLimit                   byte = 0x44
	RespCodeCharCreateForceLogin                           byte = 0x45
	RespCodeCharDeleteInProgress                           byte = 0x46
	RespCodeCharDeleteSuccess                              byte = 0x47
	RespCodeCharDeleteFailed                               byte = 0x48
	RespCodeCharDeleteFailedLockedForTransfer              byte = 0x49
	RespCodeCharDeleteFailedGuildLeader                    byte = 0x4A
	RespCodeCharDeleteFailedArenaCaptain                   byte = 0x4B
	RespCodeCharLoginInProgress                            byte = 0x4C
	RespCodeCharLoginSuccess                               byte = 0x4D
	RespCodeCharLoginNoWorld                               byte = 0x4E
	RespCodeCharLoginDuplicateCharacter                    byte = 0x4F
	RespCodeCharLoginNoInstances                           byte = 0x50
	RespCodeCharLoginFailed                                byte = 0x51
	RespCodeCharLoginDisabled                              byte = 0x52
	RespCodeCharLoginNoCharacter                           byte = 0x53
	RespCodeCharLoginLockedForTransfer                     byte = 0x54
	RespCodeCharLoginLockedByBilling                       byte = 0x55
	RespCodeCharLoginLockedByMobileAh                      byte = 0x56
	RespCodeCharNameSuccess                                byte = 0x57
	RespCodeCharNameFailure                                byte = 0x58
	RespCodeCharNameNoName                                 byte = 0x59
	RespCodeCharNameTooShort                               byte = 0x5A
	RespCodeCharNameTooLong                                byte = 0x5B
	RespCodeCharNameInvalidCharacter                       byte = 0x5C
	RespCodeCharNameMixedLanguages                         byte = 0x5D
	RespCodeCharNameProfane                                byte = 0x5E
	RespCodeCharNameReserved                               byte = 0x5F
	RespCodeCharNameInvalidApostrophe                      byte = 0x60
	RespCodeCharNameMultipleApostrophes                    byte = 0x61
	RespCodeCharNameThreeConsecutive                       byte = 0x62
	RespCodeCharNameInvalidSpace                           byte = 0x63
	RespCodeCharNameConsecutiveSpaces                      byte = 0x64
	RespCodeCharNameRussianConsecutiveSilentCharacters     byte = 0x65
	RespCodeCharNameRussianSilentCharacterAtBeginningOrEnd byte = 0x66
	RespCodeCharNameDeclensionDoesntMatchBaseName          byte = 0x67
)

const (
	ExpansionVanilla byte = 0x0
	ExpansionTbc     byte = 0x1
	ExpansionWrath   byte = 0x2
)

func main() {
	listener, err := net.Listen("tcp", ":8085")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Print("listening on port 8085")

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
	serverSeed    uint32
	authenticated bool
	crypto        *worldd.WrathHeaderCrypto
}

func handleClient(c net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recovered from panic: %v", err)

			if err := c.Close(); err != nil {
				log.Printf("error closing after recover: %v", err)
			}
		}
	}()

	log.Printf("client connected from %v\n", c.RemoteAddr().String())

	buf := make([]byte, 4096)
	client := &Client{
		conn:       c,
		serverSeed: mathrand.Uint32(),
		crypto:     worldd.NewWrathHeaderCrypto(nil /* TODO session key */),
	}

	// The server is the one who initiates the auth challenge here, unlike the login server where
	// the client is the one who initiates it
	if err := sendAuthChallenge(client); err != nil {
		log.Printf("error sending auth challenge: %v\n", err)
		c.Close()
		return
	}

	for {
		log.Println("waiting to read...")
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

// https://gtker.com/wow_messages/docs/smsg_auth_challenge.html#client-version-335
func sendAuthChallenge(c *Client) error {
	body := &bytes.Buffer{}
	body.Write([]byte{1, 0, 0, 0}) // unknown
	binary.Write(body, binary.LittleEndian, c.serverSeed)

	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return err
	}
	body.Write(seed) // seed, unused. This differs from the 4 byte server seed

	resp := &bytes.Buffer{}
	respHeader, err := makeServerHeader(OP_AUTH_CHALLENGE, uint32(body.Len()))
	if err != nil {
		return err
	}
	resp.Write(respHeader)
	resp.Write(body.Bytes())

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent auth challenge")
	return nil
}

type Header struct {
	Size   uint16
	Opcode uint32
}

func parseHeader(c *Client, data []byte) (*Header, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("parseHeader: payload should be at least 6 bytes but it's only %d", len(data))
	}

	headerData := data[:6]

	if c.authenticated {
		if c.crypto == nil {
			return nil, errors.New("parseHeader: client is authenticated but client.crypto is nil")
		}

		headerData = c.crypto.Decrypt(headerData)
	}

	h := &Header{
		Size:   binary.BigEndian.Uint16(headerData[:2]),
		Opcode: binary.LittleEndian.Uint32(headerData[2:6]),
	}

	return h, nil
}

type AuthSessionPacket struct {
	ClientBuild     uint32
	LoginServerId   uint32
	Username        string
	LoginServerType uint32
	ClientSeed      uint32
	RegionId        uint32
	BattlegroundId  uint32
	RealmId         uint32
	DOSResponse     uint64
	ClientProof     [20]byte
	AddonInfo       []byte
}

func readCString(r *bytes.Reader) (string, error) {
	s := strings.Builder{}

	for {
		b, err := r.ReadByte()

		if err != nil {
			return "", err
		} else if b == 0x0 {
			break
		}

		s.WriteByte(b)
	}

	return s.String(), nil
}

func handlePacket(c *Client, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("handlePacket: packet is empty")
	}

	header, err := parseHeader(c, data)
	if err != nil {
		return err
	}

	switch header.Opcode {
	case OP_AUTH_SESSION:
		log.Println("starting auth session")

		r := bytes.NewReader(data)

		// Skip the header
		r.Seek(6, io.SeekStart)

		p := AuthSessionPacket{}
		if err = binary.Read(r, binary.BigEndian, &p.ClientBuild); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.LoginServerId); err != nil {
			return err
		}
		if p.Username, err = readCString(r); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.LoginServerType); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.ClientSeed); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.RegionId); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.BattlegroundId); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.RealmId); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.DOSResponse); err != nil {
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

		// TODO: Check client proof

		// https://gtker.com/wow_messages/docs/smsg_auth_response.html#client-version-335
		inner := bytes.Buffer{}
		inner.WriteByte(RespCodeAuthOk)
		inner.Write([]byte{0, 0, 0, 0})   // billing time
		inner.WriteByte(0x0)              // billing flags
		inner.Write([]byte{0, 0, 0, 0})   // billing rested
		inner.WriteByte(ExpansionVanilla) // exp
		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OP_AUTH_RESPONSE, uint32(inner.Len()))
		if err != nil {
			return err
		}
		resp.Write(c.crypto.Encrypt(respHeader))
		resp.Write(inner.Bytes())

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent auth response")

		return nil
	}

	return nil
}

const (
	// 23 bits + 1 bit for LARGE_HEADER_FLAG
	SIZE_FIELD_MAX_VALUE = 0x7FFFFF

	// 15 bits (16th bit is reserved for LARGE_HEADER_FLAG)
	LARGE_HEADER_THRESHOLD = 0x7FFF

	// Set on MSB of size field (first header byte)
	LARGE_HEADER_FLAG = 0x80
)

func makeServerHeader(opcode uint16, size uint32) ([]byte, error) {
	// Include the opcode in the size
	size += 2

	if size > SIZE_FIELD_MAX_VALUE {
		return nil, fmt.Errorf("makeServerHeader: size is too large (%d bytes)", size)
	}

	var header []byte

	// The size field in the header can be 2 or 3 bytes. The most significant bit in the size field
	// is reserved as a flag to indicate this. In total, server headers are 4 or 5 bytes.
	//
	// The header format is: <size><opcode>
	// <size> is 2-3 bytes big endian
	// <opcode> is 2 bytes little endian
	if size > LARGE_HEADER_THRESHOLD {
		header = []byte{
			byte(size>>16) | LARGE_HEADER_FLAG,
			byte(size >> 8),
			byte(size),
			byte(opcode),
			byte(opcode >> 8),
		}
	} else {
		header = []byte{
			byte(size >> 8),
			byte(size),
			byte(opcode),
			byte(opcode >> 8),
		}
	}

	return header, nil
}
