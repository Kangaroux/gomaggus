package values

import (
	"reflect"

	"github.com/kangaroux/gomaggus/realmd"
)

const (
	PlayerDataSize = 1178
)

type QuestLogEntry struct {
	ID    uint32
	State uint32
	Count uint64
	Time  uint32
}

type VisibleItem struct {
	ID          uint32
	Enchantment uint32
}

type SkillEntry struct {
	ID             uint16
	Step           uint16
	SkillLevel     uint16
	SkillCap       uint16
	TempBonus      uint16
	PermanentBonus uint16
}

// Adapted from Gophercraft with some modifications
// https://github.com/Gophercraft/core/blob/master/packet/update/d12340/descriptor.go
type PlayerData struct {
	duelArbiter realmd.Guid

	// flags
	groupLeader        bool    // 0x1
	afk                bool    // 0x2
	dnd                bool    // 0x4
	gm                 bool    // 0x8
	ghost              bool    // 0x10
	resting            bool    // 0x20
	voiceChat          bool    // 0x40
	ffapvp             bool    // 0x80
	contestedPVP       bool    // 0x100
	inPVP              bool    // 0x200
	hideHelm           bool    // 0x400
	hideCloak          bool    // 0x800
	playedLongTime     bool    // 0x1000
	playedTooLong      bool    // 0x2000
	outOfBounds        bool    // 0x4000
	developer          bool    // 0x8000
	_                  bool    // 0x10000
	taxiBenchmark      bool    // 0x20000
	pvpTimer           bool    // 0x40000
	uber               bool    // 0x80000
	_                  bool    // 0x100000
	_                  bool    // 0x200000
	commentator        bool    // 0x400000
	onlyAllowAbilities bool    // 0x800000
	stopMeleeOnTab     bool    // 0x1000000
	noExperienceGain   bool    // 0x2000000
	_                  [6]bool // 0x4000000 ... 0x80000000

	guildID   uint32
	guildRank uint32

	// bytes1
	skin      uint8
	face      uint8
	hairStyle uint8
	hairColor uint8

	// bytes2
	facialHair       uint8
	restBits         uint8
	bankBagSlotCount uint8
	restState        uint8

	// bytes3
	playerGender uint8
	genderUnk    uint8
	drunkness    uint8
	pVPRank      uint8

	duelTeam       uint32
	guildTimestamp uint32

	questLog [25]QuestLogEntry

	visibleItems [19]VisibleItem

	chosenTitle     uint32
	fakeInebriation uint32

	_ uint32

	inventorySlots            [23]uint64
	packSlots                 [16]uint64 // ??
	bankSlots                 [28]uint64
	bankBagSlots              [7]uint64
	vendorBuybackSlots        [12]uint64
	keyringSlots              [32]uint64
	currencyTokenSlots        [32]uint64
	farSight                  uint64
	knownTitles               [6]uint32
	knownCurrencies           [2]uint32
	xp                        uint32
	nextLevelXP               uint32
	skills                    [128]SkillEntry
	characterPoints           [2]uint32
	trackCreatures            uint32
	trackResources            uint32
	blockPercentage           float32
	dodgePercentage           float32
	parryPercentage           float32
	expertise                 uint32
	offhandExpertise          uint32
	critPercentage            float32
	rangedCritPercentage      float32
	offhandCritPercentage     float32
	spellCritPercentage       [7]float32
	shieldBlock               uint32
	shieldBlockCritPercentage float32
	exploredZones             [128]uint32
	restStateExperience       uint32
	wealth                    int32
	modDamageDonePositive     [7]uint32
	modDamageDoneNegative     [7]uint32
	modDamageDonePercentage   [7]float32
	modHealingDonePos         uint32
	modHealingPercentage      float32
	modHealingDonePercentage  float32
	modTarResistance          uint32
	modTarPhysicalResistance  uint32

	// Field bytes
	_                             bool    // 0x1
	trackStealthed                bool    // 0x2
	_                             bool    // 0x4
	displaySpiritAutoReleaSetImer bool    // 0x8
	hideSpiritReleaseWindow       bool    // 0x10
	_                             [3]bool // 0x20 ... 0x80
	referAFriendGrantableLevel    uint8
	actionBarToggles              uint8
	lifetimeMaxPVPRank            uint8

	ammoID                 uint32
	selfResSpell           uint32
	pVPMedals              uint32
	buybackPrices          [12]uint32
	buybackTimestamps      [12]uint32
	kills                  uint32
	todayKills             uint32
	yesterdayKills         uint32
	lifetimeHonorableKills uint32

	// Field bytes 2
	_                              uint8 // TODO: flags
	ignorePowerRegenPredictionMask uint8
	overrideSpellsID               uint16

	watchedFactionIndex int32
	combatRatings       [25]uint32
	arenaTeamInfo       [21]uint32
	honorCurrency       uint32
	arenaCurrency       uint32

	maxLevel      uint32
	dailyQuests   [25]uint32
	runeRegen     [4]float32
	noReagentCost [3]uint32
	glyphSlots    [6]uint32
	glyphs        [6]uint32
	glyphsEnabled uint32
	petSpellPower uint32

	dirty *dirtyValues `value:"END"`
}

