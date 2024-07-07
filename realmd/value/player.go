package value

import (
	"reflect"

	"github.com/kangaroux/gomaggus/realmd"
)

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

	questLog [25]questLogEntry

	visibleItems [19]visibleItem

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
	skills                    [128]skillEntry
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

func NewPlayer() *Player {
	return &Player{
		dirty: newDirtyValues(getStructLayout(reflect.ValueOf(Player{}))),
	}
}

func (p *Player) Marshal(onlyDirty bool) []byte {
	return marshalValues(p, onlyDirty, p.dirty)
}

func (p *Player) DuelArbiter() realmd.Guid {
	return p.duelArbiter
}

func (p *Player) SetDuelArbiter(val realmd.Guid) {
	p.duelArbiter = val
	p.dirty.Flag("duelArbiter")
}

func (p *Player) GroupLeader() bool {
	return p.groupLeader
}

func (p *Player) SetGroupLeader(val bool) {
	p.groupLeader = val
	p.dirty.Flag("groupLeader")
}

func (p *Player) AFK() bool {
	return p.afk
}

func (p *Player) SetAFK(val bool) {
	p.afk = val
	p.dirty.Flag("afk")
}

func (p *Player) DND() bool {
	return p.dnd
}

func (p *Player) SetDND(val bool) {
	p.dnd = val
	p.dirty.Flag("dnd")
}

func (p *Player) Gm() bool {
	return p.gm
}

func (p *Player) SetGM(val bool) {
	p.gm = val
	p.dirty.Flag("gm")
}

func (p *Player) Ghost() bool {
	return p.ghost
}

func (p *Player) SetGhost(val bool) {
	p.ghost = val
	p.dirty.Flag("ghost")
}

func (p *Player) Resting() bool {
	return p.resting
}

func (p *Player) SetResting(val bool) {
	p.resting = val
	p.dirty.Flag("resting")
}

func (p *Player) VoiceChat() bool {
	return p.voiceChat
}

func (p *Player) SetVoiceChat(val bool) {
	p.voiceChat = val
	p.dirty.Flag("voiceChat")
}

func (p *Player) FFAPVP() bool {
	return p.ffapvp
}

func (p *Player) SetFFAPVP(val bool) {
	p.ffapvp = val
	p.dirty.Flag("ffapvp")
}

func (p *Player) ContestedPVP() bool {
	return p.contestedPVP
}

func (p *Player) SetContestedPVP(val bool) {
	p.contestedPVP = val
	p.dirty.Flag("contestedPVP")
}

func (p *Player) InPVP() bool {
	return p.inPVP
}

func (p *Player) SetInPVP(val bool) {
	p.inPVP = val
	p.dirty.Flag("inPVP")
}

func (p *Player) HideHelm() bool {
	return p.hideHelm
}

func (p *Player) SetHideHelm(val bool) {
	p.hideHelm = val
	p.dirty.Flag("hideHelm")
}

func (p *Player) HideCloak() bool {
	return p.hideCloak
}

func (p *Player) SetHideCloak(val bool) {
	p.hideCloak = val
	p.dirty.Flag("hideCloak")
}

func (p *Player) PlayedLongTime() bool {
	return p.playedLongTime
}

func (p *Player) SetPlayedLongTime(val bool) {
	p.playedLongTime = val
	p.dirty.Flag("playedLongTime")
}

func (p *Player) PlayedTooLong() bool {
	return p.playedTooLong
}

func (p *Player) SetPlayedTooLong(val bool) {
	p.playedTooLong = val
	p.dirty.Flag("playedTooLong")
}

func (p *Player) OutOfBounds() bool {
	return p.outOfBounds
}

func (p *Player) SetOutOfBounds(val bool) {
	p.outOfBounds = val
	p.dirty.Flag("outOfBounds")
}

func (p *Player) Developer() bool {
	return p.developer
}

func (p *Player) SetDeveloper(val bool) {
	p.developer = val
	p.dirty.Flag("developer")
}

func (p *Player) TaxiBenchmark() bool {
	return p.taxiBenchmark
}

func (p *Player) SetTaxiBenchmark(val bool) {
	p.taxiBenchmark = val
	p.dirty.Flag("taxiBenchmark")
}

func (p *Player) PVPTimer() bool {
	return p.pvpTimer
}

func (p *Player) SetPVPTimer(val bool) {
	p.pvpTimer = val
	p.dirty.Flag("pvpTimer")
}

