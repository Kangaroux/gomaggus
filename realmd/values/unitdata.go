package values

import (
	"reflect"

	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
)

const (
	UnitDataOffset = ObjectDataSize
	UnitDataSize   = 142
)

// Adapted from Gophercraft with some modifications
// https://github.com/Gophercraft/core/blob/master/packet/update/d12340/descriptor.go
type UnitData struct {
	charm         realmd.Guid
	summon        realmd.Guid
	critter       realmd.Guid
	charmedBy     realmd.Guid
	summonedBy    realmd.Guid
	createdBy     realmd.Guid
	target        realmd.Guid
	channelObject realmd.Guid
	channelSpell  uint32

	// Bytes0
	race      model.Race
	class     model.Class
	gender    model.Gender
	powerType realmd.PowerType

	health                            uint32
	mana                              uint32
	rage                              uint32
	focus                             uint32
	energy                            uint32
	happiness                         uint32
	runes                             uint32
	runicPower                        uint32
	maxHealth                         uint32
	maxMana                           uint32
	maxRage                           uint32
	maxFocus                          uint32
	maxEnergy                         uint32
	maxHappiness                      uint32
	maxRunes                          uint32
	maxRunicPower                     uint32
	powerRegenFlatModifier            [7]float32
	powerRegenInterruptedFlatModifier [7]float32
	level                             uint32
	faction                           uint32
	virtualItemSlotIDs                [3]uint32

	// Flags 1
	serverControlled    bool // 0x1
	nonAttackable       bool // 0x2
	removeClientControl bool // 0x4
	playerControlled    bool // 0x8
	rename              bool // 0x10
	petAbandon          bool // 0x20
	unk6                bool // 0x40
	unk7                bool // 0x80
	oocNotAttackable    bool // 0x100
	passive             bool // 0x200
	unk10               bool // 0x400
	unk11               bool // 0x800
	pVP                 bool // 0x1000
	isSilenced          bool // 0x2000
	isPersuaded         bool // 0x4000
	swimming            bool // 0x8000
	removeAttackIcon    bool // 0x10000
	isPacified          bool // 0x20000
	isStunned           bool // 0x40000
	inCombat            bool // 0x80000
	inTaxiFlight        bool // 0x100000
	disarmed            bool // 0x200000
	confused            bool // 0x400000
	fleeing             bool // 0x800000
	possessed           bool // 0x1000000
	notSelectable       bool // 0x2000000
	skinnable           bool // 0x4000000
	aurasVisible        bool // 0x8000000
	unk28               bool // 0x10000000
	unk29               bool // 0x20000000
	sheathe             bool // 0x40000000
	noKillReward        bool // 0x80000000

	// Flags 2
	feignDeath          bool // 0x1
	hideBodyArmor       bool // 0x2
	ignoreReputation    bool // 0x4
	comprehendLanguage  bool // 0x8
	cloned              bool // 0x10
	unitFlagUnk5        bool // 0x20
	forceMove           bool // 0x40
	disarmOffhand       bool // 0x80
	unitFlagUnk8        bool // 0x100
	unitFlagUnk9        bool // 0x200
	disarmRanged        bool // 0x400
	regeneratePower     bool // 0x800
	spellClickInGroup   bool // 0x1000
	spellClickDisabled  bool // 0x2000
	interactAnyReaction bool // 0x4000
	_                   [17]bool

	auraState         uint32
	baseAttackTime    uint32
	offhandAttackTime uint32
	rangedAttackTime  uint32
	boundingRadius    float32
	combatReach       float32
	displayID         uint32
	nativeDisplayID   uint32
	mountDisplayID    uint32
	minDamage         float32
	maxDamage         float32
	minOffhandDamage  uint32
	maxOffhandDamage  uint32
	standState        uint8
	loyaltyLevel      uint8
	shapeshiftForm    uint8
	standMiscFlags    uint8
	petNumber         uint32
	petNameTimestamp  uint32
	petExperience     uint32
	petNextLevelExp   uint32

	lootable              bool // 0x1
	trackUnit             bool // 0x2
	tapped                bool // 0x4
	tappedByPlayer        bool // 0x8
	specialInfo           bool // 0x10
	visuallyDead          bool // 0x20
	referAFriend          bool // 0x40
	tappedByAllThreatList bool // 0x80
	_                     [24]bool

	modCastSpeed   float32
	createdBySpell uint32
	gossip         bool // 0x1
	questGiver     bool // 0x2
	vendor         bool // 0x4
	flightMaster   bool // 0x8
	trainer        bool // 0x10
	spiritHealer   bool // 0x20
	spiritGuide    bool // 0x40
	innkeeper      bool // 0x80
	banker         bool // 0x100
	petitioner     bool // 0x200
	tabardDesigner bool // 0x400
	battleMaster   bool // 0x800
	auctioneer     bool // 0x1000
	stableMaster   bool // 0x2000
	repairer       bool // 0x4000
	_              [17]bool

	nPCEmoteState                  uint32
	strength                       uint32
	agility                        uint32
	stamina                        uint32
	intellect                      uint32
	spirit                         uint32
	unitPosStats                   [5]uint32
	unitNegStats                   [5]uint32
	resistances                    [7]uint32
	unitResistanceBuffModsPositive [7]uint32
	unitResistanceBuffModsNegative [7]uint32
	baseMana                       uint32
	baseHealth                     uint32
	sheathState                    uint8
	auraByteFlags                  uint8
	petRename                      uint8
	petShapeshiftForm              uint8
	attackPower                    int32
	attackPowerMods                int32
	attackPowerMultiplier          float32
	rangedAttackPower              int32
	rangedAttackPowerMods          int32
	rangedAttackPowerMultiplier    float32
	minRangedDamage                float32
	maxRangedDamage                float32
	powerCostModifier              [7]uint32
	powerCostMultiplier            [7]float32
	maxHealthModifier              float32
	hoverHeight                    float32
	_                              uint32

	dirty *dirtyValues `value:"END"`
}