func NewPlayerData() *PlayerData {
	return &PlayerData{
		dirty: newDirtyValues(getStructLayout(reflect.ValueOf(PlayerData{}))),
	}
}

func (p *PlayerData) Marshal(onlyDirty bool) ([]byte, []structSection) {
	return marshalValues(p, onlyDirty, p.dirty)
}

func (p *PlayerData) DuelArbiter() realmd.Guid {
	return p.duelArbiter
}

func (p *PlayerData) SetDuelArbiter(val realmd.Guid) {
	p.duelArbiter = val
	p.dirty.Flag("duelArbiter")
}

func (p *PlayerData) GroupLeader() bool {
	return p.groupLeader
}

func (p *PlayerData) SetGroupLeader(val bool) {
	p.groupLeader = val
	p.dirty.Flag("groupLeader")
}

func (p *PlayerData) AFK() bool {
	return p.afk
}

func (p *PlayerData) SetAFK(val bool) {
	p.afk = val
	p.dirty.Flag("afk")
}

func (p *PlayerData) DND() bool {
	return p.dnd
}

func (p *PlayerData) SetDND(val bool) {
	p.dnd = val
	p.dirty.Flag("dnd")
}

func (p *PlayerData) GM() bool {
	return p.gm
}

func (p *PlayerData) SetGM(val bool) {
	p.gm = val
	p.dirty.Flag("gm")
}

func (p *PlayerData) Ghost() bool {
	return p.ghost
}

func (p *PlayerData) SetGhost(val bool) {
	p.ghost = val
	p.dirty.Flag("ghost")
}

func (p *PlayerData) Resting() bool {
	return p.resting
}

func (p *PlayerData) SetResting(val bool) {
	p.resting = val
	p.dirty.Flag("resting")
}

func (p *PlayerData) VoiceChat() bool {
	return p.voiceChat
}

func (p *PlayerData) SetVoiceChat(val bool) {
	p.voiceChat = val
	p.dirty.Flag("voiceChat")
}

func (p *PlayerData) FFAPVP() bool {
	return p.ffapvp
}

func (p *PlayerData) SetFFAPVP(val bool) {
	p.ffapvp = val
	p.dirty.Flag("ffapvp")
}

func (p *PlayerData) ContestedPVP() bool {
	return p.contestedPVP
}

func (p *PlayerData) SetContestedPVP(val bool) {
	p.contestedPVP = val
	p.dirty.Flag("contestedPVP")
}

func (p *PlayerData) InPVP() bool {
	return p.inPVP
}

func (p *PlayerData) SetInPVP(val bool) {
	p.inPVP = val
	p.dirty.Flag("inPVP")
}

func (p *PlayerData) HideHelm() bool {
	return p.hideHelm
}

func (p *PlayerData) SetHideHelm(val bool) {
	p.hideHelm = val
	p.dirty.Flag("hideHelm")
}

func (p *PlayerData) HideCloak() bool {
	return p.hideCloak
}

func (p *PlayerData) SetHideCloak(val bool) {
	p.hideCloak = val
	p.dirty.Flag("hideCloak")
}