func (p *Player) Uber() bool {
	return p.uber
}

func (p *Player) SetUber(val bool) {
	p.uber = val
	p.dirty.Flag("uber")
}

func (p *Player) Commentator() bool {
	return p.commentator
}

func (p *Player) SetCommentator(val bool) {
	p.commentator = val
	p.dirty.Flag("commentator")
}

func (p *Player) OnlyAllowAbilities() bool {
	return p.onlyAllowAbilities
}

func (p *Player) SetOnlyAllowAbilities(val bool) {
	p.onlyAllowAbilities = val
	p.dirty.Flag("onlyAllowAbilities")
}

func (p *Player) StopMeleeOnTab() bool {
	return p.stopMeleeOnTab
}

func (p *Player) SetStopMeleeOnTab(val bool) {
	p.stopMeleeOnTab = val
	p.dirty.Flag("stopMeleeOnTab")
}

func (p *Player) NoExperienceGain() bool {
	return p.noExperienceGain
}

func (p *Player) SetNoExperienceGain(val bool) {
	p.noExperienceGain = val
	p.dirty.Flag("noExperienceGain")
}

func (p *Player) GuildID() uint32 {
	return p.guildID
}

func (p *Player) SetGuildID(val uint32) {
	p.guildID = val
	p.dirty.Flag("guildID")
}

func (p *Player) GuildRank() uint32 {
	return p.guildRank
}

func (p *Player) SetGuildRank(val uint32) {
	p.guildRank = val
	p.dirty.Flag("guildRank")
}

func (p *Player) Skin() uint8 {
	return p.skin
}

func (p *Player) SetSkin(val uint8) {
	p.skin = val
	p.dirty.Flag("skin")
}

func (p *Player) Face() uint8 {
	return p.face
}

func (p *Player) SetFace(val uint8) {
	p.face = val
	p.dirty.Flag("face")
}

func (p *Player) HairStyle() uint8 {
	return p.hairStyle
}

func (p *Player) SetHairStyle(val uint8) {
	p.hairStyle = val
	p.dirty.Flag("hairStyle")
}

func (p *Player) HairColor() uint8 {
	return p.hairColor
}

func (p *Player) SetHairColor(val uint8) {
	p.hairColor = val
	p.dirty.Flag("hairColor")
}

func (p *Player) FacialHair() uint8 {
	return p.facialHair
}

func (p *Player) SetFacialHair(val uint8) {
	p.facialHair = val
	p.dirty.Flag("facialHair")
}

func (p *Player) RestBits() uint8 {
	return p.restBits
}

func (p *Player) SetRestBits(val uint8) {
	p.restBits = val
	p.dirty.Flag("restBits")
}

func (p *Player) BankBagSlotCount() uint8 {
	return p.bankBagSlotCount
}

func (p *Player) SetBankBagSlotCount(val uint8) {
	p.bankBagSlotCount = val
	p.dirty.Flag("bankBagSlotCount")
}

func (p *Player) RestState() uint8 {
	return p.restState
}

func (p *Player) SetRestState(val uint8) {
	p.restState = val
	p.dirty.Flag("restState")
}

func (p *Player) PlayerGender() uint8 {
	return p.playerGender
}

func (p *Player) SetPlayerGender(val uint8) {
	p.playerGender = val
	p.dirty.Flag("playerGender")
}

func (p *Player) GenderUnk() uint8 {
	return p.genderUnk
}

func (p *Player) SetGenderUnk(val uint8) {
	p.genderUnk = val
	p.dirty.Flag("genderUnk")
}

func (p *Player) Drunkness() uint8 {
	return p.drunkness
}

func (p *Player) SetDrunkness(val uint8) {
	p.drunkness = val
	p.dirty.Flag("drunkness")
}

func (p *Player) PVPRank() uint8 {
	return p.pVPRank
}

func (p *Player) SetPVPRank(val uint8) {
	p.pVPRank = val
	p.dirty.Flag("pVPRank")
}

func (p *Player) DuelTeam() uint32 {
	return p.duelTeam
}

func (p *Player) SetDuelTeam(val uint32) {
	p.duelTeam = val
	p.dirty.Flag("duelTeam")
}

func (p *Player) GuildTimestamp() uint32 {
	return p.guildTimestamp
}

func (p *Player) SetGuildTimestamp(val uint32) {
	p.guildTimestamp = val
	p.dirty.Flag("guildTimestamp")
}