func NewUnitData() *UnitData {
	return &UnitData{
		dirty: newDirtyValues(getStructLayout(reflect.ValueOf(UnitData{}))),
	}
}

func (u *UnitData) Marshal(onlyDirty bool) ([]byte, []structSection) {
	return marshalValues(u, onlyDirty, u.dirty)
}

func (u *UnitData) Charm() realmd.Guid {
	return u.charm
}

func (u *UnitData) SetCharm(val realmd.Guid) {
	u.charm = val
	u.dirty.Flag("charm")
}

func (u *UnitData) Summon() realmd.Guid {
	return u.summon
}

func (u *UnitData) SetSummon(val realmd.Guid) {
	u.summon = val
	u.dirty.Flag("summon")
}

func (u *UnitData) Critter() realmd.Guid {
	return u.critter
}

func (u *UnitData) SetCritter(val realmd.Guid) {
	u.critter = val
	u.dirty.Flag("critter")
}

func (u *UnitData) CharmedBy() realmd.Guid {
	return u.charmedBy
}

func (u *UnitData) SetCharmedBy(val realmd.Guid) {
	u.charmedBy = val
	u.dirty.Flag("charmedBy")
}

func (u *UnitData) SummonedBy() realmd.Guid {
	return u.summonedBy
}

func (u *UnitData) SetSummonedBy(val realmd.Guid) {
	u.summonedBy = val
	u.dirty.Flag("summonedBy")
}

func (u *UnitData) CreatedBy() realmd.Guid {
	return u.createdBy
}

func (u *UnitData) SetCreatedBy(val realmd.Guid) {
	u.createdBy = val
	u.dirty.Flag("createdBy")
}

func (u *UnitData) Target() realmd.Guid {
	return u.target
}

func (u *UnitData) SetTarget(val realmd.Guid) {
	u.target = val
	u.dirty.Flag("target")
}

func (u *UnitData) ChannelObject() realmd.Guid {
	return u.channelObject
}

func (u *UnitData) SetChannelObject(val realmd.Guid) {
	u.channelObject = val
	u.dirty.Flag("channelObject")
}

func (u *UnitData) ChannelSpell() uint32 {
	return u.channelSpell
}

func (u *UnitData) SetChannelSpell(val uint32) {
	u.channelSpell = val
	u.dirty.Flag("channelSpell")
}

func (u *UnitData) Race() model.Race {
	return u.race
}

func (u *UnitData) SetRace(val model.Race) {
	u.race = val
	u.dirty.Flag("race")
}

func (u *UnitData) Class() model.Class {
	return u.class
}

func (u *UnitData) SetClass(val model.Class) {
	u.class = val
	u.dirty.Flag("class")
}

func (u *UnitData) Gender() model.Gender {
	return u.gender
}

func (u *UnitData) SetGender(val model.Gender) {
	u.gender = val
	u.dirty.Flag("gender")
}

func (u *UnitData) PowerType() realmd.PowerType {
	return u.powerType
}

func (u *UnitData) SetPowerType(val realmd.PowerType) {
	u.powerType = val
	u.dirty.Flag("powerType")
}

func (u *UnitData) Health() uint32 {
	return u.health
}

func (u *UnitData) SetHealth(val uint32) {
	u.health = val
	u.dirty.Flag("health")
}

func (u *UnitData) Mana() uint32 {
	return u.mana
}

func (u *UnitData) SetMana(val uint32) {
	u.mana = val
	u.dirty.Flag("mana")
}

func (u *UnitData) Rage() uint32 {
	return u.rage
}

func (u *UnitData) SetRage(val uint32) {
	u.rage = val
	u.dirty.Flag("rage")
}

func (u *UnitData) Focus() uint32 {
	return u.focus
}

func (u *UnitData) SetFocus(val uint32) {
	u.focus = val
	u.dirty.Flag("focus")
}

func (u *UnitData) Energy() uint32 {
	return u.energy
}

func (u *UnitData) SetEnergy(val uint32) {
	u.energy = val
	u.dirty.Flag("energy")
}

func (u *UnitData) Happiness() uint32 {
	return u.happiness
}

func (u *UnitData) SetHappiness(val uint32) {
	u.happiness = val
	u.dirty.Flag("happiness")
}

func (u *UnitData) Runes() uint32 {
	return u.runes
}

func (u *UnitData) SetRunes(val uint32) {
	u.runes = val
	u.dirty.Flag("runes")
}

func (u *UnitData) RunicPower() uint32 {
	return u.runicPower
}