func (p *PlayerData) PlayedLongTime() bool {
	return p.playedLongTime
}

func (p *PlayerData) SetPlayedLongTime(val bool) {
	p.playedLongTime = val
	p.dirty.Flag("playedLongTime")
}

func (p *PlayerData) PlayedTooLong() bool {
	return p.playedTooLong
}

func (p *PlayerData) SetPlayedTooLong(val bool) {
	p.playedTooLong = val
	p.dirty.Flag("playedTooLong")
}

func (p *PlayerData) OutOfBounds() bool {
	return p.outOfBounds
}

func (p *PlayerData) SetOutOfBounds(val bool) {
	p.outOfBounds = val
	p.dirty.Flag("outOfBounds")
}

func (p *PlayerData) Developer() bool {
	return p.developer
}

func (p *PlayerData) SetDeveloper(val bool) {
	p.developer = val
	p.dirty.Flag("developer")
}

func (p *PlayerData) TaxiBenchmark() bool {
	return p.taxiBenchmark
}

func (p *PlayerData) SetTaxiBenchmark(val bool) {
	p.taxiBenchmark = val
	p.dirty.Flag("taxiBenchmark")
}

func (p *PlayerData) PVPTimer() bool {
	return p.pvpTimer
}

func (p *PlayerData) SetPVPTimer(val bool) {
	p.pvpTimer = val
	p.dirty.Flag("pvpTimer")
}

func (p *PlayerData) Uber() bool {
	return p.uber
}

func (p *PlayerData) SetUber(val bool) {
	p.uber = val
	p.dirty.Flag("uber")
}

func (p *PlayerData) Commentator() bool {
	return p.commentator
}

func (p *PlayerData) SetCommentator(val bool) {
	p.commentator = val
	p.dirty.Flag("commentator")
}

func (p *PlayerData) OnlyAllowAbilities() bool {
	return p.onlyAllowAbilities
}

func (p *PlayerData) SetOnlyAllowAbilities(val bool) {
	p.onlyAllowAbilities = val
	p.dirty.Flag("onlyAllowAbilities")
}

func (p *PlayerData) StopMeleeOnTab() bool {
	return p.stopMeleeOnTab
}

func (p *PlayerData) SetStopMeleeOnTab(val bool) {
	p.stopMeleeOnTab = val
	p.dirty.Flag("stopMeleeOnTab")
}

func (p *PlayerData) NoExperienceGain() bool {
	return p.noExperienceGain
}

func (p *PlayerData) SetNoExperienceGain(val bool) {
	p.noExperienceGain = val
	p.dirty.Flag("noExperienceGain")
}

func (p *PlayerData) GuildID() uint32 {
	return p.guildID
}

func (p *PlayerData) SetGuildID(val uint32) {
	p.guildID = val
	p.dirty.Flag("guildID")
}

func (p *PlayerData) GuildRank() uint32 {
	return p.guildRank
}

func (p *PlayerData) SetGuildRank(val uint32) {
	p.guildRank = val
	p.dirty.Flag("guildRank")
}

func (p *PlayerData) Skin() uint8 {
	return p.skin
}

func (p *PlayerData) SetSkin(val uint8) {
	p.skin = val
	p.dirty.Flag("skin")
}

func (p *PlayerData) Face() uint8 {
	return p.face
}

func (p *PlayerData) SetFace(val uint8) {
	p.face = val
	p.dirty.Flag("face")
}

func (p *PlayerData) HairStyle() uint8 {
	return p.hairStyle
}

func (p *PlayerData) SetHairStyle(val uint8) {
	p.hairStyle = val
	p.dirty.Flag("hairStyle")
}

func (p *PlayerData) HairColor() uint8 {
	return p.hairColor
}

func (p *PlayerData) SetHairColor(val uint8) {
	p.hairColor = val
	p.dirty.Flag("hairColor")
}

func (p *PlayerData) FacialHair() uint8 {
	return p.facialHair
}

func (p *PlayerData) SetFacialHair(val uint8) {
	p.facialHair = val
	p.dirty.Flag("facialHair")
}

func (p *PlayerData) RestBits() uint8 {
	return p.restBits
}