func (p *Player) QuestLog() [25]questLogEntry {
	return p.questLog
}

func (p *Player) SetQuestLog(val [25]questLogEntry) {
	p.questLog = val
	p.dirty.Flag("questLog")
}

func (p *Player) VisibleItems() [19]visibleItem {
	return p.visibleItems
}

func (p *Player) SetVisibleItems(val [19]visibleItem) {
	p.visibleItems = val
	p.dirty.Flag("visibleItems")
}

func (p *Player) ChosenTitle() uint32 {
	return p.chosenTitle
}

func (p *Player) SetChosenTitle(val uint32) {
	p.chosenTitle = val
	p.dirty.Flag("chosenTitle")
}

func (p *Player) FakeInebriation() uint32 {
	return p.fakeInebriation
}

func (p *Player) SetFakeInebriation(val uint32) {
	p.fakeInebriation = val
	p.dirty.Flag("fakeInebriation")
}

func (p *Player) InventorySlots() [23]uint64 {
	return p.inventorySlots
}

func (p *Player) SetInventorySlots(val [23]uint64) {
	p.inventorySlots = val
	p.dirty.Flag("inventorySlots")
}

func (p *Player) PackSlots() [16]uint64 {
	return p.packSlots
}

func (p *Player) SetPackSlots(val [16]uint64) {
	p.packSlots = val
	p.dirty.Flag("packSlots")
}

func (p *Player) BankSlots() [28]uint64 {
	return p.bankSlots
}

func (p *Player) SetBankSlots(val [28]uint64) {
	p.bankSlots = val
	p.dirty.Flag("bankSlots")
}

func (p *Player) BankBagSlots() [7]uint64 {
	return p.bankBagSlots
}

func (p *Player) SetBankBagSlots(val [7]uint64) {
	p.bankBagSlots = val
	p.dirty.Flag("bankBagSlots")
}

func (p *Player) VendorBuybackSlots() [12]uint64 {
	return p.vendorBuybackSlots
}

func (p *Player) SetVendorBuybackSlots(val [12]uint64) {
	p.vendorBuybackSlots = val
	p.dirty.Flag("vendorBuybackSlots")
}

func (p *Player) KeyringSlots() [32]uint64 {
	return p.keyringSlots
}

func (p *Player) SetKeyringSlots(val [32]uint64) {
	p.keyringSlots = val
	p.dirty.Flag("keyringSlots")
}

func (p *Player) CurrencyTokenSlots() [32]uint64 {
	return p.currencyTokenSlots
}

func (p *Player) SetCurrencyTokenSlots(val [32]uint64) {
	p.currencyTokenSlots = val
	p.dirty.Flag("currencyTokenSlots")
}

func (p *Player) FarSight() uint64 {
	return p.farSight
}

func (p *Player) SetFarSight(val uint64) {
	p.farSight = val
	p.dirty.Flag("farSight")
}

func (p *Player) KnownTitles() [6]uint32 {
	return p.knownTitles
}

func (p *Player) SetKnownTitles(val [6]uint32) {
	p.knownTitles = val
	p.dirty.Flag("knownTitles")
}

func (p *Player) KnownCurrencies() [2]uint32 {
	return p.knownCurrencies
}

func (p *Player) SetKnownCurrencies(val [2]uint32) {
	p.knownCurrencies = val
	p.dirty.Flag("knownCurrencies")
}

func (p *Player) Xp() uint32 {
	return p.xp
}

func (p *Player) SetXp(val uint32) {
	p.xp = val
	p.dirty.Flag("xp")
}

func (p *Player) NextLevelXP() uint32 {
	return p.nextLevelXP
}

func (p *Player) SetNextLevelXP(val uint32) {
	p.nextLevelXP = val
	p.dirty.Flag("nextLevelXP")
}

func (p *Player) Skills() [128]skillEntry {
	return p.skills
}

func (p *Player) SetSkills(val [128]skillEntry) {
	p.skills = val
	p.dirty.Flag("skills")
}

func (p *Player) CharacterPoints() [2]uint32 {
	return p.characterPoints
}

func (p *Player) SetCharacterPoints(val [2]uint32) {
	p.characterPoints = val
	p.dirty.Flag("characterPoints")
}

func (p *Player) TrackCreatures() uint32 {
	return p.trackCreatures
}

