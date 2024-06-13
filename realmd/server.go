package realmd

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"net"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/models"
)

const (
	DefaultListenAddr = ":8085"
)

type Server struct {
	listenAddr string

	accountsDb models.AccountService
	charsDb    models.CharacterService
	realmsDb   models.RealmService
	sessionsDb models.SessionService
}

func NewServer(db *sqlx.DB, listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		accountsDb: models.NewDbAccountService(db),
		charsDb:    models.NewDbCharacterervice(db),
		realmsDb:   models.NewDbRealmService(db),
		sessionsDb: models.NewDbSessionService(db),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp4", s.listenAddr)

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Printf("listening on %s\n", listener.Addr().String())

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recovered from panic: %v", err)

			if err := conn.Close(); err != nil {
				log.Printf("error closing after recover: %v", err)
			}
		}
	}()

	log.Printf("client connected from %v\n", conn.RemoteAddr().String())

	client := &Client{conn: conn}
	binary.BigEndian.PutUint32(client.serverSeed[:], mrand.Uint32())

	// The server is the one who initiates the auth challenge here, unlike the login server where
	// the client is the one who initiates it
	if err := s.sendAuthChallenge(client); err != nil {
		log.Printf("error sending auth challenge: %v\n", err)
		conn.Close()
		return
	}

	buf := make([]byte, 4096)

	for {
		log.Println("waiting to read...")
		n, err := conn.Read(buf)

		if err == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		} else if err != nil {
			log.Printf("error reading from client: %v\n", err)
			conn.Close()
			return
		}

		log.Printf("read %d bytes\n", n)

		if err := s.handlePacket(client, buf[:n]); err != nil {
			log.Printf("error handling packet: %v\n", err)
			conn.Close()
			return
		}
	}
}

// https://gtker.com/wow_messages/docs/smsg_auth_challenge.html#client-version-335
func (s *Server) sendAuthChallenge(c *Client) error {
	body := &bytes.Buffer{}
	body.Write([]byte{1, 0, 0, 0}) // unknown
	binary.Write(body, binary.BigEndian, c.serverSeed)

	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return err
	}
	body.Write(seed) // seed, unused. This differs from the 4 byte server seed

	resp := &bytes.Buffer{}
	respHeader, err := makeServerHeader(OpServerAuthChallenge, uint32(body.Len()))
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
		Opcode: ClientOpcode(binary.LittleEndian.Uint32(headerData[2:6])),
	}

	return h, nil
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