func (p *PlayerData) SetRestBits(val uint8) {
	p.restBits = val
	p.dirty.Flag("restBits")
}

func (p *PlayerData) BankBagSlotCount() uint8 {
	return p.bankBagSlotCount
}

func (p *PlayerData) SetBankBagSlotCount(val uint8) {
	p.bankBagSlotCount = val
	p.dirty.Flag("bankBagSlotCount")
}

func (p *PlayerData) RestState() uint8 {
	return p.restState
}

func (p *PlayerData) SetRestState(val uint8) {
	p.restState = val
	p.dirty.Flag("restState")
}

func (p *PlayerData) PlayerGender() uint8 {
	return p.playerGender
}

func (p *PlayerData) SetPlayerGender(val uint8) {
	p.playerGender = val
	p.dirty.Flag("playerGender")
}

func (p *PlayerData) GenderUnk() uint8 {
	return p.genderUnk
}

func (p *PlayerData) SetGenderUnk(val uint8) {
	p.genderUnk = val
	p.dirty.Flag("genderUnk")
}

func (p *PlayerData) Drunkness() uint8 {
	return p.drunkness
}

func (p *PlayerData) SetDrunkness(val uint8) {
	p.drunkness = val
	p.dirty.Flag("drunkness")
}

func (p *PlayerData) PVPRank() uint8 {
	return p.pVPRank
}

func (p *PlayerData) SetPVPRank(val uint8) {
	p.pVPRank = val
	p.dirty.Flag("pVPRank")
}

func (p *PlayerData) DuelTeam() uint32 {
	return p.duelTeam
}

func (p *PlayerData) SetDuelTeam(val uint32) {
	p.duelTeam = val
	p.dirty.Flag("duelTeam")
}

func (p *PlayerData) GuildTimestamp() uint32 {
	return p.guildTimestamp
}

func (p *PlayerData) SetGuildTimestamp(val uint32) {
	p.guildTimestamp = val
	p.dirty.Flag("guildTimestamp")
}

func (p *PlayerData) QuestLog() [25]QuestLogEntry {
	return p.questLog
}

func (p *PlayerData) SetQuestLog(val [25]QuestLogEntry) {
	p.questLog = val
	p.dirty.Flag("questLog")
}

func (p *PlayerData) VisibleItems() [19]VisibleItem {
	return p.visibleItems
}

func (p *PlayerData) SetVisibleItems(val [19]VisibleItem) {
	p.visibleItems = val
	p.dirty.Flag("visibleItems")
}

func (p *PlayerData) ChosenTitle() uint32 {
	return p.chosenTitle
}

func (p *PlayerData) SetChosenTitle(val uint32) {
	p.chosenTitle = val
	p.dirty.Flag("chosenTitle")
}

func (p *PlayerData) FakeInebriation() uint32 {
	return p.fakeInebriation
}

func (p *PlayerData) SetFakeInebriation(val uint32) {
	p.fakeInebriation = val
	p.dirty.Flag("fakeInebriation")
}

func (p *PlayerData) InventorySlots() [23]uint64 {
	return p.inventorySlots
}

func (p *PlayerData) SetInventorySlots(val [23]uint64) {
	p.inventorySlots = val
	p.dirty.Flag("inventorySlots")
}

func (p *PlayerData) PackSlots() [16]uint64 {
	return p.packSlots
}

func (p *PlayerData) SetPackSlots(val [16]uint64) {
	p.packSlots = val
	p.dirty.Flag("packSlots")
}

func (p *PlayerData) BankSlots() [28]uint64 {
	return p.bankSlots
}

func (p *PlayerData) SetBankSlots(val [28]uint64) {
	p.bankSlots = val
	p.dirty.Flag("bankSlots")
}

func (p *PlayerData) BankBagSlots() [7]uint64 {
	return p.bankBagSlots
}

func (p *PlayerData) SetBankBagSlots(val [7]uint64) {
	p.bankBagSlots = val
	p.dirty.Flag("bankBagSlots")
}

func (p *PlayerData) VendorBuybackSlots() [12]uint64 {
	return p.vendorBuybackSlots
}