func (u *UnitData) SetRunicPower(val uint32) {
	u.runicPower = val
	u.dirty.Flag("runicPower")
}

func (u *UnitData) MaxHealth() uint32 {
	return u.maxHealth
}

func (u *UnitData) SetMaxHealth(val uint32) {
	u.maxHealth = val
	u.dirty.Flag("maxHealth")
}

func (u *UnitData) MaxMana() uint32 {
	return u.maxMana
}

func (u *UnitData) SetMaxMana(val uint32) {
	u.maxMana = val
	u.dirty.Flag("maxMana")
}

func (u *UnitData) MaxRage() uint32 {
	return u.maxRage
}

func (u *UnitData) SetMaxRage(val uint32) {
	u.maxRage = val
	u.dirty.Flag("maxRage")
}

func (u *UnitData) MaxFocus() uint32 {
	return u.maxFocus
}

func (u *UnitData) SetMaxFocus(val uint32) {
	u.maxFocus = val
	u.dirty.Flag("maxFocus")
}

func (u *UnitData) MaxEnergy() uint32 {
	return u.maxEnergy
}

func (u *UnitData) SetMaxEnergy(val uint32) {
	u.maxEnergy = val
	u.dirty.Flag("maxEnergy")
}

func (u *UnitData) MaxHappiness() uint32 {
	return u.maxHappiness
}

func (u *UnitData) SetMaxHappiness(val uint32) {
	u.maxHappiness = val
	u.dirty.Flag("maxHappiness")
}

func (u *UnitData) MaxRunes() uint32 {
	return u.maxRunes
}

func (u *UnitData) SetMaxRunes(val uint32) {
	u.maxRunes = val
	u.dirty.Flag("maxRunes")
}

func (u *UnitData) MaxRunicPower() uint32 {
	return u.maxRunicPower
}

func (u *UnitData) SetMaxRunicPower(val uint32) {
	u.maxRunicPower = val
	u.dirty.Flag("maxRunicPower")
}

func (u *UnitData) PowerRegenFlatModifier() [7]float32 {
	return u.powerRegenFlatModifier
}

func (u *UnitData) SetPowerRegenFlatModifier(val [7]float32) {
	u.powerRegenFlatModifier = val
	u.dirty.Flag("powerRegenFlatModifier")
}

func (u *UnitData) PowerRegenInterruptedFlatModifier() [7]float32 {
	return u.powerRegenInterruptedFlatModifier
}

func (u *UnitData) SetPowerRegenInterruptedFlatModifier(val [7]float32) {
	u.powerRegenInterruptedFlatModifier = val
	u.dirty.Flag("powerRegenInterruptedFlatModifier")
}

func (u *UnitData) Level() uint32 {
	return u.level
}

func (u *UnitData) SetLevel(val uint32) {
	u.level = val
	u.dirty.Flag("level")
}

func (u *UnitData) Faction() uint32 {
	return u.faction
}

func (u *UnitData) SetFaction(val uint32) {
	u.faction = val
	u.dirty.Flag("faction")
}

func (u *UnitData) VirtualItemSlotIDs() [3]uint32 {
	return u.virtualItemSlotIDs
}

func (u *UnitData) SetVirtualItemSlotIDs(val [3]uint32) {
	u.virtualItemSlotIDs = val
	u.dirty.Flag("virtualItemSlotIDs")
}

func (u *UnitData) ServerControlled() bool {
	return u.serverControlled
}

func (u *UnitData) SetServerControlled(val bool) {
	u.serverControlled = val
	u.dirty.Flag("serverControlled")
}

func (u *UnitData) NonAttackable() bool {
	return u.nonAttackable
}

func (u *UnitData) SetNonAttackable(val bool) {
	u.nonAttackable = val
	u.dirty.Flag("nonAttackable")
}

func (u *UnitData) RemoveClientControl() bool {
	return u.removeClientControl
}

func (u *UnitData) SetRemoveClientControl(val bool) {
	u.removeClientControl = val
	u.dirty.Flag("removeClientControl")
}

func (u *UnitData) PlayerControlled() bool {
	return u.playerControlled
}

func (u *UnitData) SetPlayerControlled(val bool) {
	u.playerControlled = val
	u.dirty.Flag("playerControlled")
}

func (u *UnitData) Rename() bool {
	return u.rename
}

func (u *UnitData) SetRename(val bool) {
	u.rename = val
	u.dirty.Flag("rename")
}

func (u *UnitData) PetAbandon() bool {
	return u.petAbandon
}

func (u *UnitData) SetPetAbandon(val bool) {
	u.petAbandon = val
	u.dirty.Flag("petAbandon")
}

func (u *UnitData) Unk6() bool {
	return u.unk6
}

func (u *UnitData) SetUnk6(val bool) {
	u.unk6 = val
	u.dirty.Flag("unk6")
}

func (u *UnitData) Unk7() bool {
	return u.unk7
}

func (u *UnitData) SetUnk7(val bool) {
	u.unk7 = val
	u.dirty.Flag("unk7")
}

func (u *UnitData) OOCNotAttackable() bool {
	return u.oocNotAttackable
}

