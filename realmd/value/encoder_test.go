package value

import (
	"testing"
)

type PlayerData struct {
	DuelArbiter struct {
		High uint32
		Low  uint32
	}
	GroupLeader    bool
	AFK            bool
	DND            bool
	GM             bool
	Ghost          bool
	Resting        bool
	VoiceChat      bool
	FFAPVP         bool
	ContestedPVP   bool
	PVPDesired     bool
	HideHelm       bool
	HideCloak      bool
	PlayedLongTime bool
	PlayedTooLong  bool
	OutOfBounds    bool
	GhostEffect    bool
	Sanctuary      bool
	TaxiBenchmark  bool
	PVPTimer       bool
	_              [13]bool
	_              uint32

	GuildID          uint32
	GuildRank        uint32
	Skin             uint8
	Face             uint8
	HairStyle        uint8
	HairColor        uint8
	FacialHair       uint8
	RestBits         uint8
	BankBagSlotCount uint8
	RestState        uint8
	PlayerGender     uint8
	GenderUnk        uint8
	Drunkness        uint8
	PVPRank          uint8
	DuelTeam         uint32
	GuildTimestamp   uint32
	QuestLog         [25]struct {
		QuestID    uint32
		CountState uint32
		QuestUnk   uint32
		QuestUnk2  uint32
		Time       uint32
	}

	VisibleItems [19]struct {
		Entry       uint32
		Enchantment uint32
	}

	ChosenTitle     uint32
	FakeInebriation uint32

	// StartSlotPad       update.ChunkPad
	StartSlotPad       uint32
	InventorySlots     [39]int32 `update:"private"`
	BankSlots          [28]int32 `update:"private"`
	BankBagSlots       [7]int32  `update:"private"`
	VendorBuybackSlots [12]int32 `update:"private"`
	KeyringSlots       [32]int32 `update:"private"`
	CurrencyTokenSlots [32]int32 `update:"private"`
	FarSight           int32
	KnownTitles        [6]uint32
	KnownCurrencies    [2]uint32
	XP                 uint32
	NextLevelXP        uint32
	SkillInfos         [128]struct {
		ID         uint16
		Step       uint16
		SkillLevel uint16
		SkillCap   uint16
		Bonus      uint32
	} `update:"private"`
	CharacterPoints             [2]uint32 `update:"private"`
	TrackCreatures              uint32    `update:"private"`
	TrackResources              uint32    `update:"private"`
	BlockPercentage             float32
	DodgePercentage             float32
	ParryPercentage             float32
	Expertise                   uint32
	OffhandExpertise            uint32
	CritPercentage              float32
	RangedCritPercentage        float32
	OffhandCritPercentage       float32
	SpellCritPercentage         [7]float32
	ShieldBlock                 uint32
	ShieldBlockCritPercentage   float32
	ExploredZones               [128]uint32 // TODO: use Bitmask type with length tag to refer to this field.
	RestStateExperience         uint32
	Coinage                     int32 `update:"private"`
	ModDamageDonePositive       [7]uint32
	ModDamageDoneNegative       [7]uint32
	ModDamageDonePercentage     [7]float32
	ModHealingDonePos           uint32
	ModHealingPercentage        float32
	ModHealingDonePercentage    float32
	ModTargetResistance         uint32
	ModTargetPhysicalResistance uint32
	// Flags
	PlayerFieldBytes0UnkBit0      bool
	TrackStealthed                bool
	DisplaySpiritAutoReleaseTimer bool
	HideSpiritReleaseWindow       bool
	_                             [4]bool
	RAFGrantableLevel             uint8 // parser should automatically frame this to next byte.
	ActionBarToggles              uint8
	LifetimeMaxPVPRank            uint8

	AmmoID                 uint32
	SelfResSpell           uint32
	PVPMedals              uint32
	BuybackPrices          [12]uint32 `update:"private"`
	BuybackTimestamps      [12]uint32 `update:"private"`
	Kills                  uint32
	TodayKills             uint32
	YesterdayKills         uint32
	LifetimeHonorableKills uint32
	HonorRankPoints        uint8
	DetectionFlagUnk       bool
	DetectAmore0           bool
	DetectAmore1           bool
	DetectAmore2           bool
	DetectAmore3           bool
	DetectStealth          bool
	DetectInvisibilityGlow bool
	_                      bool
	_                      uint16

	WatchedFactionIndex int32
	CombatRatings       [25]uint32
	ArenaTeamInfo       [21]uint32
	HonorCurrency       uint32
	ArenaCurrency       uint32

	MaxLevel      uint32
	DailyQuests   [25]uint32
	RuneRegen     [4]float32
	NoReagentCost [3]uint32
	GlyphSlots    [6]uint32
	Glyphs        [6]uint32
	GlyphsEnabled uint32
	PetSpellPower uint32
}

func TestV2(t *testing.T) {
	e := &encoder{}

	// type Foo struct {
	// 	A uint16
	// 	B uint8
	// 	C uint16
	// 	D uint8
	// 	E uint8
	// }

	// fmt.Println("first", e.Encode(Foo{1, 2, 3, 4, 5}))

	// fmt.Printf("second %08b\n", e.Encode(struct {
	// 	A     bool
	// 	W     uint32
	// 	flags [7]bool
	// 	L     uint16
	// 	I     uint8
	// 	X     uint8
	// 	Y     uint16
	// }{true, ^uint32(0), [7]bool{false, true, false, true, false, true, false}, 3, 5, 6, 7}))

	e.Encode(&PlayerData{})
}

func BenchmarkDeez(b *testing.B) {
	e := encoder{}

	for i := 0; i < b.N; i++ {
		e.Encode(&PlayerData{})
	}
}