func (p *PlayerData) SetVendorBuybackSlots(val [12]uint64) {
	p.vendorBuybackSlots = val
	p.dirty.Flag("vendorBuybackSlots")
}

func (p *PlayerData) KeyringSlots() [32]uint64 {
	return p.keyringSlots
}

func (p *PlayerData) SetKeyringSlots(val [32]uint64) {
	p.keyringSlots = val
	p.dirty.Flag("keyringSlots")
}

func (p *PlayerData) CurrencyTokenSlots() [32]uint64 {
	return p.currencyTokenSlots
}

func (p *PlayerData) SetCurrencyTokenSlots(val [32]uint64) {
	p.currencyTokenSlots = val
	p.dirty.Flag("currencyTokenSlots")
}

func (p *PlayerData) FarSight() uint64 {
	return p.farSight
}

func (p *PlayerData) SetFarSight(val uint64) {
	p.farSight = val
	p.dirty.Flag("farSight")
}

func (p *PlayerData) KnownTitles() [6]uint32 {
	return p.knownTitles
}

func (p *PlayerData) SetKnownTitles(val [6]uint32) {
	p.knownTitles = val
	p.dirty.Flag("knownTitles")
}

func (p *PlayerData) KnownCurrencies() [2]uint32 {
	return p.knownCurrencies
}

func (p *PlayerData) SetKnownCurrencies(val [2]uint32) {
	p.knownCurrencies = val
	p.dirty.Flag("knownCurrencies")
}

func (p *PlayerData) Xp() uint32 {
	return p.xp
}

func (p *PlayerData) SetXp(val uint32) {
	p.xp = val
	p.dirty.Flag("xp")
}

func (p *PlayerData) NextLevelXP() uint32 {
	return p.nextLevelXP
}

func (p *PlayerData) SetNextLevelXP(val uint32) {
	p.nextLevelXP = val
	p.dirty.Flag("nextLevelXP")
}

func (p *PlayerData) Skills() [128]SkillEntry {
	return p.skills
}

func (p *PlayerData) SetSkills(val [128]SkillEntry) {
	p.skills = val
	p.dirty.Flag("skills")
}

func (p *PlayerData) CharacterPoints() [2]uint32 {
	return p.characterPoints
}

func (p *PlayerData) SetCharacterPoints(val [2]uint32) {
	p.characterPoints = val
	p.dirty.Flag("characterPoints")
}

func (p *PlayerData) TrackCreatures() uint32 {
	return p.trackCreatures
}

func (p *PlayerData) SetTrackCreatures(val uint32) {
	p.trackCreatures = val
	p.dirty.Flag("trackCreatures")
}

func (p *PlayerData) TrackResources() uint32 {
	return p.trackResources
}

func (p *PlayerData) SetTrackResources(val uint32) {
	p.trackResources = val
	p.dirty.Flag("trackResources")
}

func (p *PlayerData) BlockPercentage() float32 {
	return p.blockPercentage
}

func (p *PlayerData) SetBlockPercentage(val float32) {
	p.blockPercentage = val
	p.dirty.Flag("blockPercentage")
}

func (p *PlayerData) DodgePercentage() float32 {
	return p.dodgePercentage
}

func (p *PlayerData) SetDodgePercentage(val float32) {
	p.dodgePercentage = val
	p.dirty.Flag("dodgePercentage")
}

func (p *PlayerData) ParryPercentage() float32 {
	return p.parryPercentage
}

func (p *PlayerData) SetParryPercentage(val float32) {
	p.parryPercentage = val
	p.dirty.Flag("parryPercentage")
}

func (p *PlayerData) Expertise() uint32 {
	return p.expertise
}

func (p *PlayerData) SetExpertise(val uint32) {
	p.expertise = val
	p.dirty.Flag("expertise")
}

func (p *PlayerData) OffhandExpertise() uint32 {
	return p.offhandExpertise
}

func (p *PlayerData) SetOffhandExpertise(val uint32) {
	p.offhandExpertise = val
	p.dirty.Flag("offhandExpertise")
}

func (p *PlayerData) CritPercentage() float32 {
	return p.critPercentage
}

