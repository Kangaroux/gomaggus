package value

import "github.com/kangaroux/gomaggus/realmd"

type questLogEntry struct {
	ID    uint32
	State uint32
	Count uint64
	Time  uint32
}

type visibleItem struct {
	ID          uint32
	Enchantment uint32
}

type skillEntry struct {
	ID             uint16
	Step           uint16
	SkillLevel     uint16
	SkillCap       uint16
	TempBonus      uint16
	PermanentBonus uint16
}

// Adapted from Gophercraft with some modifications
// https://github.com/Gophercraft/core/blob/master/packet/update/d12340/descriptor.go
type Player struct {
	DuelArbiter realmd.Guid

	// flags
	GroupLeader        bool    // 0x1
	AFK                bool    // 0x2
	DND                bool    // 0x4
	GM                 bool    // 0x8
	Ghost              bool    // 0x10
	Resting            bool    // 0x20
	VoiceChat          bool    // 0x40
	FFAPVP             bool    // 0x80
	ContestedPVP       bool    // 0x100
	InPVP              bool    // 0x200
	HideHelm           bool    // 0x400
	HideCloak          bool    // 0x800
	PlayedLongTime     bool    // 0x1000
	PlayedTooLong      bool    // 0x2000
	OutOfBounds        bool    // 0x4000
	Developer          bool    // 0x8000
	_                  bool    // 0x10000
	TaxiBenchmark      bool    // 0x20000
	PVPTimer           bool    // 0x40000
	Uber               bool    // 0x80000
	_                  bool    // 0x100000
	_                  bool    // 0x200000
	Commentator        bool    // 0x400000
	OnlyAllowAbilities bool    // 0x800000
	StopMeleeOnTab     bool    // 0x1000000
	NoExperienceGain   bool    // 0x2000000
	_                  [6]bool // 0x4000000 ... 0x80000000

	GuildID   uint32
	GuildRank uint32

	// bytes1
	Skin      uint8
	Face      uint8
	HairStyle uint8
	HairColor uint8

	// bytes2
	FacialHair       uint8
	RestBits         uint8
	BankBagSlotCount uint8
	RestState        uint8

	// bytes3
	PlayerGender uint8
	GenderUnk    uint8
	Drunkness    uint8
	PVPRank      uint8

	DuelTeam       uint32
	GuildTimestamp uint32

	QuestLog [25]questLogEntry

	VisibleItems [19]visibleItem

	ChosenTitle     uint32
	FakeInebriation uint32

	_ uint32

	InventorySlots              [23]uint64
	PackSlots                   [16]uint64 // ??
	BankSlots                   [28]uint64
	BankBagSlots                [7]uint64
	VendorBuybackSlots          [12]uint64
	KeyringSlots                [32]uint64
	CurrencyTokenSlots          [32]uint64
	FarSight                    uint64
	KnownTitles                 [6]uint32
	KnownCurrencies             [2]uint32
	XP                          uint32
	NextLevelXP                 uint32
	Skills                      [128]skillEntry
	CharacterPoints             [2]uint32
	TrackCreatures              uint32
	TrackResources              uint32
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
	ExploredZones               [128]uint32
	RestStateExperience         uint32
	Wealth                      int32
	ModDamageDonePositive       [7]uint32
	ModDamageDoneNegative       [7]uint32
	ModDamageDonePercentage     [7]float32
	ModHealingDonePos           uint32
	ModHealingPercentage        float32
	ModHealingDonePercentage    float32
	ModTargetResistance         uint32
	ModTargetPhysicalResistance uint32

	// Field bytes
	_                             bool    // 0x1
	TrackStealthed                bool    // 0x2
	_                             bool    // 0x4
	DisplaySpiritAutoReleaseTimer bool    // 0x8
	HideSpiritReleaseWindow       bool    // 0x10
	_                             [3]bool // 0x20 ... 0x80
	ReferAFriendGrantableLevel    uint8
	ActionBarToggles              uint8
	LifetimeMaxPVPRank            uint8

	AmmoID                 uint32
	SelfResSpell           uint32
	PVPMedals              uint32
	BuybackPrices          [12]uint32
	BuybackTimestamps      [12]uint32
	Kills                  uint32
	TodayKills             uint32
	YesterdayKills         uint32
	LifetimeHonorableKills uint32

	// Field bytes 2
	_                              uint8 // TODO: flags
	IgnorePowerRegenPredictionMask uint8
	OverrideSpellsID               uint16

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