func (u *UnitData) SetOOCNotAttackable(val bool) {
	u.oocNotAttackable = val
	u.dirty.Flag("oocNotAttackable")
}

func (u *UnitData) Passive() bool {
	return u.passive
}

func (u *UnitData) SetPassive(val bool) {
	u.passive = val
	u.dirty.Flag("passive")
}

func (u *UnitData) Unk10() bool {
	return u.unk10
}

func (u *UnitData) SetUnk10(val bool) {
	u.unk10 = val
	u.dirty.Flag("unk10")
}

func (u *UnitData) Unk11() bool {
	return u.unk11
}

func (u *UnitData) SetUnk11(val bool) {
	u.unk11 = val
	u.dirty.Flag("unk11")
}

func (u *UnitData) PVP() bool {
	return u.pVP
}

func (u *UnitData) SetPVP(val bool) {
	u.pVP = val
	u.dirty.Flag("pVP")
}

func (u *UnitData) IsSilenced() bool {
	return u.isSilenced
}

func (u *UnitData) SetIsSilenced(val bool) {
	u.isSilenced = val
	u.dirty.Flag("isSilenced")
}

func (u *UnitData) IsPersuaded() bool {
	return u.isPersuaded
}

func (u *UnitData) SetIsPersuaded(val bool) {
	u.isPersuaded = val
	u.dirty.Flag("isPersuaded")
}

func (u *UnitData) Swimming() bool {
	return u.swimming
}

func (u *UnitData) SetSwimming(val bool) {
	u.swimming = val
	u.dirty.Flag("swimming")
}

func (u *UnitData) RemoveAttackIcon() bool {
	return u.removeAttackIcon
}

func (u *UnitData) SetRemoveAttackIcon(val bool) {
	u.removeAttackIcon = val
	u.dirty.Flag("removeAttackIcon")
}

func (u *UnitData) IsPacified() bool {
	return u.isPacified
}

func (u *UnitData) SetIsPacified(val bool) {
	u.isPacified = val
	u.dirty.Flag("isPacified")
}

func (u *UnitData) IsStunned() bool {
	return u.isStunned
}

func (u *UnitData) SetIsStunned(val bool) {
	u.isStunned = val
	u.dirty.Flag("isStunned")
}

func (u *UnitData) InCombat() bool {
	return u.inCombat
}

func (u *UnitData) SetInCombat(val bool) {
	u.inCombat = val
	u.dirty.Flag("inCombat")
}

func (u *UnitData) InTaxiFlight() bool {
	return u.inTaxiFlight
}

func (u *UnitData) SetInTaxiFlight(val bool) {
	u.inTaxiFlight = val
	u.dirty.Flag("inTaxiFlight")
}

func (u *UnitData) Disarmed() bool {
	return u.disarmed
}

func (u *UnitData) SetDisarmed(val bool) {
	u.disarmed = val
	u.dirty.Flag("disarmed")
}

func (u *UnitData) Confused() bool {
	return u.confused
}

func (u *UnitData) SetConfused(val bool) {
	u.confused = val
	u.dirty.Flag("confused")
}

func (u *UnitData) Fleeing() bool {
	return u.fleeing
}

func (u *UnitData) SetFleeing(val bool) {
	u.fleeing = val
	u.dirty.Flag("fleeing")
}

func (u *UnitData) Possessed() bool {
	return u.possessed
}

func (u *UnitData) SetPossessed(val bool) {
	u.possessed = val
	u.dirty.Flag("possessed")
}

func (u *UnitData) NotSelectable() bool {
	return u.notSelectable
}

func (u *UnitData) SetNotSelectable(val bool) {
	u.notSelectable = val
	u.dirty.Flag("notSelectable")
}

func (u *UnitData) Skinnable() bool {
	return u.skinnable
}

func (u *UnitData) SetSkinnable(val bool) {
	u.skinnable = val
	u.dirty.Flag("skinnable")
}

func (u *UnitData) AurasVisible() bool {
	return u.aurasVisible
}

func (u *UnitData) SetAurasVisible(val bool) {
	u.aurasVisible = val
	u.dirty.Flag("aurasVisible")
}

func (u *UnitData) Unk28() bool {
	return u.unk28
}

func (u *UnitData) SetUnk28(val bool) {
	u.unk28 = val
	u.dirty.Flag("unk28")
}

func (u *UnitData) Unk29() bool {
	return u.unk29
}

func (u *UnitData) SetUnk29(val bool) {
	u.unk29 = val
	u.dirty.Flag("unk29")
}

func (u *UnitData) Sheathe() bool {
	return u.sheathe
}

func (u *UnitData) SetSheathe(val bool) {
	u.sheathe = val
	u.dirty.Flag("sheathe")
}

func (u *UnitData) NoKillReward() bool {
	return u.noKillReward
}

func (u *UnitData) SetNoKillReward(val bool) {
	u.noKillReward = val
	u.dirty.Flag("noKillReward")
}

func (u *UnitData) FeignDeath() bool {
	return u.feignDeath
}

func (u *UnitData) SetFeignDeath(val bool) {
	u.feignDeath = val
	u.dirty.Flag("feignDeath")
}

func (u *UnitData) HideBodyArmor() bool {
	return u.hideBodyArmor
}

