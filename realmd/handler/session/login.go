package session

import (
	"bytes"
	"database/sql"
	"math"
	"time"

	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/kangaroux/gomaggus/realmd/handler/account"
	"github.com/kangaroux/gomaggus/realmd/objupdate"
	"github.com/mixcode/binarystruct"
)

// https://gtker.com/wow_messages/docs/cmsg_player_login.html
type loginRequest struct {
	CharacterId uint64
}

// https: gtker.com/wow_messages/docs/smsg_character_login_failed.html#client-version-335
type loginFailed struct {
	ResponseCode realmd.ResponseCode
}

func LoginHandler(svc *realmd.Service, client *realmd.Client, data []byte) error {
	req := loginRequest{}
	if _, err := binarystruct.Unmarshal(data, binarystruct.LittleEndian, &req); err != nil {
		return err
	}

	char, err := svc.Characters.Get(uint32(req.CharacterId))
	if err != nil {
		return err
	}

	// Notify the client something isn't right and close the connection
	if char == nil || char.AccountId != client.Account.Id || char.RealmId != client.Realm.Id {
		resp := loginFailed{ResponseCode: realmd.RespCodeCharLoginFailed}

		if err := client.SendPacket(realmd.OpServerCharLoginFailed, &resp); err != nil {
			return err
		}

		client.Log.Warn().Uint64("char", req.CharacterId).Msg("client tried logging in as invalid character")

		// TODO: set client state as invalid
		return &realmd.ErrKickClient{Reason: "invalid login"}
	}

	client.Character = char

	if err := sendVerifyWorld(client); err != nil {
		return err
	}
	if err := sendCharacterStorageTimes(svc, client); err != nil {
		return err
	}
	if err := sendTutorialFlags(client); err != nil {
		return err
	}
	if err := sendWorldTimeSpeed(client); err != nil {
		return err
	}
	if err := sendSystemFeatures(client); err != nil {
		return err
	}
	if err := sendHearthLocation(client); err != nil {
		return err
	}
	if err := sendIntroCinematic(client); err != nil {
		return err
	}
	if err := sendInitialSpells(client); err != nil {
		return err
	}
	if err := sendActionButtons(client); err != nil {
		return err
	}
	if err := sendFactionReputation(client); err != nil {
		return err
	}
	if err := sendInitialWorldStates(client); err != nil {
		return err
	}
	if err := sendSpawnPlayer(client); err != nil {
		return err
	}
	if err := sendMOTD(client); err != nil {
		return err
	}

	client.Character.LastLogin = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	if _, err := svc.Characters.Update(client.Character); err != nil {
		return err
	}

	client.Log.Info().Str("char", client.Character.String()).Msg("player login")

	return nil
}

// https://gtker.com/wow_messages/docs/smsg_login_verify_world.html
type verifyWorldResponse struct {
	Map      uint32
	Position realmd.Vector4
}

func sendVerifyWorld(client *realmd.Client) error {
	resp := verifyWorldResponse{
		Map: 0x0, // Eastern kingdoms
		Position: realmd.Vector4{
			X:        float32(-8949.95),
			Y:        float32(-132.493),
			Z:        float32(83.5312),
			Rotation: float32(0),
		},
	}

	return client.SendPacket(realmd.OpServerCharLoginVerifyWorld, &resp)
}

func sendCharacterStorageTimes(svc *realmd.Service, client *realmd.Client) error {
	h := &account.StorageTimesHandler{
		Client:  client,
		Service: svc,
	}
	return h.Send(model.AllCharacterStorage)
}

// https://gtker.com/wow_messages/docs/smsg_tutorial_flags.html
type tutorialFlags struct {
	Flags []byte `binary:"[32]byte"`
}

func sendTutorialFlags(client *realmd.Client) error {
	// Enable all tutorial flags
	resp := tutorialFlags{Flags: bytes.Repeat([]byte{0xFF}, 32)}
	return client.SendPacket(realmd.OpServerTutorialFlags, &resp)
}

type worldTimeSpeed struct {
	DateTime realmd.DateTime

	// TimeScale is the rate the in-game realm time passes. It's represented as a ratio of minutes:seconds,
	// where minutes is in-game minutes and seconds is real-world seconds. This affects how fast the
	// in-game clock ticks as well as day/night cycles.
	TimeScale float32
	Unknown   uint32
}