func (s *Server) handlePacket(c *Client, data []byte) error {
	var err error

	if len(data) == 0 {
		return fmt.Errorf("handlePacket: packet is empty")
	}

	header, err := parseHeader(c, data)
	if err != nil {
		return err
	}

	switch header.Opcode {
	case OpClientAuthSession:
		log.Println("starting auth session")

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

		c.authenticated, err = s.authenticateClient(c, &p)
		if err != nil {
			return err
		}

		if !c.authenticated {
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
		resp.Write(c.crypto.Encrypt(respHeader))
		resp.Write(inner.Bytes())

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent auth response")

		return nil
	case OpClientPing:
		log.Println("starting ping")

		r := bytes.NewReader(data[6:])
		p := PingPacket{}
		if err = binary.Read(r, binary.LittleEndian, &p.SequenceId); err != nil {
			return err
		}
		if err = binary.Read(r, binary.LittleEndian, &p.RoundTripTime); err != nil {
			return err
		}

		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OpServerPong, 4)
		if err != nil {
			return err
		}
		resp.Write(c.crypto.Encrypt(respHeader))
		binary.Write(&resp, binary.LittleEndian, p.SequenceId)

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent pong")

		return nil

	case OpClientReadyForAccountDataTimes:
		log.Println("starting account data times")

		inner := bytes.Buffer{}
		binary.Write(&inner, binary.LittleEndian, uint32(time.Now().Unix()))
		inner.WriteByte(1)                 // activated (bool)
		inner.Write([]byte{0, 0, 0, 0xFF}) // cache mask (all)
		// cache times
		for i := 0; i < 8; i++ {
			inner.Write([]byte{0, 0, 0, 0})
		}

		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OpServerAccountDataTimes, uint32(inner.Len()))
		if err != nil {
			return err
		}
		resp.Write(c.crypto.Encrypt(respHeader))
		resp.Write(inner.Bytes())

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent account data times")

		return nil

	case OpClientCharEnum:
		log.Println("starting character list")

		accountChars, err := s.charsDb.List(&models.CharacterListParams{
			AccountId: c.account.Id,
			RealmId:   c.realm.Id,
		})
		if err != nil {
			return err
		}

		// https://gtker.com/wow_messages/docs/smsg_char_enum.html#client-version-335
		inner := bytes.Buffer{}
		inner.WriteByte(byte(len(accountChars)))

		for _, char := range accountChars {
			binary.Write(&inner, binary.LittleEndian, uint64(char.Id))
			inner.WriteString(char.Name)
			inner.WriteByte(0) // NUL-terminated
			inner.WriteByte(byte(char.Race))
			inner.WriteByte(byte(char.Class))
			inner.WriteByte(byte(char.Gender))
			inner.WriteByte(char.SkinColor)
			inner.WriteByte(char.Face)
			inner.WriteByte(char.HairStyle)
			inner.WriteByte(char.HairColor)
			inner.WriteByte(char.FacialHair)
			inner.WriteByte(1)                                    // level
			inner.Write([]byte{12, 0, 0, 0})                      // area (hardcoded as elwynn forest)
			inner.Write([]byte{0, 0, 0, 0})                       // map (hardcoded as eastern kingdoms)
			binary.Write(&inner, binary.LittleEndian, float32(0)) // x
			binary.Write(&inner, binary.LittleEndian, float32(0)) // y
			binary.Write(&inner, binary.LittleEndian, float32(0)) // z
			inner.Write([]byte{0, 0, 0, 0})                       // guild id
			inner.Write([]byte{0, 0, 0, 0})                       // flags
			inner.Write([]byte{0, 0, 0, 0})                       // recustomization_flags (?)

			if !char.LastLogin.Valid {
				inner.WriteByte(1) // first login, show tutorial
			} else {
				inner.WriteByte(0) // not first login
			}

			inner.Write([]byte{0, 0, 0, 0}) // pet display id
			inner.Write([]byte{0, 0, 0, 0}) // pet level
			inner.Write([]byte{0, 0, 0, 0}) // pet family

			// equipment (from head to holdable)
			// https://gtker.com/wow_messages/docs/inventorytype.html
			for i := 1; i <= 23; i++ {
				inner.Write([]byte{0, 0, 0, 0}) // equipment display id
				inner.WriteByte(byte(i))        // slot
				inner.Write([]byte{0, 0, 0, 0}) // enchantment
			}
		}

		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OpServerCharEnum, uint32(inner.Len()))
		if err != nil {
			return err
		}
		resp.Write(c.crypto.Encrypt(respHeader))
		resp.Write(inner.Bytes())

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent character list")

		return nil

	case OpClientRealmSplit:
		log.Println("starting realm split")

		r := bytes.NewReader(data[6:])
		p := RealmSplitPacket{}
		binary.Read(r, binary.LittleEndian, &p.RealmId)

		// https://gtker.com/wow_messages/docs/smsg_realm_split.html
		inner := bytes.Buffer{}
		binary.Write(&inner, binary.LittleEndian, p.RealmId)
		inner.Write([]byte{0, 0, 0, 0})   // split state, 0 = normal
		inner.WriteString("01/01/01\x00") // send a bogus date (NUL-terminated)

		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OpServerRealmSplit, uint32(inner.Len()))
		if err != nil {
			return err
		}
		resp.Write(c.crypto.Encrypt(respHeader))
		resp.Write(inner.Bytes())

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent realm split")

		return nil

	case OpClientCharCreate:
		log.Println("starting character create")

		// TODO: check if account is full
		// accountChars, err := s.charsDb.List(&models.CharacterListParams{
		// 	AccountId: c.account.Id,
		// 	RealmId:   c.realm.Id,
		// })
		// if err != nil {
		// 	return err
		// }

		p := CharCreatePacket{}
		r := bytes.NewReader(data[6:])
		charName, err := readCString(r)
		if err != nil {
			return err
		}
		charName = strings.TrimSpace(charName)

		if err := binary.Read(r, binary.BigEndian, &p); err != nil {
			return err
		}

		log.Println("client wants to create character", charName)

		existing, err := s.charsDb.GetName(charName, c.realm.Id)
		if err != nil {
			return err
		}

		if existing == nil {
			char := &models.Character{
				Name:       charName,
				AccountId:  c.account.Id,
				RealmId:    c.realm.Id,
				Race:       p.Race,   // TODO
				Class:      p.Class,  // TODO
				Gender:     p.Gender, // TODO
				SkinColor:  p.SkinColor,
				Face:       p.Face,
				HairStyle:  p.HairStyle,
				HairColor:  p.HairColor,
				FacialHair: p.FacialHair,
				OutfitId:   p.OutfitId,
			}
			if err := s.charsDb.Create(char); err != nil {
				return err
			}
			log.Println("created char with id", char.Id)
		}

		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OpServerCharCreate, 1)
		if err != nil {
			return err
		}
		resp.Write(c.crypto.Encrypt(respHeader))

		if existing != nil {
			resp.WriteByte(byte(RespCodeCharCreateNameInUse))
		} else {
			resp.WriteByte(byte(RespCodeCharCreateSuccess))
		}

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("finished character create")

		return nil

	case OpClientCharDelete:
		log.Println("start character delete")

		r := bytes.NewReader(data[6:])
		p := CharDeletePacket{}
		if err := binary.Read(r, binary.LittleEndian, &p.CharacterId); err != nil {
			return err
		}

		char, err := s.charsDb.Get(uint32(p.CharacterId))
		if err != nil {
			return err
		}

		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OpServerCharDelete, 1)
		if err != nil {
			return err
		}
		resp.Write(c.crypto.Encrypt(respHeader))

		if char == nil || char.AccountId != c.account.Id || char.RealmId != c.realm.Id {
			resp.WriteByte(byte(RespCodeCharDeleteFailed))
		} else {
			if _, err := s.charsDb.Delete(char.Id); err != nil {
				return err
			}
			resp.WriteByte(byte(RespCodeCharDeleteSuccess))
		}

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("finished character create")

		return nil

	case OpClientPlayerLogin:
		log.Println("start character login")

		r := bytes.NewReader(data[6:])
		p := CharLoginPacket{}
		if err := binary.Read(r, binary.LittleEndian, &p.CharacterId); err != nil {
			return err
		}

		char, err := s.charsDb.Get(uint32(p.CharacterId))
		if err != nil {
			return err
		}

		resp := bytes.Buffer{}
		ok := char != nil && char.AccountId == c.account.Id && char.RealmId == c.realm.Id

		if !ok {
			// https: gtker.com/wow_messages/docs/smsg_character_login_failed.html#client-version-335
			respHeader, err := makeServerHeader(OpServerCharLoginFailed, 1)
			if err != nil {
				return err
			}
			resp.Write(c.crypto.Encrypt(respHeader))
			resp.WriteByte(byte(RespCodeCharLoginFailed))
		} else {
			// https://gtker.com/wow_messages/docs/smsg_login_verify_world.html
			inner := bytes.Buffer{}
			inner.Write([]byte{0, 0, 0, 0})                              // map (hardcoded as eastern kingdoms)
			binary.Write(&inner, binary.LittleEndian, float32(-8949.95)) // x
			binary.Write(&inner, binary.LittleEndian, float32(-132.493)) // y
			binary.Write(&inner, binary.LittleEndian, float32(83.5312))  // z
			binary.Write(&inner, binary.LittleEndian, float32(0))        // orientation

			respHeader, err := makeServerHeader(OpServerCharLoginVerifyWorld, uint32(inner.Len()))
			if err != nil {
				return err
			}
			resp.Write(c.crypto.Encrypt(respHeader))
			resp.Write(inner.Bytes())
		}

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent verify world")

		if ok {
			// https://gtker.com/wow_messages/docs/smsg_tutorial_flags.html
			resp := bytes.Buffer{}
			respHeader, err := makeServerHeader(OpServerTutorialFlags, 32)
			if err != nil {
				return err
			}
			resp.Write(c.crypto.Encrypt(respHeader))
			resp.Write(bytes.Repeat([]byte{255}, 32))

			if _, err := c.conn.Write(resp.Bytes()); err != nil {
				return err
			}

			log.Println("sent tutorial flags")

			// https://gtker.com/wow_messages/docs/smsg_feature_system_status.html#client-version-335
			inner := bytes.Buffer{}
			inner.WriteByte(2) // auto ignore?
			inner.WriteByte(0) // voip enabled

			resp = bytes.Buffer{}
			respHeader, err = makeServerHeader(OpServerSystemFeatures, uint32(inner.Len()))
			if err != nil {
				return err
			}
			resp.Write(c.crypto.Encrypt(respHeader))
			resp.Write(inner.Bytes())

			if _, err := c.conn.Write(resp.Bytes()); err != nil {
				return err
			}

			log.Println("sent system features")

			// https://gtker.com/wow_messages/docs/smsg_bindpointupdate.html#client-version-335
			inner = bytes.Buffer{}
			binary.Write(&inner, binary.LittleEndian, float32(-8949.95)) // hearth x
			binary.Write(&inner, binary.LittleEndian, float32(-132.493)) // hearth y
			binary.Write(&inner, binary.LittleEndian, float32(83.5312))  // hearth z
			inner.Write([]byte{0, 0, 0, 0})                              // map: eastern kingdoms
			inner.Write([]byte{12, 0, 0, 0})                             // area: elwynn forest

			resp = bytes.Buffer{}
			respHeader, err = makeServerHeader(OpServerHearthLocation, uint32(inner.Len()))
			if err != nil {
				return err
			}
			resp.Write(c.crypto.Encrypt(respHeader))
			resp.Write(inner.Bytes())

			if _, err := c.conn.Write(resp.Bytes()); err != nil {
				return err
			}

			log.Println("sent hearth location")

			// https://gtker.com/wow_messages/docs/smsg_trigger_cinematic.html
			// inner = bytes.Buffer{}
			// binary.Write(&inner, binary.LittleEndian, uint32(81)) // human

			// resp = bytes.Buffer{}
			// respHeader, err = makeServerHeader(OP_SRV_PLAY_CINEMATIC, uint32(inner.Len()))
			// if err != nil {
			// 	return err
			// }
			// resp.Write(c.crypto.Encrypt(respHeader))
			// resp.Write(inner.Bytes())

			// if _, err := c.conn.Write(resp.Bytes()); err != nil {
			// 	return err
			// }

			// log.Println("sent play cinematic")

			// https://gtker.com/wow_messages/docs/smsg_update_object.html#client-version-335
			inner = bytes.Buffer{}
			inner.Write([]byte{1, 0, 0, 0}) // number of objects

			// nested object start
			inner.WriteByte(byte(UpdateTypeCreateObject2)) // update type: CREATE_OBJECT2
			inner.Write(packGuid(uint64(char.Id)))         // packed guid
			inner.WriteByte(byte(ObjectTypePlayer))

			// movement block start
			// inner.WriteByte()
			binary.Write(&inner, binary.LittleEndian, UpdateFlagSelf|UpdateFlagLiving)
			inner.Write([]byte{0, 0, 0, 0, 0, 0})                        // movement flags
			inner.Write([]byte{0, 0, 0, 0})                              // timestamp
			binary.Write(&inner, binary.LittleEndian, float32(-8949.95)) // x
			binary.Write(&inner, binary.LittleEndian, float32(-132.493)) // y
			binary.Write(&inner, binary.LittleEndian, float32(83.5312))  // z
			binary.Write(&inner, binary.LittleEndian, float32(0))        // orientation
			inner.Write([]byte{0, 0, 0, 0})                              // fall time

			binary.Write(&inner, binary.LittleEndian, float32(1))       // walk speed
			binary.Write(&inner, binary.LittleEndian, float32(70))      // run speed
			binary.Write(&inner, binary.LittleEndian, float32(4.5))     // reverse speed
			binary.Write(&inner, binary.LittleEndian, float32(0))       // swim speed
			binary.Write(&inner, binary.LittleEndian, float32(0))       // swim reverse speed
			binary.Write(&inner, binary.LittleEndian, float32(0))       // flight speed
			binary.Write(&inner, binary.LittleEndian, float32(0))       // flight reverse speed
			binary.Write(&inner, binary.LittleEndian, float32(3.14159)) // turn speed
			binary.Write(&inner, binary.LittleEndian, float32(0))       // pitch rate
			// movement block end

			// field mask start
			updateMask := NewUpdateMask()
			valuesBuf := bytes.Buffer{}

			// Without this, client gets stuck on loading screen and floods server with 0x2CE opcode
			updateMask.SetFieldMask(FieldMaskObjectGuid)
			binary.Write(&valuesBuf, binary.LittleEndian, uint32(char.Id)) // low guid
			binary.Write(&valuesBuf, binary.LittleEndian, uint32(0))       // high guid

			// // Character seems to load fine without this
			// updateMask.SetFieldMask(FieldMaskObjectType)
			// binary.Write(&valuesBuf, binary.LittleEndian, uint32(1<<ObjectTypeObject|1<<ObjectTypeUnit|1<<ObjectTypePlayer))

			// // Without this, character model scale is zero and camera starts in first person
			// updateMask.SetFieldMask(FieldMaskObjectScaleX)
			// valuesBuf.Write([]byte{0x00, 0x00, 0x80, 0x3f})

			// // Without this, talent screen is blank
			// updateMask.SetFieldMask(FieldMaskUnitBytes0)
			// valuesBuf.WriteByte(char.Race)
			// valuesBuf.WriteByte(char.Class)
			// valuesBuf.WriteByte(char.Gender)
			// valuesBuf.WriteByte(getPowerTypeForClass(char.Class))

			// // Without this, character spawns in as a corpse
			// updateMask.SetFieldMask(FieldMaskUnitHealth)
			// valuesBuf.Write([]byte{100, 0, 0, 0})

			// // Without this, UI doesn't show max health
			// updateMask.SetFieldMask(FieldMaskUnitMaxHealth)
			// valuesBuf.Write([]byte{100, 0, 0, 0})

			// // Without this, character level appears as 0
			// updateMask.SetFieldMask(FieldMaskUnitLevel)
			// valuesBuf.Write([]byte{10, 0, 0, 0})

			// Without this, client segfaults
			updateMask.SetFieldMask(FieldMaskUnitFactionTemplate)
			valuesBuf.Write([]byte{byte(char.Race), 0, 0, 0})

			// Without this, client segfaults
			updateMask.SetFieldMask(FieldMaskUnitDisplayId)
			valuesBuf.Write([]byte{0x0C, 0x4D, 0x00, 0x00}) // human female

			// Without this, client segfaults
			updateMask.SetFieldMask(FieldMaskUnitNativeDisplayId)
			valuesBuf.Write([]byte{0x0C, 0x4D, 0x00, 0x00}) // human female

			mask := updateMask.Mask()
			inner.WriteByte(byte(len(mask)))
			binary.Write(&inner, binary.LittleEndian, mask)
			inner.Write(valuesBuf.Bytes())
			// field mask end

			// nested object end

			resp = bytes.Buffer{}
			respHeader, err = makeServerHeader(OpServerUpdateObject, uint32(inner.Len()))
			if err != nil {
				return err
			}
			resp.Write(c.crypto.Encrypt(respHeader))
			resp.Write(inner.Bytes())

			fmt.Printf("%x\n", respHeader)
			fmt.Printf("%x\n", resp.Bytes())

			if _, err := c.conn.Write(resp.Bytes()); err != nil {
				return err
			}

			log.Println("sent object update")
		}

		log.Println("finished character login")

		return nil

	default:
		log.Printf("unknown opcode: 0x%x\n", header.Opcode)
	}

	return nil
}