func (u *UnitData) SetHideBodyArmor(val bool) {
	u.hideBodyArmor = val
	u.dirty.Flag("hideBodyArmor")
}

func (u *UnitData) IgnoreReputation() bool {
	return u.ignoreReputation
}

func (u *UnitData) SetIgnoreReputation(val bool) {
	u.ignoreReputation = val
	u.dirty.Flag("ignoreReputation")
}

func (u *UnitData) ComprehendLanguage() bool {
	return u.comprehendLanguage
}

func (u *UnitData) SetComprehendLanguage(val bool) {
	u.comprehendLanguage = val
	u.dirty.Flag("comprehendLanguage")
}

func (u *UnitData) Cloned() bool {
	return u.cloned
}

func (u *UnitData) SetCloned(val bool) {
	u.cloned = val
	u.dirty.Flag("cloned")
}

func (u *UnitData) UnitFlagUnk5() bool {
	return u.unitFlagUnk5
}

func (u *UnitData) SetUnitFlagUnk5(val bool) {
	u.unitFlagUnk5 = val
	u.dirty.Flag("unitFlagUnk5")
}

func (u *UnitData) ForceMove() bool {
	return u.forceMove
}

func (u *UnitData) SetForceMove(val bool) {
	u.forceMove = val
	u.dirty.Flag("forceMove")
}

func (u *UnitData) DisarmOffhand() bool {
	return u.disarmOffhand
}

func (u *UnitData) SetDisarmOffhand(val bool) {
	u.disarmOffhand = val
	u.dirty.Flag("disarmOffhand")
}

func (u *UnitData) UnitFlagUnk8() bool {
	return u.unitFlagUnk8
}

func (u *UnitData) SetUnitFlagUnk8(val bool) {
	u.unitFlagUnk8 = val
	u.dirty.Flag("unitFlagUnk8")
}

func (u *UnitData) UnitFlagUnk9() bool {
	return u.unitFlagUnk9
}

func (u *UnitData) SetUnitFlagUnk9(val bool) {
	u.unitFlagUnk9 = val
	u.dirty.Flag("unitFlagUnk9")
}

func (u *UnitData) DisarmRanged() bool {
	return u.disarmRanged
}

func (u *UnitData) SetDisarmRanged(val bool) {
	u.disarmRanged = val
	u.dirty.Flag("disarmRanged")
}

func (u *UnitData) RegeneratePower() bool {
	return u.regeneratePower
}

func (u *UnitData) SetRegeneratePower(val bool) {
	u.regeneratePower = val
	u.dirty.Flag("regeneratePower")
}

func (u *UnitData) SpellClickInGroup() bool {
	return u.spellClickInGroup
}

func (u *UnitData) SetSpellClickInGroup(val bool) {
	u.spellClickInGroup = val
	u.dirty.Flag("spellClickInGroup")
}

func (u *UnitData) SpellClickDisabled() bool {
	return u.spellClickDisabled
}

func (u *UnitData) SetSpellClickDisabled(val bool) {
	u.spellClickDisabled = val
	u.dirty.Flag("spellClickDisabled")
}

func (u *UnitData) InteractAnyReaction() bool {
	return u.interactAnyReaction
}

func (u *UnitData) SetInteractAnyReaction(val bool) {
	u.interactAnyReaction = val
	u.dirty.Flag("interactAnyReaction")
}

func (u *UnitData) AuraState() uint32 {
	return u.auraState
}

func (u *UnitData) SetAuraState(val uint32) {
	u.auraState = val
	u.dirty.Flag("auraState")
}

func (u *UnitData) BaseAttackTime() uint32 {
	return u.baseAttackTime
}

func (u *UnitData) SetBaseAttackTime(val uint32) {
	u.baseAttackTime = val
	u.dirty.Flag("baseAttackTime")
}

func (u *UnitData) OffhandAttackTime() uint32 {
	return u.offhandAttackTime
}

func (u *UnitData) SetOffhandAttackTime(val uint32) {
	u.offhandAttackTime = val
	u.dirty.Flag("offhandAttackTime")
}

func (u *UnitData) RangedAttackTime() uint32 {
	return u.rangedAttackTime
}

func (u *UnitData) SetRangedAttackTime(val uint32) {
	u.rangedAttackTime = val
	u.dirty.Flag("rangedAttackTime")
}

func (u *UnitData) BoundingRadius() float32 {
	return u.boundingRadius
}

func (u *UnitData) SetBoundingRadius(val float32) {
	u.boundingRadius = val
	u.dirty.Flag("boundingRadius")
}

func (u *UnitData) CombatReach() float32 {
	return u.combatReach
}

func (u *UnitData) SetCombatReach(val float32) {
	u.combatReach = val
	u.dirty.Flag("combatReach")
}

func (u *UnitData) DisplayID() uint32 {
	return u.displayID
}

func (u *UnitData) SetDisplayID(val uint32) {
	u.displayID = val
	u.dirty.Flag("displayID")
}

func (u *UnitData) NativeDisplayID() uint32 {
	return u.nativeDisplayID
}

func (u *UnitData) SetNativeDisplayID(val uint32) {
	u.nativeDisplayID = val
	u.dirty.Flag("nativeDisplayID")
}

func (u *UnitData) MountDisplayID() uint32 {
	return u.mountDisplayID
}

