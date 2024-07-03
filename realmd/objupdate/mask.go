package objupdate

import (
	"bytes"
	"encoding/binary"
)

// https://gtker.com/wow_messages/types/update-mask.html#version-335
type FieldMask struct {
	// Size is the number of uint32 blocks used for the field data.
	Size uint32
	// Offset is the bit number that is set to indicate this field is included.
	Offset uint32
	Name   string
	// Length is used for array fields, and describes how many values the array can hold.
	// Each value is Size bytes long. Length=0 is interpreted as Length=1.
	Length int
}

func (fm *FieldMask) String() string {
	return fm.Name
}

var (
	FieldMaskObjectGuid    = FieldMask{Size: 2, Offset: 0x0, Name: "FieldMaskObjectGuid"}
	FieldMaskObjectType    = FieldMask{Size: 1, Offset: 0x2, Name: "FieldMaskObjectType"}
	FieldMaskObjectEntry   = FieldMask{Size: 1, Offset: 0x3, Name: "FieldMaskObjectEntry"}
	FieldMaskObjectScaleX  = FieldMask{Size: 1, Offset: 0x4, Name: "FieldMaskObjectScaleX"}
	FieldMaskObjectPadding = FieldMask{Size: 1, Offset: 0x5, Name: "FieldMaskObjectPadding"}

	FieldMaskItemOwner              = FieldMask{Size: 2, Offset: 0x6, Name: "FieldMaskItemOwner"}
	FieldMaskItemContained          = FieldMask{Size: 2, Offset: 0x8, Name: "FieldMaskItemContained"}
	FieldMaskItemCreator            = FieldMask{Size: 2, Offset: 0xA, Name: "FieldMaskItemCreator"}
	FieldMaskItemGiftcreator        = FieldMask{Size: 2, Offset: 0xC, Name: "FieldMaskItemGiftcreator"}
	FieldMaskItemStackCount         = FieldMask{Size: 1, Offset: 0xE, Name: "FieldMaskItemStackCount"}
	FieldMaskItemDuration           = FieldMask{Size: 1, Offset: 0xF, Name: "FieldMaskItemDuration"}
	FieldMaskItemSpellCharges       = FieldMask{Size: 1, Offset: 0x10, Length: 5, Name: "FieldMaskItemSpellCharges"}
	FieldMaskItemFlags              = FieldMask{Size: 1, Offset: 0x15, Name: "FieldMaskItemFlags"}
	FieldMaskItemEnchantment        = FieldMask{Size: 3, Offset: 0x16, Length: 12, Name: "FieldMaskItemEnchantment"}
	FieldMaskItemPropertySeed       = FieldMask{Size: 1, Offset: 0x3A, Name: "FieldMaskItemPropertySeed"}
	FieldMaskItemRandomPropertiesId = FieldMask{Size: 1, Offset: 0x3B, Name: "FieldMaskItemRandomPropertiesId"}
	FieldMaskItemDurability         = FieldMask{Size: 1, Offset: 0x3C, Name: "FieldMaskItemDurability"}
	FieldMaskItemMaxdurability      = FieldMask{Size: 1, Offset: 0x3D, Name: "FieldMaskItemMaxdurability"}
	FieldMaskItemCreatePlayedTime   = FieldMask{Size: 1, Offset: 0x3E, Name: "FieldMaskItemCreatePlayedTime"}

	FieldMaskContainerNumSlots = FieldMask{Size: 1, Offset: 0x40, Name: "FieldMaskContainerNumSlots"}
	FieldMaskContainerSlot1    = FieldMask{Size: 2, Offset: 0x42, Length: 36, Name: "FieldMaskContainerSlot"}

	FieldMaskUnitCharm                             = FieldMask{Size: 2, Offset: 0x6, Name: "FieldMaskUnitCharm"}
	FieldMaskUnitSummon                            = FieldMask{Size: 2, Offset: 0x8, Name: "FieldMaskUnitSummon"}
	FieldMaskUnitCritter                           = FieldMask{Size: 2, Offset: 0xA, Name: "FieldMaskUnitCritter"}
	FieldMaskUnitCharmedby                         = FieldMask{Size: 2, Offset: 0xC, Name: "FieldMaskUnitCharmedby"}
	FieldMaskUnitSummonedby                        = FieldMask{Size: 2, Offset: 0xE, Name: "FieldMaskUnitSummonedby"}
	FieldMaskUnitCreatedby                         = FieldMask{Size: 2, Offset: 0x10, Name: "FieldMaskUnitCreatedby"}
	FieldMaskUnitTarget                            = FieldMask{Size: 2, Offset: 0x12, Name: "FieldMaskUnitTarget"}
	FieldMaskUnitChannelObject                     = FieldMask{Size: 2, Offset: 0x14, Name: "FieldMaskUnitChannelObject"}
	FieldMaskUnitChannelSpell                      = FieldMask{Size: 1, Offset: 0x16, Name: "FieldMaskUnitChannelSpell"}
	FieldMaskUnitRaceClassGenderPower              = FieldMask{Size: 1, Offset: 0x17, Name: "FieldMaskUnitRaceClassGenderPower"} // Bytes0
	FieldMaskUnitHealth                            = FieldMask{Size: 1, Offset: 0x18, Name: "FieldMaskUnitHealth"}
	FieldMaskUnitPower                             = FieldMask{Size: 1, Offset: 0x19, Length: 7, Name: "FieldMaskUnitPower"}
	FieldMaskUnitMaxHealth                         = FieldMask{Size: 1, Offset: 0x20, Name: "FieldMaskUnitMaxHealth"}
	FieldMaskUnitMaxPower                          = FieldMask{Size: 1, Offset: 0x21, Length: 7, Name: "FieldMaskUnitMaxPower"}
	FieldMaskUnitPowerRegenFlatModifier            = FieldMask{Size: 1, Offset: 0x28, Length: 7, Name: "FieldMaskUnitPowerRegenFlatModifier"}
	FieldMaskUnitPowerRegenInterruptedFlatModifier = FieldMask{Size: 1, Offset: 0x2F, Length: 7, Name: "FieldMaskUnitPowerRegenInterruptedFlatModifier"}
	FieldMaskUnitLevel                             = FieldMask{Size: 1, Offset: 0x36, Name: "FieldMaskUnitLevel"}
	FieldMaskUnitFactionTemplate                   = FieldMask{Size: 1, Offset: 0x37, Name: "FieldMaskUnitFactionTemplate"}
	FieldMaskUnitVirtualItemSlotId                 = FieldMask{Size: 1, Offset: 0x38, Length: 3, Name: "FieldMaskUnitVirtualItemSlotId"}
	FieldMaskUnitFlags                             = FieldMask{Size: 1, Offset: 0x3B, Name: "FieldMaskUnitFlags"}
	FieldMaskUnitFlags2                            = FieldMask{Size: 1, Offset: 0x3C, Name: "FieldMaskUnitFlags2"}
	FieldMaskUnitAurastate                         = FieldMask{Size: 1, Offset: 0x3D, Name: "FieldMaskUnitAurastate"}
	FieldMaskUnitBaseattacktime                    = FieldMask{Size: 2, Offset: 0x3E, Name: "FieldMaskUnitBaseattacktime"}
	FieldMaskUnitRangedattacktime                  = FieldMask{Size: 1, Offset: 0x40, Name: "FieldMaskUnitRangedattacktime"}
	FieldMaskUnitBoundingradius                    = FieldMask{Size: 1, Offset: 0x41, Name: "FieldMaskUnitBoundingradius"}
	FieldMaskUnitCombatreach                       = FieldMask{Size: 1, Offset: 0x42, Name: "FieldMaskUnitCombatreach"}
	FieldMaskUnitDisplayId                         = FieldMask{Size: 1, Offset: 0x43, Name: "FieldMaskUnitDisplayId"}
	FieldMaskUnitNativeDisplayId                   = FieldMask{Size: 1, Offset: 0x44, Name: "FieldMaskUnitNativeDisplayId"}
	FieldMaskUnitMountDisplayid                    = FieldMask{Size: 1, Offset: 0x45, Name: "FieldMaskUnitMountDisplayid"}
	FieldMaskUnitMinDamage                         = FieldMask{Size: 1, Offset: 0x46, Name: "FieldMaskUnitMinDamage"}
	FieldMaskUnitMaxDamage                         = FieldMask{Size: 1, Offset: 0x47, Name: "FieldMaskUnitMaxDamage"}
	FieldMaskUnitMinOffhandDamage                  = FieldMask{Size: 1, Offset: 0x48, Name: "FieldMaskUnitMinOffhandDamage"}
	FieldMaskUnitMaxOffhandDamage                  = FieldMask{Size: 1, Offset: 0x49, Name: "FieldMaskUnitMaxOffhandDamage"}
	FieldMaskUnitBytes1                            = FieldMask{Size: 1, Offset: 0x4A, Name: "FieldMaskUnitBytes1"}
	FieldMaskUnitPetNumber                         = FieldMask{Size: 1, Offset: 0x4B, Name: "FieldMaskUnitPetNumber"}
	FieldMaskUnitPetNameTimestamp                  = FieldMask{Size: 1, Offset: 0x4C, Name: "FieldMaskUnitPetNameTimestamp"}
	FieldMaskUnitPetExperience                     = FieldMask{Size: 1, Offset: 0x4D, Name: "FieldMaskUnitPetExperience"}
	FieldMaskUnitPetExperienceToNextLevel          = FieldMask{Size: 1, Offset: 0x4E, Name: "FieldMaskUnitPetExperienceToNextLevel"}
	FieldMaskUnitDynamicFlags                      = FieldMask{Size: 1, Offset: 0x4F, Name: "FieldMaskUnitDynamicFlags"}
	FieldMaskUnitModCastSpeed                      = FieldMask{Size: 1, Offset: 0x50, Name: "FieldMaskUnitModCastSpeed"}
	FieldMaskUnitCreatedBySpell                    = FieldMask{Size: 1, Offset: 0x51, Name: "FieldMaskUnitCreatedBySpell"}
	FieldMaskUnitNpcFlags                          = FieldMask{Size: 1, Offset: 0x52, Name: "FieldMaskUnitNpcFlags"}
	FieldMaskUnitNpcEmotestate                     = FieldMask{Size: 1, Offset: 0x53, Name: "FieldMaskUnitNpcEmotestate"}
	FieldMaskUnitStrength                          = FieldMask{Size: 1, Offset: 0x54, Name: "FieldMaskUnitStrength"}
	FieldMaskUnitAgility                           = FieldMask{Size: 1, Offset: 0x55, Name: "FieldMaskUnitAgility"}
	FieldMaskUnitStamina                           = FieldMask{Size: 1, Offset: 0x56, Name: "FieldMaskUnitStamina"}
	FieldMaskUnitIntellect                         = FieldMask{Size: 1, Offset: 0x57, Name: "FieldMaskUnitIntellect"}
	FieldMaskUnitSpirit                            = FieldMask{Size: 1, Offset: 0x58, Name: "FieldMaskUnitSpirit"}
	FieldMaskUnitPosStat                           = FieldMask{Size: 1, Offset: 0x59, Length: 5, Name: "FieldMaskUnitPosStat"}
	FieldMaskUnitNegStat                           = FieldMask{Size: 1, Offset: 0x5E, Length: 5, Name: "FieldMaskUnitNegStat"}
	FieldMaskUnitResistances                       = FieldMask{Size: 1, Offset: 0x63, Length: 7, Name: "FieldMaskUnitResistances"}
	FieldMaskUnitResistanceBuffModsPositive        = FieldMask{Size: 1, Offset: 0x6A, Length: 7, Name: "FieldMaskUnitResistanceBuffModsPositive"}
	FieldMaskUnitResistanceBuffModsNegative        = FieldMask{Size: 1, Offset: 0x71, Length: 7, Name: "FieldMaskUnitResistanceBuffModsNegative"}
	FieldMaskUnitBaseMana                          = FieldMask{Size: 1, Offset: 0x78, Name: "FieldMaskUnitBaseMana"}
	FieldMaskUnitBaseHealth                        = FieldMask{Size: 1, Offset: 0x79, Name: "FieldMaskUnitBaseHealth"}
	FieldMaskUnitBytes2                            = FieldMask{Size: 1, Offset: 0x7A, Name: "FieldMaskUnitBytes2"}
	FieldMaskUnitAttackPower                       = FieldMask{Size: 1, Offset: 0x7B, Name: "FieldMaskUnitAttackPower"}
	FieldMaskUnitAttackPowerMods                   = FieldMask{Size: 1, Offset: 0x7C, Name: "FieldMaskUnitAttackPowerMods"}
	FieldMaskUnitAttackPowerMultiplier             = FieldMask{Size: 1, Offset: 0x7D, Name: "FieldMaskUnitAttackPowerMultiplier"}
	FieldMaskUnitRangedAttackPower                 = FieldMask{Size: 1, Offset: 0x7E, Name: "FieldMaskUnitRangedAttackPower"}
	FieldMaskUnitRangedAttackPowerMods             = FieldMask{Size: 1, Offset: 0x7F, Name: "FieldMaskUnitRangedAttackPowerMods"}
	FieldMaskUnitRangedAttackPowerMultiplier       = FieldMask{Size: 1, Offset: 0x80, Name: "FieldMaskUnitRangedAttackPowerMultiplier"}
	FieldMaskUnitMinRangedDamage                   = FieldMask{Size: 1, Offset: 0x81, Name: "FieldMaskUnitMinRangedDamage"}
	FieldMaskUnitMaxRangedDamage                   = FieldMask{Size: 1, Offset: 0x82, Name: "FieldMaskUnitMaxRangedDamage"}
	FieldMaskUnitPowerCostModifier                 = FieldMask{Size: 1, Offset: 0x83, Length: 7, Name: "FieldMaskUnitPowerCostModifier"}
	FieldMaskUnitPowerCostMultiplier               = FieldMask{Size: 1, Offset: 0x8A, Length: 7, Name: "FieldMaskUnitPowerCostMultiplier"}
	FieldMaskUnitMaxHealthModifier                 = FieldMask{Size: 1, Offset: 0x91, Name: "FieldMaskUnitMaxHealthModifier"}
	FieldMaskUnitHoverHeight                       = FieldMask{Size: 1, Offset: 0x92, Name: "FieldMaskUnitHoverHeight"}

	FieldMaskPlayerDuelArbiter                 = FieldMask{Size: 2, Offset: 0x94, Name: "FieldMaskPlayerDuelArbiter"}
	FieldMaskPlayerFlags                       = FieldMask{Size: 1, Offset: 0x96, Name: "FieldMaskPlayerFlags"}
	FieldMaskPlayerGuildid                     = FieldMask{Size: 1, Offset: 0x97, Name: "FieldMaskPlayerGuildid"}
	FieldMaskPlayerGuildrank                   = FieldMask{Size: 1, Offset: 0x98, Name: "FieldMaskPlayerGuildrank"}
	FieldMaskPlayerBytes1                      = FieldMask{Size: 1, Offset: 0x99, Name: "FieldMaskPlayerFieldBytes"}
	FieldMaskPlayerBytes2                      = FieldMask{Size: 1, Offset: 0x9A, Name: "FieldMaskPlayerBytes2"}
	FieldMaskPlayerBytes3                      = FieldMask{Size: 1, Offset: 0x9B, Name: "FieldMaskPlayerBytes3"}
	FieldMaskPlayerDuelTeam                    = FieldMask{Size: 1, Offset: 0x9C, Name: "FieldMaskPlayerDuelTeam"}
	FieldMaskPlayerGuildTimestamp              = FieldMask{Size: 1, Offset: 0x9D, Name: "FieldMaskPlayerGuildTimestamp"}
	FieldMaskPlayerQuestLog                    = FieldMask{Size: 5, Offset: 0x9E, Length: 25, Name: "FieldMaskPlayerQuestLog"}
	FieldMaskPlayerVisibleItem                 = FieldMask{Size: 2, Offset: 0x11B, Length: 19, Name: "FieldMaskPlayerVisibleItem"}
	FieldMaskPlayerChosenTitle                 = FieldMask{Size: 1, Offset: 0x141, Name: "FieldMaskPlayerChosenTitle"}
	FieldMaskPlayerFakeInebriation             = FieldMask{Size: 1, Offset: 0x142, Name: "FieldMaskPlayerFakeInebriation"}
	FieldMaskPlayerInventorySlots              = FieldMask{Size: 3, Offset: 0x144, Name: "FieldMaskPlayerInventorySlots"}
	FieldMaskPlayerBankSlots                   = FieldMask{Size: 3, Offset: 0x144, Name: "FieldMaskPlayerBankSlots"}
	FieldMaskPlayerBankBags                    = FieldMask{Size: 3, Offset: 0x144, Name: "FieldMaskPlayerBankBags"}
	FieldMaskPlayerBuybackSlots                = FieldMask{Size: 3, Offset: 0x144, Name: "FieldMaskPlayerBuybackSlots"}
	FieldMaskPlayerKeyRing                     = FieldMask{Size: 3, Offset: 0x144, Name: "FieldMaskPlayerKeyRing"}
	FieldMaskPlayerCurrencyTokenSlots          = FieldMask{Size: 3, Offset: 0x144, Name: "FieldMaskPlayerCurrencyTokenSlots"}
	FieldMaskPlayerFarsight                    = FieldMask{Size: 2, Offset: 0x270, Name: "FieldMaskPlayerFarsight"}
	FieldMaskPlayerKnownTitles                 = FieldMask{Size: 2, Offset: 0x272, Length: 3, Name: "FieldMaskPlayerKnownTitles"}
	FieldMaskPlayerKnownCurrencies             = FieldMask{Size: 1, Offset: 0x278, Length: 2, Name: "FieldMaskPlayerKnownCurrencies"}
	FieldMaskPlayerExperience                  = FieldMask{Size: 1, Offset: 0x27A, Name: "FieldMaskPlayerXp"}
	FieldMaskPlayerExperienceToNextLevel       = FieldMask{Size: 1, Offset: 0x27B, Name: "FieldMaskPlayerNextLevelXp"}
	FieldMaskPlayerSkillInfoArray              = FieldMask{Size: 3, Offset: 0x27C, Length: 128, Name: "FieldMaskPlayerSkillInfo"}
	FieldMaskPlayerCharacterPoints             = FieldMask{Size: 1, Offset: 0x3FC, Length: 2, Name: "FieldMaskPlayerCharacterPoints"}
	FieldMaskPlayerTrackCreatures              = FieldMask{Size: 1, Offset: 0x3FE, Name: "FieldMaskPlayerTrackCreatures"}
	FieldMaskPlayerTrackResources              = FieldMask{Size: 1, Offset: 0x3FF, Name: "FieldMaskPlayerTrackResources"}
	FieldMaskPlayerBlockPercentage             = FieldMask{Size: 1, Offset: 0x400, Name: "FieldMaskPlayerBlockPercentage"}
	FieldMaskPlayerDodgePercentage             = FieldMask{Size: 1, Offset: 0x401, Name: "FieldMaskPlayerDodgePercentage"}
	FieldMaskPlayerParryPercentage             = FieldMask{Size: 1, Offset: 0x402, Name: "FieldMaskPlayerParryPercentage"}
	FieldMaskPlayerExpertise                   = FieldMask{Size: 1, Offset: 0x403, Name: "FieldMaskPlayerExpertise"}
	FieldMaskPlayerOffhandExpertise            = FieldMask{Size: 1, Offset: 0x404, Name: "FieldMaskPlayerOffhandExpertise"}
	FieldMaskPlayerCritPercentage              = FieldMask{Size: 1, Offset: 0x405, Name: "FieldMaskPlayerCritPercentage"}
	FieldMaskPlayerRangedCritPercentage        = FieldMask{Size: 1, Offset: 0x406, Name: "FieldMaskPlayerRangedCritPercentage"}
	FieldMaskPlayerOffhandCritPercentage       = FieldMask{Size: 1, Offset: 0x407, Name: "FieldMaskPlayerOffhandCritPercentage"}
	FieldMaskPlayerSpellCritPercentage         = FieldMask{Size: 1, Offset: 0x408, Length: 7, Name: "FieldMaskPlayerSpellCritPercentage"}
	FieldMaskPlayerShieldBlock                 = FieldMask{Size: 1, Offset: 0x40F, Name: "FieldMaskPlayerShieldBlock"}
	FieldMaskPlayerShieldBlockCritPercentage   = FieldMask{Size: 1, Offset: 0x410, Name: "FieldMaskPlayerShieldBlockCritPercentage"}
	FieldMaskPlayerExploredZones               = FieldMask{Size: 1, Offset: 0x411, Length: 128, Name: "FieldMaskPlayerExploredZones"}
	FieldMaskPlayerRestStateExperience         = FieldMask{Size: 1, Offset: 0x491, Name: "FieldMaskPlayerRestStateExperience"}
	FieldMaskPlayerWealth                      = FieldMask{Size: 1, Offset: 0x492, Name: "FieldMaskPlayerCoinage"}
	FieldMaskPlayerModDamageDonePos            = FieldMask{Size: 1, Offset: 0x493, Length: 7, Name: "FieldMaskPlayerModDamageDonePos"}
	FieldMaskPlayerModDamageDoneNeg            = FieldMask{Size: 1, Offset: 0x49A, Length: 7, Name: "FieldMaskPlayerModDamageDoneNeg"}
	FieldMaskPlayerModDamageDonePercentage     = FieldMask{Size: 1, Offset: 0x4A1, Length: 7, Name: "FieldMaskPlayerModDamageDonePercentage"}
	FieldMaskPlayerModHealingDonePos           = FieldMask{Size: 1, Offset: 0x4A8, Name: "FieldMaskPlayerModHealingDonePos"}
	FieldMaskPlayerModHealingPercentage        = FieldMask{Size: 1, Offset: 0x4A9, Name: "FieldMaskPlayerModHealingPercentage"}
	FieldMaskPlayerModHealingDonePercentage    = FieldMask{Size: 1, Offset: 0x4AA, Name: "FieldMaskPlayerModHealingDonePercentage"}
	FieldMaskPlayerModTargetResistance         = FieldMask{Size: 1, Offset: 0x4AB, Name: "FieldMaskPlayerModTargetResistance"}
	FieldMaskPlayerModTargetPhysicalResistance = FieldMask{Size: 1, Offset: 0x4AC, Name: "FieldMaskPlayerModTargetPhysicalResistance"}
	FieldMaskPlayerFeatures                    = FieldMask{Size: 1, Offset: 0x4AD, Name: "FieldMaskPlayerFeatures"}
	FieldMaskPlayerAmmoId                      = FieldMask{Size: 1, Offset: 0x4AE, Name: "FieldMaskPlayerAmmoId"}
	FieldMaskPlayerSelfResSpell                = FieldMask{Size: 1, Offset: 0x4AF, Name: "FieldMaskPlayerSelfResSpell"}
	FieldMaskPlayerPvpMedals                   = FieldMask{Size: 1, Offset: 0x4B0, Name: "FieldMaskPlayerPvpMedals"}
	FieldMaskPlayerBuybackPrice                = FieldMask{Size: 1, Offset: 0x4B1, Length: 12, Name: "FieldMaskPlayerBuybackPrice"}
	FieldMaskPlayerBuybackTimestamp            = FieldMask{Size: 1, Offset: 0x4BD, Length: 12, Name: "FieldMaskPlayerBuybackTimestamp"}
	FieldMaskPlayerKills                       = FieldMask{Size: 1, Offset: 0x4C9, Name: "FieldMaskPlayerKills"}
	FieldMaskPlayerTodayContribution           = FieldMask{Size: 1, Offset: 0x4CA, Name: "FieldMaskPlayerTodayContribution"}
	FieldMaskPlayerYesterdayContribution       = FieldMask{Size: 1, Offset: 0x4CB, Name: "FieldMaskPlayerYesterdayContribution"}
	FieldMaskPlayerLifetimeHonorableKills      = FieldMask{Size: 1, Offset: 0x4CC, Name: "FieldMaskPlayerLifetimeHonorableKills"}
	FieldMaskPlayerBytes4                      = FieldMask{Size: 1, Offset: 0x4CD, Name: "FieldMaskPlayerBytes4"}
	FieldMaskPlayerWatchedFactionIndex         = FieldMask{Size: 1, Offset: 0x4CE, Name: "FieldMaskPlayerWatchedFactionIndex"}
	FieldMaskPlayerCombatRating                = FieldMask{Size: 1, Offset: 0x4CF, Length: 25, Name: "FieldMaskPlayerCombatRating"}
	FieldMaskPlayerArenaTeamInfo               = FieldMask{Size: 1, Offset: 0x4E8, Length: 21, Name: "FieldMaskPlayerArenaTeamInfo"}
	FieldMaskPlayerHonorCurrency               = FieldMask{Size: 1, Offset: 0x4FD, Name: "FieldMaskPlayerHonorCurrency"}
	FieldMaskPlayerArenaCurrency               = FieldMask{Size: 1, Offset: 0x4FE, Name: "FieldMaskPlayerArenaCurrency"}
	FieldMaskPlayerMaxLevel                    = FieldMask{Size: 1, Offset: 0x4FF, Name: "FieldMaskPlayerMaxLevel"}
	FieldMaskPlayerDailyQuests                 = FieldMask{Size: 1, Offset: 0x500, Length: 25, Name: "FieldMaskPlayerDailyQuests"}
	FieldMaskPlayerRuneRegen                   = FieldMask{Size: 1, Offset: 0x519, Length: 4, Name: "FieldMaskPlayerRuneRegen"}
	FieldMaskPlayerNoReagentCost               = FieldMask{Size: 1, Offset: 0x51D, Length: 3, Name: "FieldMaskPlayerNoReagentCost"}
	FieldMaskPlayerGlyphSlots                  = FieldMask{Size: 1, Offset: 0x520, Length: 6, Name: "FieldMaskPlayerGlyphSlots"}
	FieldMaskPlayerGlyphs                      = FieldMask{Size: 1, Offset: 0x526, Length: 6, Name: "FieldMaskPlayerGlyphs"}
	FieldMaskPlayerGlyphsEnabled               = FieldMask{Size: 1, Offset: 0x52C, Name: "FieldMaskPlayerGlyphsEnabled"}
	FieldMaskPlayerPetSpellPower               = FieldMask{Size: 1, Offset: 0x52D, Name: "FieldMaskPlayerPetSpellPower"}

	FieldMaskGameObjectDisplayid      = FieldMask{Size: 1, Offset: 0x8, Name: "FieldMaskGameObjectDisplayid"}
	FieldMaskGameObjectFlags          = FieldMask{Size: 1, Offset: 0x9, Name: "FieldMaskGameObjectFlags"}
	FieldMaskGameObjectParentrotation = FieldMask{Size: 4, Offset: 0xA, Name: "FieldMaskGameObjectParentrotation"}
	FieldMaskGameObjectDynamic        = FieldMask{Size: 1, Offset: 0xE, Name: "FieldMaskGameObjectDynamic"}
	FieldMaskGameObjectFaction        = FieldMask{Size: 1, Offset: 0xF, Name: "FieldMaskGameObjectFaction"}
	FieldMaskGameObjectLevel          = FieldMask{Size: 1, Offset: 0x10, Name: "FieldMaskGameObjectLevel"}
	FieldMaskGameObjectBytes1         = FieldMask{Size: 1, Offset: 0x11, Name: "FieldMaskGameObjectBytes1"}

	FieldMaskDynamicObjectCaster   = FieldMask{Size: 2, Offset: 0x6, Name: "FieldMaskDynamicObjectCaster"}
	FieldMaskDynamicObjectBytes    = FieldMask{Size: 1, Offset: 0x8, Name: "FieldMaskDynamicObjectBytes"}
	FieldMaskDynamicObjectSpellid  = FieldMask{Size: 1, Offset: 0x9, Name: "FieldMaskDynamicObjectSpellid"}
	FieldMaskDynamicObjectRadius   = FieldMask{Size: 1, Offset: 0xA, Name: "FieldMaskDynamicObjectRadius"}
	FieldMaskDynamicObjectCasttime = FieldMask{Size: 1, Offset: 0xB, Name: "FieldMaskDynamicObjectCasttime"}

	FieldMaskCorpseOwner        = FieldMask{Size: 2, Offset: 0x6, Name: "FieldMaskCorpseOwner"}
	FieldMaskCorpseParty        = FieldMask{Size: 2, Offset: 0x8, Name: "FieldMaskCorpseParty"}
	FieldMaskCorpseDisplayId    = FieldMask{Size: 1, Offset: 0xA, Name: "FieldMaskCorpseDisplayId"}
	FieldMaskCorpseItem         = FieldMask{Size: 1, Offset: 0xB, Name: "FieldMaskCorpseItem"}
	FieldMaskCorpseBytes1       = FieldMask{Size: 1, Offset: 0x1E, Name: "FieldMaskCorpseBytes1"}
	FieldMaskCorpseBytes2       = FieldMask{Size: 1, Offset: 0x1F, Name: "FieldMaskCorpseBytes2"}
	FieldMaskCorpseGuild        = FieldMask{Size: 1, Offset: 0x20, Name: "FieldMaskCorpseGuild"}
	FieldMaskCorpseFlags        = FieldMask{Size: 1, Offset: 0x21, Name: "FieldMaskCorpseFlags"}
	FieldMaskCorpseDynamicFlags = FieldMask{Size: 1, Offset: 0x22, Name: "FieldMaskCorpseDynamicFlags"}
)