func (p *PlayerData) SetCritPercentage(val float32) {
	p.critPercentage = val
	p.dirty.Flag("critPercentage")
}

func (p *PlayerData) RangedCritPercentage() float32 {
	return p.rangedCritPercentage
}

func (p *PlayerData) SetRangedCritPercentage(val float32) {
	p.rangedCritPercentage = val
	p.dirty.Flag("rangedCritPercentage")
}

func (p *PlayerData) OffhandCritPercentage() float32 {
	return p.offhandCritPercentage
}

func (p *PlayerData) SetOffhandCritPercentage(val float32) {
	p.offhandCritPercentage = val
	p.dirty.Flag("offhandCritPercentage")
}

func (p *PlayerData) SpellCritPercentage() [7]float32 {
	return p.spellCritPercentage
}

func (p *PlayerData) SetSpellCritPercentage(val [7]float32) {
	p.spellCritPercentage = val
	p.dirty.Flag("spellCritPercentage")
}

func (p *PlayerData) ShieldBlock() uint32 {
	return p.shieldBlock
}

func (p *PlayerData) SetShieldBlock(val uint32) {
	p.shieldBlock = val
	p.dirty.Flag("shieldBlock")
}

func (p *PlayerData) ShieldBlockCritPercentage() float32 {
	return p.shieldBlockCritPercentage
}

func (p *PlayerData) SetShieldBlockCritPercentage(val float32) {
	p.shieldBlockCritPercentage = val
	p.dirty.Flag("shieldBlockCritPercentage")
}

func (p *PlayerData) ExploredZones() [128]uint32 {
	return p.exploredZones
}

func (p *PlayerData) SetExploredZones(val [128]uint32) {
	p.exploredZones = val
	p.dirty.Flag("exploredZones")
}

func (p *PlayerData) RestStateExperience() uint32 {
	return p.restStateExperience
}

func (p *PlayerData) SetRestStateExperience(val uint32) {
	p.restStateExperience = val
	p.dirty.Flag("restStateExperience")
}

func (p *PlayerData) Wealth() int32 {
	return p.wealth
}

func (p *PlayerData) SetWealth(val int32) {
	p.wealth = val
	p.dirty.Flag("wealth")
}

func (p *PlayerData) ModDamageDonePositive() [7]uint32 {
	return p.modDamageDonePositive
}

func (p *PlayerData) SetModDamageDonePositive(val [7]uint32) {
	p.modDamageDonePositive = val
	p.dirty.Flag("modDamageDonePositive")
}

func (p *PlayerData) ModDamageDoneNegative() [7]uint32 {
	return p.modDamageDoneNegative
}

func (p *PlayerData) SetModDamageDoneNegative(val [7]uint32) {
	p.modDamageDoneNegative = val
	p.dirty.Flag("modDamageDoneNegative")
}

func (p *PlayerData) ModDamageDonePercentage() [7]float32 {
	return p.modDamageDonePercentage
}

func (p *PlayerData) SetModDamageDonePercentage(val [7]float32) {
	p.modDamageDonePercentage = val
	p.dirty.Flag("modDamageDonePercentage")
}

func (p *PlayerData) ModHealingDonePos() uint32 {
	return p.modHealingDonePos
}

func (p *PlayerData) SetModHealingDonePos(val uint32) {
	p.modHealingDonePos = val
	p.dirty.Flag("modHealingDonePos")
}

func (p *PlayerData) ModHealingPercentage() float32 {
	return p.modHealingPercentage
}

func (p *PlayerData) SetModHealingPercentage(val float32) {
	p.modHealingPercentage = val
	p.dirty.Flag("modHealingPercentage")
}

func (p *PlayerData) ModHealingDonePercentage() float32 {
	return p.modHealingDonePercentage
}

func (p *PlayerData) SetModHealingDonePercentage(val float32) {
	p.modHealingDonePercentage = val
	p.dirty.Flag("modHealingDonePercentage")
}

func (p *PlayerData) ModTarResistance() uint32 {
	return p.modTarResistance
}

func (p *PlayerData) SetModTarResistance(val uint32) {
	p.modTarResistance = val
	p.dirty.Flag("modTarResistance")
}