func (u *UnitData) SetMountDisplayID(val uint32) {
	u.mountDisplayID = val
	u.dirty.Flag("mountDisplayID")
}

func (u *UnitData) MinDamage() float32 {
	return u.minDamage
}

func (u *UnitData) SetMinDamage(val float32) {
	u.minDamage = val
	u.dirty.Flag("minDamage")
}

func (u *UnitData) MaxDamage() float32 {
	return u.maxDamage
}

func (u *UnitData) SetMaxDamage(val float32) {
	u.maxDamage = val
	u.dirty.Flag("maxDamage")
}

func (u *UnitData) MinOffhandDamage() uint32 {
	return u.minOffhandDamage
}

func (u *UnitData) SetMinOffhandDamage(val uint32) {
	u.minOffhandDamage = val
	u.dirty.Flag("minOffhandDamage")
}

func (u *UnitData) MaxOffhandDamage() uint32 {
	return u.maxOffhandDamage
}

func (u *UnitData) SetMaxOffhandDamage(val uint32) {
	u.maxOffhandDamage = val
	u.dirty.Flag("maxOffhandDamage")
}

func (u *UnitData) StandState() uint8 {
	return u.standState
}

func (u *UnitData) SetStandState(val uint8) {
	u.standState = val
	u.dirty.Flag("standState")
}

func (u *UnitData) LoyaltyLevel() uint8 {
	return u.loyaltyLevel
}

func (u *UnitData) SetLoyaltyLevel(val uint8) {
	u.loyaltyLevel = val
	u.dirty.Flag("loyaltyLevel")
}

func (u *UnitData) ShapeshiftForm() uint8 {
	return u.shapeshiftForm
}

func (u *UnitData) SetShapeshiftForm(val uint8) {
	u.shapeshiftForm = val
	u.dirty.Flag("shapeshiftForm")
}

func (u *UnitData) StandMiscFlags() uint8 {
	return u.standMiscFlags
}

func (u *UnitData) SetStandMiscFlags(val uint8) {
	u.standMiscFlags = val
	u.dirty.Flag("standMiscFlags")
}

func (u *UnitData) PetNumber() uint32 {
	return u.petNumber
}

func (u *UnitData) SetPetNumber(val uint32) {
	u.petNumber = val
	u.dirty.Flag("petNumber")
}

func (u *UnitData) PetNameTimestamp() uint32 {
	return u.petNameTimestamp
}

func (u *UnitData) SetPetNameTimestamp(val uint32) {
	u.petNameTimestamp = val
	u.dirty.Flag("petNameTimestamp")
}

func (u *UnitData) PetExperience() uint32 {
	return u.petExperience
}

func (u *UnitData) SetPetExperience(val uint32) {
	u.petExperience = val
	u.dirty.Flag("petExperience")
}

func (u *UnitData) PetNextLevelExp() uint32 {
	return u.petNextLevelExp
}

func (u *UnitData) SetPetNextLevelExp(val uint32) {
	u.petNextLevelExp = val
	u.dirty.Flag("petNextLevelExp")
}

func (u *UnitData) Lootable() bool {
	return u.lootable
}

func (u *UnitData) SetLootable(val bool) {
	u.lootable = val
	u.dirty.Flag("lootable")
}

func (u *UnitData) TrackUnit() bool {
	return u.trackUnit
}

func (u *UnitData) SetTrackUnit(val bool) {
	u.trackUnit = val
	u.dirty.Flag("trackUnit")
}

func (u *UnitData) Tapped() bool {
	return u.tapped
}

func (u *UnitData) SetTapped(val bool) {
	u.tapped = val
	u.dirty.Flag("tapped")
}

func (u *UnitData) TappedByPlayer() bool {
	return u.tappedByPlayer
}

func (u *UnitData) SetTappedByPlayer(val bool) {
	u.tappedByPlayer = val
	u.dirty.Flag("tappedByPlayer")
}

func (u *UnitData) SpecialInfo() bool {
	return u.specialInfo
}

func (u *UnitData) SetSpecialInfo(val bool) {
	u.specialInfo = val
	u.dirty.Flag("specialInfo")
}

func (u *UnitData) VisuallyDead() bool {
	return u.visuallyDead
}

func (u *UnitData) SetVisuallyDead(val bool) {
	u.visuallyDead = val
	u.dirty.Flag("visuallyDead")
}

func (u *UnitData) ReferAFriend() bool {
	return u.referAFriend
}

func (u *UnitData) SetReferAFriend(val bool) {
	u.referAFriend = val
	u.dirty.Flag("referAFriend")
}

func (u *UnitData) TappedByAllThreatList() bool {
	return u.tappedByAllThreatList
}

func (u *UnitData) SetTappedByAllThreatList(val bool) {
	u.tappedByAllThreatList = val
	u.dirty.Flag("tappedByAllThreatList")
}

func (u *UnitData) ModCastSpeed() float32 {
	return u.modCastSpeed
}

func (u *UnitData) SetModCastSpeed(val float32) {
	u.modCastSpeed = val
	u.dirty.Flag("modCastSpeed")
}

func (u *UnitData) CreatedBySpell() uint32 {
	return u.createdBySpell
}