// https://gtker.com/wow_messages/docs/smsg_login_settimespeed.html#client-version-312-client-version-32-client-version-33
func sendWorldTimeSpeed(client *realmd.Client) error {
	resp := worldTimeSpeed{
		DateTime:  realmd.NewDateTime(time.Now()),
		TimeScale: 1.0 / 60, // 1 realm minute = 60 real world seconds
		Unknown:   0,
	}
	return client.SendPacket(realmd.OpServerSetTimeSpeed, &resp)
}

// https://gtker.com/wow_messages/docs/smsg_feature_system_status.html#client-version-335
type systemFeatures struct {
	// player reporting?
	ComplaintStatus uint8 // 0=disabled, 1=enabled (no auto ignore), 2=enabled (auto ignore)
	VoipEnabled     bool
}

func sendSystemFeatures(client *realmd.Client) error {
	resp := systemFeatures{
		ComplaintStatus: 0x2,
		VoipEnabled:     false,
	}
	return client.SendPacket(realmd.OpServerSystemFeatures, &resp)
}

// https://gtker.com/wow_messages/docs/smsg_bindpointupdate.html#client-version-335
type hearthLocation struct {
	Position realmd.Vector3
	Map      uint32
	Area     uint32
}

func sendHearthLocation(client *realmd.Client) error {
	resp := hearthLocation{
		Position: realmd.Vector3{
			X: float32(-8949.95),
			Y: float32(-132.493),
			Z: float32(83.5312),
		},
		Map:  0x0, // Eastern kingdoms
		Area: 0xB, // Elwynn forest
	}
	return client.SendPacket(realmd.OpServerHearthLocation, &resp)
}

// https://gtker.com/wow_messages/docs/smsg_trigger_cinematic.html
type playCinematic struct {
	CinematicId uint32
}

func sendIntroCinematic(client *realmd.Client) error {
	// Only play the cinematic on first login
	if client.Character.LastLogin.Valid {
		return nil
	}

	resp := playCinematic{CinematicId: 81} // 81 = human
	return client.SendPacket(realmd.OpServerPlayCinematic, &resp)
}

type spell struct {
	SpellID uint32
	Unknown uint16 // gophercraft has this as TargetFlags
}

type initialSpellsResponse struct {
	Unknown       uint8 // gophercraft has this as TalentSpec but it's unused
	SpellCount    uint16
	Spells        []spell `binary:"[SpellCount]Any"`
	CooldownCount uint16
	// TODO: CooldownSpells
}

// https://gtker.com/wow_messages/docs/smsg_initial_spells.html#client-version-335
func sendInitialSpells(client *realmd.Client) error {
	spells := []spell{
		{
			SpellID: 668, // Language (common)
			Unknown: 0,
		},
		{
			SpellID: 669, // Language (orcish)
			Unknown: 0,
		},
		{
			SpellID: 122, // Frost nova
			Unknown: 0,
		},
	}
	resp := initialSpellsResponse{
		Unknown:       0, // always zero
		SpellCount:    uint16(len(spells)),
		Spells:        spells,
		CooldownCount: 0,
	}

	return client.SendPacket(realmd.OpServerInitialSpells, &resp)
}

type actionButtonSetType uint8

const (
	buttonsInitial actionButtonSetType = 0 // Unused
	buttonsSet     actionButtonSetType = 1
	buttonsClear   actionButtonSetType = 2
)

type actionButton struct {
	// ActionPacked stores an action in the lower 24 bits and the action type in the upper 8 bits
	ActionPacked uint32
}

const (
	actionButtonCount = 144
)

// actionButtonSetResponse will overwrite all of the player's action bars
type actionButtonSetResponse struct {
	Type    actionButtonSetType
	Buttons [actionButtonCount]actionButton
}

// actionButtonClearResponse will clear all of the player's action bars
type actionButtonClearResponse struct {
	Type actionButtonSetType
}

// https://gtker.com/wow_messages/docs/smsg_action_buttons.html#client-version-335
func sendActionButtons(client *realmd.Client) error {
	resp := actionButtonSetResponse{
		// Trinity says there were issues using Initial, so Set is used instead
		Type: buttonsSet,
		// TODO: send buttons
	}
	return client.SendPacket(realmd.OpServerActionButtons, &resp)
}