func (s *Server) authenticateClient(c *Client, p *AuthSessionPacket) (bool, error) {
	var err error

	if c.account, err = s.accountsDb.Get(&models.AccountGetParams{Username: p.Username}); err != nil {
		return false, err
	} else if c.account == nil {
		log.Printf("no account with username %s exists", p.Username)
		return false, nil
	}

	if c.realm, err = s.realmsDb.Get(p.RealmId); err != nil {
		return false, err
	} else if c.realm == nil {
		log.Printf("no realm with id %d exists", p.RealmId)
		return false, nil
	}

	if c.session, err = s.sessionsDb.Get(c.account.Id); err != nil {
		return false, err
	} else if c.session == nil {
		log.Printf("no session for username %s exists", c.account.Username)
		return false, nil
	}

	if err := c.session.Decode(); err != nil {
		return false, err
	}

	c.crypto = NewWrathHeaderCrypto(c.session.SessionKey())
	if err := c.crypto.Init(); err != nil {
		return false, err
	}

	proof := CalculateWorldProof(p.Username, p.ClientSeed[:], c.serverSeed[:], c.session.SessionKey())

	if !bytes.Equal(proof, p.ClientProof[:]) {
		log.Println("proofs don't match")
		log.Printf("got:    %x\n", p.ClientProof)
		log.Printf("wanted: %x\n", proof)
		return false, nil
	}

	log.Println("client authenticated successfully")

	return true, nil
}