func (u *UnitData) SetCreatedBySpell(val uint32) {
	u.createdBySpell = val
	u.dirty.Flag("createdBySpell")
}

func (u *UnitData) Gossip() bool {
	return u.gossip
}

func (u *UnitData) SetGossip(val bool) {
	u.gossip = val
	u.dirty.Flag("gossip")
}

func (u *UnitData) QuestGiver() bool {
	return u.questGiver
}

func (u *UnitData) SetQuestGiver(val bool) {
	u.questGiver = val
	u.dirty.Flag("questGiver")
}

func (u *UnitData) Vendor() bool {
	return u.vendor
}

func (u *UnitData) SetVendor(val bool) {
	u.vendor = val
	u.dirty.Flag("vendor")
}

func (u *UnitData) FlightMaster() bool {
	return u.flightMaster
}

func (u *UnitData) SetFlightMaster(val bool) {
	u.flightMaster = val
	u.dirty.Flag("flightMaster")
}

func (u *UnitData) Trainer() bool {
	return u.trainer
}

func (u *UnitData) SetTrainer(val bool) {
	u.trainer = val
	u.dirty.Flag("trainer")
}

func (u *UnitData) SpiritHealer() bool {
	return u.spiritHealer
}

func (u *UnitData) SetSpiritHealer(val bool) {
	u.spiritHealer = val
	u.dirty.Flag("spiritHealer")
}

func (u *UnitData) SpiritGuide() bool {
	return u.spiritGuide
}

func (u *UnitData) SetSpiritGuide(val bool) {
	u.spiritGuide = val
	u.dirty.Flag("spiritGuide")
}

func (u *UnitData) Innkeeper() bool {
	return u.innkeeper
}

func (u *UnitData) SetInnkeeper(val bool) {
	u.innkeeper = val
	u.dirty.Flag("innkeeper")
}

func (u *UnitData) Banker() bool {
	return u.banker
}

func (u *UnitData) SetBanker(val bool) {
	u.banker = val
	u.dirty.Flag("banker")
}

func (u *UnitData) Petitioner() bool {
	return u.petitioner
}

func (u *UnitData) SetPetitioner(val bool) {
	u.petitioner = val
	u.dirty.Flag("petitioner")
}

func (u *UnitData) TabardDesigner() bool {
	return u.tabardDesigner
}

func (u *UnitData) SetTabardDesigner(val bool) {
	u.tabardDesigner = val
	u.dirty.Flag("tabardDesigner")
}

func (u *UnitData) BattleMaster() bool {
	return u.battleMaster
}

func (u *UnitData) SetBattleMaster(val bool) {
	u.battleMaster = val
	u.dirty.Flag("battleMaster")
}

func (u *UnitData) Auctioneer() bool {
	return u.auctioneer
}

func (u *UnitData) SetAuctioneer(val bool) {
	u.auctioneer = val
	u.dirty.Flag("auctioneer")
}

func (u *UnitData) StableMaster() bool {
	return u.stableMaster
}

func (u *UnitData) SetStableMaster(val bool) {
	u.stableMaster = val
	u.dirty.Flag("stableMaster")
}

func (u *UnitData) Repairer() bool {
	return u.repairer
}

func (u *UnitData) SetRepairer(val bool) {
	u.repairer = val
	u.dirty.Flag("repairer")
}

func (u *UnitData) NPCEmoteState() uint32 {
	return u.nPCEmoteState
}

func (u *UnitData) SetNPCEmoteState(val uint32) {
	u.nPCEmoteState = val
	u.dirty.Flag("nPCEmoteState")
}

func (u *UnitData) Strength() uint32 {
	return u.strength
}

func (u *UnitData) SetStrength(val uint32) {
	u.strength = val
	u.dirty.Flag("strength")
}

func (u *UnitData) Agility() uint32 {
	return u.agility
}

func (u *UnitData) SetAgility(val uint32) {
	u.agility = val
	u.dirty.Flag("agility")
}

func (u *UnitData) Stamina() uint32 {
	return u.stamina
}

func (u *UnitData) SetStamina(val uint32) {
	u.stamina = val
	u.dirty.Flag("stamina")
}

func (u *UnitData) Intellect() uint32 {
	return u.intellect
}

func (u *UnitData) SetIntellect(val uint32) {
	u.intellect = val
	u.dirty.Flag("intellect")
}

func (u *UnitData) Spirit() uint32 {
	return u.spirit
}

func (u *UnitData) SetSpirit(val uint32) {
	u.spirit = val
	u.dirty.Flag("spirit")
}

func (u *UnitData) UnitPosStats() [5]uint32 {
	return u.unitPosStats
}

func (u *UnitData) SetUnitPosStats(val [5]uint32) {
	u.unitPosStats = val
	u.dirty.Flag("unitPosStats")
}

func (u *UnitData) UnitNegStats() [5]uint32 {
	return u.unitNegStats
}

func (u *UnitData) SetUnitNegStats(val [5]uint32) {
	u.unitNegStats = val
	u.dirty.Flag("unitNegStats")
}

func (u *UnitData) Resistances() [7]uint32 {
	return u.resistances
}

func (u *UnitData) SetResistances(val [7]uint32) {
	u.resistances = val
	u.dirty.Flag("resistances")
}