const (
	FactionCount = 128
)

type faction struct {
	Flags    uint8
	Standing uint32
}

type factionReputationResponse struct {
	Count   uint32
	Faction [FactionCount]faction
}

// https://gtker.com/wow_messages/docs/smsg_initialize_factions.html#client-version-3
func sendFactionReputation(client *realmd.Client) error {
	resp := factionReputationResponse{
		Count: FactionCount,
	}

	for _, f := range resp.Faction {
		f.Flags = 0x1 // visible
	}

	return client.SendPacket(realmd.OpServerFactionReputation, &resp)
}

type worldState struct {
	ID    uint32
	Value uint32
}

type initialWorldStateResponse struct {
	Map        uint32
	Area       uint32
	SubArea    uint32
	StateCount uint16
	States     []worldState `binary:"[StateCount]Any"`
}

// World states are used for telling the client about map or zone specific information.
// A state is a key:value mapping between a 32 bit ID and a 32 bit value.
// For example, PvP battlegrounds use world states for tracking captures.
// https://gtker.com/wow_messages/docs/smsg_init_world_states.html#client-version-335
func sendInitialWorldStates(client *realmd.Client) error {
	resp := initialWorldStateResponse{
		Map:  0,  // Eastern kingdoms
		Area: 12, // Elwynn Forest
	}
	return client.SendPacket(realmd.OpServerInitialWorldStates, &resp)
}

// https://gtker.com/wow_messages/docs/smsg_update_object.html#client-version-335
// TODO: need a builder for these packets
func sendSpawnPlayer(client *realmd.Client) error {
	char := client.Character
	inner := bytes.Buffer{}
	inner.Write([]byte{1, 0, 0, 0}) // number of objects

	// nested object start
	inner.WriteByte(byte(objupdate.UpdateTypeCreateObject))
	inner.Write(realmd.PackGuid(uint64(char.Id)))
	inner.WriteByte(byte(objupdate.ObjectTypePlayer))

	movement := objupdate.MovementValues{}
	movement.Self()

	living := movement.Living()
	living.Data(&objupdate.LivingData{
		Timestamp: uint32(time.Now().UnixMilli()),
		PositionRotation: realmd.Vector4{
			X:        float32(-8949.95),
			Y:        float32(-132.493),
			Z:        float32(83.5312),
			Rotation: float32(0),
		},
		FallTime:           0,
		WalkSpeed:          1,
		RunSpeed:           70,
		ReverseSpeed:       4.5,
		SwimSpeed:          0,
		SwimReverseSpeed:   0,
		FlightSpeed:        0,
		FlightReverseSpeed: 0,
		TurnRate:           math.Pi,
		PitchRate:          0,
	})

	inner.Write(movement.Bytes())

	values := objupdate.Values{}
	obj := values.Object()
	obj.Guid(realmd.Guid(char.Id))
	obj.Type(objupdate.ObjectTypeObject, objupdate.ObjectTypeUnit, objupdate.ObjectTypePlayer)
	obj.ScaleX(1)

	unit := values.Unit()
	unit.Race(char.Race)
	unit.Class(char.Class)
	unit.Gender(char.Gender)
	unit.PowerType(realmd.PowerTypeForClass(char.Class))
	unit.Health(100)
	unit.MaxHealth(100)
	unit.Level(1)
	unit.Faction(char.Race)
	unit.DisplayModel(0x4D0C)       // human female
	unit.NativeDisplayModel(0x4D0C) // human female
	unit.Flags(objupdate.PlayerControlled | objupdate.AurasVisible)
	unit.Agility(10)
	unit.Intellect(10)
	unit.Stamina(10)
	unit.Strength(10)
	unit.Spirit(10)

	inner.Write(values.Bytes())

	return client.SendPacketBytes(realmd.OpServerUpdateObject, inner.Bytes())
}

type motdResponse struct {
	MOTDCount uint32
	MOTD      []string `binary:"[MOTDCount]zstring"`
}

// https://gtker.com/wow_messages/docs/smsg_motd.html#client-version-243-client-version-3
func sendMOTD(client *realmd.Client) error {
	motd := []string{
		"Testing 123",
	}
	resp := motdResponse{
		MOTDCount: uint32(len(motd)),
		MOTD:      motd,
	}
	return client.SendPacket(realmd.OpServerMOTD, &resp)
}