func (p *Player) SetTrackCreatures(val uint32) {
	p.trackCreatures = val
	p.dirty.Flag("trackCreatures")
}

func (p *Player) TrackResources() uint32 {
	return p.trackResources
}

func (p *Player) SetTrackResources(val uint32) {
	p.trackResources = val
	p.dirty.Flag("trackResources")
}

func (p *Player) BlockPercentage() float32 {
	return p.blockPercentage
}

func (p *Player) SetBlockPercentage(val float32) {
	p.blockPercentage = val
	p.dirty.Flag("blockPercentage")
}

func (p *Player) DodgePercentage() float32 {
	return p.dodgePercentage
}

func (p *Player) SetDodgePercentage(val float32) {
	p.dodgePercentage = val
	p.dirty.Flag("dodgePercentage")
}

func (p *Player) ParryPercentage() float32 {
	return p.parryPercentage
}

func (p *Player) SetParryPercentage(val float32) {
	p.parryPercentage = val
	p.dirty.Flag("parryPercentage")
}

func (p *Player) Expertise() uint32 {
	return p.expertise
}

func (p *Player) SetExpertise(val uint32) {
	p.expertise = val
	p.dirty.Flag("expertise")
}

func (p *Player) OffhandExpertise() uint32 {
	return p.offhandExpertise
}

func (p *Player) SetOffhandExpertise(val uint32) {
	p.offhandExpertise = val
	p.dirty.Flag("offhandExpertise")
}

func (p *Player) CritPercentage() float32 {
	return p.critPercentage
}

func (p *Player) SetCritPercentage(val float32) {
	p.critPercentage = val
	p.dirty.Flag("critPercentage")
}

func (p *Player) RangedCritPercentage() float32 {
	return p.rangedCritPercentage
}

func (p *Player) SetRangedCritPercentage(val float32) {
	p.rangedCritPercentage = val
	p.dirty.Flag("rangedCritPercentage")
}

func (p *Player) OffhandCritPercentage() float32 {
	return p.offhandCritPercentage
}

func (p *Player) SetOffhandCritPercentage(val float32) {
	p.offhandCritPercentage = val
	p.dirty.Flag("offhandCritPercentage")
}

func (p *Player) SpellCritPercentage() [7]float32 {
	return p.spellCritPercentage
}

func (p *Player) SetSpellCritPercentage(val [7]float32) {
	p.spellCritPercentage = val
	p.dirty.Flag("spellCritPercentage")
}

func (p *Player) ShieldBlock() uint32 {
	return p.shieldBlock
}

func (p *Player) SetShieldBlock(val uint32) {
	p.shieldBlock = val
	p.dirty.Flag("shieldBlock")
}

func (p *Player) ShieldBlockCritPercentage() float32 {
	return p.shieldBlockCritPercentage
}

func (p *Player) SetShieldBlockCritPercentage(val float32) {
	p.shieldBlockCritPercentage = val
	p.dirty.Flag("shieldBlockCritPercentage")
}

func (p *Player) ExploredZones() [128]uint32 {
	return p.exploredZones
}

func (p *Player) SetExploredZones(val [128]uint32) {
	p.exploredZones = val
	p.dirty.Flag("exploredZones")
}

func (p *Player) RestStateExperience() uint32 {
	return p.restStateExperience
}

func (p *Player) SetRestStateExperience(val uint32) {
	p.restStateExperience = val
	p.dirty.Flag("restStateExperience")
}

func (p *Player) Wealth() int32 {
	return p.wealth
}

func (p *Player) SetWealth(val int32) {
	p.wealth = val
	p.dirty.Flag("wealth")
}

func (p *Player) ModDamageDonePositive() [7]uint32 {
	return p.modDamageDonePositive
}

func (p *Player) SetModDamageDonePositive(val [7]uint32) {
	p.modDamageDonePositive = val
	p.dirty.Flag("modDamageDonePositive")
}

func (p *Player) ModDamageDoneNegative() [7]uint32 {
	return p.modDamageDoneNegative
}

func (p *Player) SetModDamageDoneNegative(val [7]uint32) {
	p.modDamageDoneNegative = val
	p.dirty.Flag("modDamageDoneNegative")
}

func (p *Player) ModDamageDonePercentage() [7]float32 {
	return p.modDamageDonePercentage
}

func (p *Player) SetModDamageDonePercentage(val [7]float32) {
	p.modDamageDonePercentage = val
	p.dirty.Flag("modDamageDonePercentage")
}