func makeServerHeader(opcode ServerOpcode, size uint32) ([]byte, error) {
	// Include the opcode in the size
	size += 2

	if size > SizeFieldMaxValue {
		return nil, fmt.Errorf("makeServerHeader: size is too large (%d bytes)", size)
	}

	var header []byte

	// The size field in the header can be 2 or 3 bytes. The most significant bit in the size field
	// is reserved as a flag to indicate this. In total, server headers are 4 or 5 bytes.
	//
	// The header format is: <size><opcode>
	// <size> is 2-3 bytes big endian
	// <opcode> is 2 bytes little endian
	if size > LargeHeaderThreshold {
		header = []byte{
			byte(size>>16) | LargeHeaderFlag,
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

func getPowerTypeForClass(c models.Class) PowerType {
	switch c {
	case models.ClassWarrior:
		return PowerTypeRage

	case models.ClassPaladin,
		models.ClassHunter,
		models.ClassPriest,
		models.ClassShaman,
		models.ClassMage,
		models.ClassWarlock,
		models.ClassDruid:
		return PowerTypeMana

	case models.ClassRogue:
		return PowerTypeEnergy

	default:
		log.Println("getPowerTypeForClass: got unexpected class", c)
		return PowerTypeMana
	}
}

// packGuid returns a packed *little-endian* representation of an 8-byte integer. The packing works
// by creating a bit mask to mark which bytes are non-zero. Any bytes which are zero are discarded.
// The result is a byte array with the first byte as the bitmask, followed by the remaining
// undiscarded bytes. The bytes after the bitmask are little-endian.
func packGuid(val uint64) []byte {
	// At its largest, a packed guid takes up 9 bytes (1 byte mask + 8 bytes)
	result := make([]byte, 9)
	n := 0

	for i := 0; i < 8; i++ {
		if byte(val) > 0 {
			// Set the mask bit
			result[0] |= 1 << i
			// Add the byte to the result. The loop traverses the bytes from right-to-left but they
			// are written to the result from left-to-right, which swaps it to little-endian.
			result[1] = byte(val)
			n++
		}
		// Move to the next byte
		val >>= 8
	}

	return result[:n+1]
}

type UpdateMask struct {
	largestBit int
	mask       []uint32
}

func NewUpdateMask() *UpdateMask {
	return &UpdateMask{mask: make([]uint32, 16)}
}

// Mask returns the smallest []uint32 to represent all of the mask bits that were set.
func (m *UpdateMask) Mask() []uint32 {
	largestBitIndex := m.largestBit / 32
	return m.mask[:largestBitIndex+1]
}

// SetFieldMask sets all the bits necessary for the provided field mask.
func (m *UpdateMask) SetFieldMask(fieldMask FieldMask) {
	for i := 0; i < fieldMask.Size; i++ {
		m.SetBit(int(fieldMask.Offset) + i)
	}
}

// SetBit sets the nth bit in the update mask. The bit is zero-indexed with the first bit being zero.
func (m *UpdateMask) SetBit(bit int) {
	index := bit / 32
	bitPos := bit % 32
	m.resize(index)

	if bit > m.largestBit {
		m.largestBit = bit
	}

	m.mask[index] |= 1 << bitPos
}

// Resizes the mask to fit up to n uint32s.
func (m *UpdateMask) resize(n int) {
	if len(m.mask) > n {
		return
	}

	// Grow the array exponentially
	newSize := len(m.mask)
	newSize *= newSize

	// If it's still too small just use the desired size
	if newSize < n {
		newSize = n
	}

	oldMask := m.mask
	m.mask = make([]uint32, newSize)
	copy(m.mask, oldMask)
}