// ValuesMask stores the byte mask that specifices which value fields are set.
type ValuesMask struct {
	// Tracks whether any bits in the mask have been set. Used to disambiguate between no bits set
	// and only the first bit set
	anyBits    bool
	largestBit uint32
	mask       []uint32
}

const (
	maskMinSize = 16
)

// Bit returns whether the provided bit has been set.
func (m *ValuesMask) Bit(bit uint32) bool {
	if !m.anyBits || bit > m.largestBit {
		return false
	}

	index := bit / 32
	bitPos := bit % 32

	return m.mask[index]&(1<<bitPos) > 0
}

// Bytes returns a little-endian byte array of the value mask. The first byte is the size of the mask.
func (m *ValuesMask) Bytes() []byte {
	buf := bytes.Buffer{}
	size := m.Len()
	buf.WriteByte(byte(size))
	binary.Write(&buf, binary.LittleEndian, m.mask[:size])

	return buf.Bytes()
}

// FieldMask returns whether all the bits for the provided mask have been set.
func (m *ValuesMask) FieldMask(fieldMask FieldMask) bool {
	for i := uint32(0); i < fieldMask.Size; i++ {
		if !m.Bit(fieldMask.Offset + i) {
			return false
		}
	}

	return true
}

// Len returns the number of uint32s used to represent the mask.
func (m *ValuesMask) Len() int {
	size := m.largestBit / 32

	if m.anyBits {
		size++
	}

	return int(size)
}

// SetFieldMask sets all the bits necessary for the provided field mask.
func (m *ValuesMask) SetFieldMask(fieldMask FieldMask) {
	for i := uint32(0); i < fieldMask.Size; i++ {
		m.SetBit(fieldMask.Offset + i)
	}
}

// SetBit sets the nth bit in the update mask. The bit is zero-indexed with the first bit being zero.
func (m *ValuesMask) SetBit(bit uint32) {
	index := bit / 32
	bitPos := bit % 32
	m.resize(index + 1)

	if bit > m.largestBit {
		m.largestBit = bit
	}

	m.mask[index] |= 1 << bitPos
	m.anyBits = true
}

// Resizes the mask to fit up to n uint32s.
func (m *ValuesMask) resize(n uint32) {
	maskLen := uint32(len(m.mask))

	if maskLen > n {
		return
	}

	var newSize uint32

	if maskLen < maskMinSize {
		newSize = maskMinSize
	} else {
		// Grow the array exponentially
		newSize = maskLen * 2
	}

	// If it's still too small just use the desired size
	if newSize < n {
		newSize = n
	}

	oldMask := m.mask
	m.mask = make([]uint32, newSize)
	copy(m.mask, oldMask)
}