func (p *Player) ModHealingDonePos() uint32 {
	return p.modHealingDonePos
}

func (p *Player) SetModHealingDonePos(val uint32) {
	p.modHealingDonePos = val
	p.dirty.Flag("modHealingDonePos")
}

func (p *Player) ModHealingPercentage() float32 {
	return p.modHealingPercentage
}

func (p *Player) SetModHealingPercentage(val float32) {
	p.modHealingPercentage = val
	p.dirty.Flag("modHealingPercentage")
}

func (p *Player) ModHealingDonePercentage() float32 {
	return p.modHealingDonePercentage
}

func (p *Player) SetModHealingDonePercentage(val float32) {
	p.modHealingDonePercentage = val
	p.dirty.Flag("modHealingDonePercentage")
}

func (p *Player) ModTarResistance() uint32 {
	return p.modTarResistance
}

func (p *Player) SetModTarResistance(val uint32) {
	p.modTarResistance = val
	p.dirty.Flag("modTarResistance")
}

func (p *Player) ModTarPhysicalResistance() uint32 {
	return p.modTarPhysicalResistance
}

func (p *Player) SetModTarPhysicalResistance(val uint32) {
	p.modTarPhysicalResistance = val
	p.dirty.Flag("modTarPhysicalResistance")
}

func (p *Player) TrackStealthed() bool {
	return p.trackStealthed
}

func (p *Player) SetTrackStealthed(val bool) {
	p.trackStealthed = val
	p.dirty.Flag("trackStealthed")
}

func (p *Player) DisplaySpiritAutoReleaSetImer() bool {
	return p.displaySpiritAutoReleaSetImer
}

func (p *Player) SetDisplaySpiritAutoReleaSetImer(val bool) {
	p.displaySpiritAutoReleaSetImer = val
	p.dirty.Flag("displaySpiritAutoReleaSetImer")
}

func (p *Player) HideSpiritReleaseWindow() bool {
	return p.hideSpiritReleaseWindow
}

func (p *Player) SetHideSpiritReleaseWindow(val bool) {
	p.hideSpiritReleaseWindow = val
	p.dirty.Flag("hideSpiritReleaseWindow")
}

func (p *Player) ReferAFriendGrantableLevel() uint8 {
	return p.referAFriendGrantableLevel
}

func (p *Player) SetReferAFriendGrantableLevel(val uint8) {
	p.referAFriendGrantableLevel = val
	p.dirty.Flag("referAFriendGrantableLevel")
}

func (p *Player) ActionBarToggles() uint8 {
	return p.actionBarToggles
}

func (p *Player) SetActionBarToggles(val uint8) {
	p.actionBarToggles = val
	p.dirty.Flag("actionBarToggles")
}

func (p *Player) LifetimeMaxPVPRank() uint8 {
	return p.lifetimeMaxPVPRank
}

func (p *Player) SetLifetimeMaxPVPRank(val uint8) {
	p.lifetimeMaxPVPRank = val
	p.dirty.Flag("lifetimeMaxPVPRank")
}

func (p *Player) AmmoID() uint32 {
	return p.ammoID
}

func (p *Player) SetAmmoID(val uint32) {
	p.ammoID = val
	p.dirty.Flag("ammoID")
}

func (p *Player) SelfResSpell() uint32 {
	return p.selfResSpell
}

func (p *Player) SetSelfResSpell(val uint32) {
	p.selfResSpell = val
	p.dirty.Flag("selfResSpell")
}

func (p *Player) PVPMedals() uint32 {
	return p.pVPMedals
}

func (p *Player) SetPVPMedals(val uint32) {
	p.pVPMedals = val
	p.dirty.Flag("pVPMedals")
}

func (p *Player) BuybackPrices() [12]uint32 {
	return p.buybackPrices
}

func (p *Player) SetBuybackPrices(val [12]uint32) {
	p.buybackPrices = val
	p.dirty.Flag("buybackPrices")
}

func (p *Player) BuybackTimestamps() [12]uint32 {
	return p.buybackTimestamps
}

func (p *Player) SetBuybackTimestamps(val [12]uint32) {
	p.buybackTimestamps = val
	p.dirty.Flag("buybackTimestamps")
}

func (p *Player) Kills() uint32 {
	return p.kills
}

func (p *Player) SetKills(val uint32) {
	p.kills = val
	p.dirty.Flag("kills")
}