func (p *PlayerData) ModTarPhysicalResistance() uint32 {
	return p.modTarPhysicalResistance
}

func (p *PlayerData) SetModTarPhysicalResistance(val uint32) {
	p.modTarPhysicalResistance = val
	p.dirty.Flag("modTarPhysicalResistance")
}

func (p *PlayerData) TrackStealthed() bool {
	return p.trackStealthed
}

func (p *PlayerData) SetTrackStealthed(val bool) {
	p.trackStealthed = val
	p.dirty.Flag("trackStealthed")
}

func (p *PlayerData) DisplaySpiritAutoReleaSetImer() bool {
	return p.displaySpiritAutoReleaSetImer
}

func (p *PlayerData) SetDisplaySpiritAutoReleaSetImer(val bool) {
	p.displaySpiritAutoReleaSetImer = val
	p.dirty.Flag("displaySpiritAutoReleaSetImer")
}

func (p *PlayerData) HideSpiritReleaseWindow() bool {
	return p.hideSpiritReleaseWindow
}

func (p *PlayerData) SetHideSpiritReleaseWindow(val bool) {
	p.hideSpiritReleaseWindow = val
	p.dirty.Flag("hideSpiritReleaseWindow")
}

func (p *PlayerData) ReferAFriendGrantableLevel() uint8 {
	return p.referAFriendGrantableLevel
}

func (p *PlayerData) SetReferAFriendGrantableLevel(val uint8) {
	p.referAFriendGrantableLevel = val
	p.dirty.Flag("referAFriendGrantableLevel")
}

func (p *PlayerData) ActionBarToggles() uint8 {
	return p.actionBarToggles
}

func (p *PlayerData) SetActionBarToggles(val uint8) {
	p.actionBarToggles = val
	p.dirty.Flag("actionBarToggles")
}

func (p *PlayerData) LifetimeMaxPVPRank() uint8 {
	return p.lifetimeMaxPVPRank
}

func (p *PlayerData) SetLifetimeMaxPVPRank(val uint8) {
	p.lifetimeMaxPVPRank = val
	p.dirty.Flag("lifetimeMaxPVPRank")
}

func (p *PlayerData) AmmoID() uint32 {
	return p.ammoID
}

func (p *PlayerData) SetAmmoID(val uint32) {
	p.ammoID = val
	p.dirty.Flag("ammoID")
}

func (p *PlayerData) SelfResSpell() uint32 {
	return p.selfResSpell
}

func (p *PlayerData) SetSelfResSpell(val uint32) {
	p.selfResSpell = val
	p.dirty.Flag("selfResSpell")
}

func (p *PlayerData) PVPMedals() uint32 {
	return p.pVPMedals
}

func (p *PlayerData) SetPVPMedals(val uint32) {
	p.pVPMedals = val
	p.dirty.Flag("pVPMedals")
}

func (p *PlayerData) BuybackPrices() [12]uint32 {
	return p.buybackPrices
}

func (p *PlayerData) SetBuybackPrices(val [12]uint32) {
	p.buybackPrices = val
	p.dirty.Flag("buybackPrices")
}

func (p *PlayerData) BuybackTimestamps() [12]uint32 {
	return p.buybackTimestamps
}

func (p *PlayerData) SetBuybackTimestamps(val [12]uint32) {
	p.buybackTimestamps = val
	p.dirty.Flag("buybackTimestamps")
}

func (p *PlayerData) Kills() uint32 {
	return p.kills
}

func (p *PlayerData) SetKills(val uint32) {
	p.kills = val
	p.dirty.Flag("kills")
}

func (p *PlayerData) TodayKills() uint32 {
	return p.todayKills
}

func (p *PlayerData) SetTodayKills(val uint32) {
	p.todayKills = val
	p.dirty.Flag("todayKills")
}

func (p *PlayerData) YesterdayKills() uint32 {
	return p.yesterdayKills
}

func (p *PlayerData) SetYesterdayKills(val uint32) {
	p.yesterdayKills = val
	p.dirty.Flag("yesterdayKills")
}

func (p *PlayerData) LifetimeHonorableKills() uint32 {
	return p.lifetimeHonorableKills
}