func (u *UnitData) UnitResistanceBuffModsPositive() [7]uint32 {
	return u.unitResistanceBuffModsPositive
}

func (u *UnitData) SetUnitResistanceBuffModsPositive(val [7]uint32) {
	u.unitResistanceBuffModsPositive = val
	u.dirty.Flag("unitResistanceBuffModsPositive")
}

func (u *UnitData) UnitResistanceBuffModsNegative() [7]uint32 {
	return u.unitResistanceBuffModsNegative
}

func (u *UnitData) SetUnitResistanceBuffModsNegative(val [7]uint32) {
	u.unitResistanceBuffModsNegative = val
	u.dirty.Flag("unitResistanceBuffModsNegative")
}

func (u *UnitData) BaseMana() uint32 {
	return u.baseMana
}

func (u *UnitData) SetBaseMana(val uint32) {
	u.baseMana = val
	u.dirty.Flag("baseMana")
}

func (u *UnitData) BaseHealth() uint32 {
	return u.baseHealth
}

func (u *UnitData) SetBaseHealth(val uint32) {
	u.baseHealth = val
	u.dirty.Flag("baseHealth")
}

func (u *UnitData) SheathState() uint8 {
	return u.sheathState
}

func (u *UnitData) SetSheathState(val uint8) {
	u.sheathState = val
	u.dirty.Flag("sheathState")
}

func (u *UnitData) AuraByteFlags() uint8 {
	return u.auraByteFlags
}

func (u *UnitData) SetAuraByteFlags(val uint8) {
	u.auraByteFlags = val
	u.dirty.Flag("auraByteFlags")
}

func (u *UnitData) PetRename() uint8 {
	return u.petRename
}

func (u *UnitData) SetPetRename(val uint8) {
	u.petRename = val
	u.dirty.Flag("petRename")
}

func (u *UnitData) PetShapeshiftForm() uint8 {
	return u.petShapeshiftForm
}

func (u *UnitData) SetPetShapeshiftForm(val uint8) {
	u.petShapeshiftForm = val
	u.dirty.Flag("petShapeshiftForm")
}

func (u *UnitData) AttackPower() int32 {
	return u.attackPower
}

func (u *UnitData) SetAttackPower(val int32) {
	u.attackPower = val
	u.dirty.Flag("attackPower")
}

func (u *UnitData) AttackPowerMods() int32 {
	return u.attackPowerMods
}

func (u *UnitData) SetAttackPowerMods(val int32) {
	u.attackPowerMods = val
	u.dirty.Flag("attackPowerMods")
}

func (u *UnitData) AttackPowerMultiplier() float32 {
	return u.attackPowerMultiplier
}

func (u *UnitData) SetAttackPowerMultiplier(val float32) {
	u.attackPowerMultiplier = val
	u.dirty.Flag("attackPowerMultiplier")
}

func (u *UnitData) RangedAttackPower() int32 {
	return u.rangedAttackPower
}

func (u *UnitData) SetRangedAttackPower(val int32) {
	u.rangedAttackPower = val
	u.dirty.Flag("rangedAttackPower")
}

func (u *UnitData) RangedAttackPowerMods() int32 {
	return u.rangedAttackPowerMods
}

func (u *UnitData) SetRangedAttackPowerMods(val int32) {
	u.rangedAttackPowerMods = val
	u.dirty.Flag("rangedAttackPowerMods")
}

func (u *UnitData) RangedAttackPowerMultiplier() float32 {
	return u.rangedAttackPowerMultiplier
}

func (u *UnitData) SetRangedAttackPowerMultiplier(val float32) {
	u.rangedAttackPowerMultiplier = val
	u.dirty.Flag("rangedAttackPowerMultiplier")
}

func (u *UnitData) MinRangedDamage() float32 {
	return u.minRangedDamage
}

func (u *UnitData) SetMinRangedDamage(val float32) {
	u.minRangedDamage = val
	u.dirty.Flag("minRangedDamage")
}

func (u *UnitData) MaxRangedDamage() float32 {
	return u.maxRangedDamage
}

func (u *UnitData) SetMaxRangedDamage(val float32) {
	u.maxRangedDamage = val
	u.dirty.Flag("maxRangedDamage")
}

func (u *UnitData) PowerCostModifier() [7]uint32 {
	return u.powerCostModifier
}

func (u *UnitData) SetPowerCostModifier(val [7]uint32) {
	u.powerCostModifier = val
	u.dirty.Flag("powerCostModifier")
}

func (u *UnitData) PowerCostMultiplier() [7]float32 {
	return u.powerCostMultiplier
}

func (u *UnitData) SetPowerCostMultiplier(val [7]float32) {
	u.powerCostMultiplier = val
	u.dirty.Flag("powerCostMultiplier")
}

func (u *UnitData) MaxHealthModifier() float32 {
	return u.maxHealthModifier
}

func (u *UnitData) SetMaxHealthModifier(val float32) {
	u.maxHealthModifier = val
	u.dirty.Flag("maxHealthModifier")
}

func (u *UnitData) HoverHeight() float32 {
	return u.hoverHeight
}

func (u *UnitData) SetHoverHeight(val float32) {
	u.hoverHeight = val
	u.dirty.Flag("hoverHeight")
}