func (p *Player) TodayKills() uint32 {
	return p.todayKills
}

func (p *Player) SetTodayKills(val uint32) {
	p.todayKills = val
	p.dirty.Flag("todayKills")
}

func (p *Player) YesterdayKills() uint32 {
	return p.yesterdayKills
}

func (p *Player) SetYesterdayKills(val uint32) {
	p.yesterdayKills = val
	p.dirty.Flag("yesterdayKills")
}

func (p *Player) LifetimeHonorableKills() uint32 {
	return p.lifetimeHonorableKills
}

func (p *Player) SetLifetimeHonorableKills(val uint32) {
	p.lifetimeHonorableKills = val
	p.dirty.Flag("lifetimeHonorableKills")
}

func (p *Player) IgnorePowerRegenPredictionMask() uint8 {
	return p.ignorePowerRegenPredictionMask
}

func (p *Player) SetIgnorePowerRegenPredictionMask(val uint8) {
	p.ignorePowerRegenPredictionMask = val
	p.dirty.Flag("ignorePowerRegenPredictionMask")
}

func (p *Player) OverrideSpellsID() uint16 {
	return p.overrideSpellsID
}

func (p *Player) SetOverrideSpellsID(val uint16) {
	p.overrideSpellsID = val
	p.dirty.Flag("overrideSpellsID")
}

func (p *Player) WatchedFactionIndex() int32 {
	return p.watchedFactionIndex
}

func (p *Player) SetWatchedFactionIndex(val int32) {
	p.watchedFactionIndex = val
	p.dirty.Flag("watchedFactionIndex")
}

func (p *Player) CombatRatings() [25]uint32 {
	return p.combatRatings
}

func (p *Player) SetCombatRatings(val [25]uint32) {
	p.combatRatings = val
	p.dirty.Flag("combatRatings")
}

func (p *Player) ArenaTeamInfo() [21]uint32 {
	return p.arenaTeamInfo
}

func (p *Player) SetArenaTeamInfo(val [21]uint32) {
	p.arenaTeamInfo = val
	p.dirty.Flag("arenaTeamInfo")
}

func (p *Player) HonorCurrency() uint32 {
	return p.honorCurrency
}

func (p *Player) SetHonorCurrency(val uint32) {
	p.honorCurrency = val
	p.dirty.Flag("honorCurrency")
}

func (p *Player) ArenaCurrency() uint32 {
	return p.arenaCurrency
}

func (p *Player) SetArenaCurrency(val uint32) {
	p.arenaCurrency = val
	p.dirty.Flag("arenaCurrency")
}

func (p *Player) MaxLevel() uint32 {
	return p.maxLevel
}

func (p *Player) SetMaxLevel(val uint32) {
	p.maxLevel = val
	p.dirty.Flag("maxLevel")
}

func (p *Player) DailyQuests() [25]uint32 {
	return p.dailyQuests
}

func (p *Player) SetDailyQuests(val [25]uint32) {
	p.dailyQuests = val
	p.dirty.Flag("dailyQuests")
}

func (p *Player) RuneRegen() [4]float32 {
	return p.runeRegen
}

func (p *Player) SetRuneRegen(val [4]float32) {
	p.runeRegen = val
	p.dirty.Flag("runeRegen")
}

func (p *Player) NoReagentCost() [3]uint32 {
	return p.noReagentCost
}

func (p *Player) SetNoReagentCost(val [3]uint32) {
	p.noReagentCost = val
	p.dirty.Flag("noReagentCost")
}

func (p *Player) GlyphSlots() [6]uint32 {
	return p.glyphSlots
}

func (p *Player) SetGlyphSlots(val [6]uint32) {
	p.glyphSlots = val
	p.dirty.Flag("glyphSlots")
}

func (p *Player) Glyphs() [6]uint32 {
	return p.glyphs
}

func (p *Player) SetGlyphs(val [6]uint32) {
	p.glyphs = val
	p.dirty.Flag("glyphs")
}

func (p *Player) GlyphsEnabled() uint32 {
	return p.glyphsEnabled
}

func (p *Player) SetGlyphsEnabled(val uint32) {
	p.glyphsEnabled = val
	p.dirty.Flag("glyphsEnabled")
}

func (p *Player) PetSpellPower() uint32 {
	return p.petSpellPower
}

func (p *Player) SetPetSpellPower(val uint32) {
	p.petSpellPower = val
	p.dirty.Flag("petSpellPower")
}