func (p *PlayerData) SetLifetimeHonorableKills(val uint32) {
	p.lifetimeHonorableKills = val
	p.dirty.Flag("lifetimeHonorableKills")
}

func (p *PlayerData) IgnorePowerRegenPredictionMask() uint8 {
	return p.ignorePowerRegenPredictionMask
}

func (p *PlayerData) SetIgnorePowerRegenPredictionMask(val uint8) {
	p.ignorePowerRegenPredictionMask = val
	p.dirty.Flag("ignorePowerRegenPredictionMask")
}

func (p *PlayerData) OverrideSpellsID() uint16 {
	return p.overrideSpellsID
}

func (p *PlayerData) SetOverrideSpellsID(val uint16) {
	p.overrideSpellsID = val
	p.dirty.Flag("overrideSpellsID")
}

func (p *PlayerData) WatchedFactionIndex() int32 {
	return p.watchedFactionIndex
}

func (p *PlayerData) SetWatchedFactionIndex(val int32) {
	p.watchedFactionIndex = val
	p.dirty.Flag("watchedFactionIndex")
}

func (p *PlayerData) CombatRatings() [25]uint32 {
	return p.combatRatings
}

func (p *PlayerData) SetCombatRatings(val [25]uint32) {
	p.combatRatings = val
	p.dirty.Flag("combatRatings")
}

func (p *PlayerData) ArenaTeamInfo() [21]uint32 {
	return p.arenaTeamInfo
}

func (p *PlayerData) SetArenaTeamInfo(val [21]uint32) {
	p.arenaTeamInfo = val
	p.dirty.Flag("arenaTeamInfo")
}

func (p *PlayerData) HonorCurrency() uint32 {
	return p.honorCurrency
}

func (p *PlayerData) SetHonorCurrency(val uint32) {
	p.honorCurrency = val
	p.dirty.Flag("honorCurrency")
}

func (p *PlayerData) ArenaCurrency() uint32 {
	return p.arenaCurrency
}

func (p *PlayerData) SetArenaCurrency(val uint32) {
	p.arenaCurrency = val
	p.dirty.Flag("arenaCurrency")
}

func (p *PlayerData) MaxLevel() uint32 {
	return p.maxLevel
}

func (p *PlayerData) SetMaxLevel(val uint32) {
	p.maxLevel = val
	p.dirty.Flag("maxLevel")
}

func (p *PlayerData) DailyQuests() [25]uint32 {
	return p.dailyQuests
}

func (p *PlayerData) SetDailyQuests(val [25]uint32) {
	p.dailyQuests = val
	p.dirty.Flag("dailyQuests")
}

func (p *PlayerData) RuneRegen() [4]float32 {
	return p.runeRegen
}

func (p *PlayerData) SetRuneRegen(val [4]float32) {
	p.runeRegen = val
	p.dirty.Flag("runeRegen")
}

func (p *PlayerData) NoReagentCost() [3]uint32 {
	return p.noReagentCost
}

func (p *PlayerData) SetNoReagentCost(val [3]uint32) {
	p.noReagentCost = val
	p.dirty.Flag("noReagentCost")
}

func (p *PlayerData) GlyphSlots() [6]uint32 {
	return p.glyphSlots
}

func (p *PlayerData) SetGlyphSlots(val [6]uint32) {
	p.glyphSlots = val
	p.dirty.Flag("glyphSlots")
}

func (p *PlayerData) Glyphs() [6]uint32 {
	return p.glyphs
}

func (p *PlayerData) SetGlyphs(val [6]uint32) {
	p.glyphs = val
	p.dirty.Flag("glyphs")
}

func (p *PlayerData) GlyphsEnabled() uint32 {
	return p.glyphsEnabled
}

func (p *PlayerData) SetGlyphsEnabled(val uint32) {
	p.glyphsEnabled = val
	p.dirty.Flag("glyphsEnabled")
}

func (p *PlayerData) PetSpellPower() uint32 {
	return p.petSpellPower
}

func (p *PlayerData) SetPetSpellPower(val uint32) {
	p.petSpellPower = val
	p.dirty.Flag("petSpellPower")
}
