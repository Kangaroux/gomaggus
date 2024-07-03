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
	FieldMaskItemSpellCharges       = FieldMask{Size: 5, Offset: 0x10, Name: "FieldMaskItemSpellCharges"}
	FieldMaskItemFlags              = FieldMask{Size: 1, Offset: 0x15, Name: "FieldMaskItemFlags"}
	FieldMaskItemEnchantment1_1     = FieldMask{Size: 2, Offset: 0x16, Name: "FieldMaskItemEnchantment1_1"}
	FieldMaskItemEnchantment1_3     = FieldMask{Size: 1, Offset: 0x18, Name: "FieldMaskItemEnchantment1_3"}
	FieldMaskItemEnchantment2_1     = FieldMask{Size: 2, Offset: 0x19, Name: "FieldMaskItemEnchantment2_1"}
	FieldMaskItemEnchantment2_3     = FieldMask{Size: 1, Offset: 0x1B, Name: "FieldMaskItemEnchantment2_3"}
	FieldMaskItemEnchantment3_1     = FieldMask{Size: 2, Offset: 0x1C, Name: "FieldMaskItemEnchantment3_1"}
	FieldMaskItemEnchantment3_3     = FieldMask{Size: 1, Offset: 0x1E, Name: "FieldMaskItemEnchantment3_3"}
	FieldMaskItemEnchantment4_1     = FieldMask{Size: 2, Offset: 0x1F, Name: "FieldMaskItemEnchantment4_1"}
	FieldMaskItemEnchantment4_3     = FieldMask{Size: 1, Offset: 0x21, Name: "FieldMaskItemEnchantment4_3"}
	FieldMaskItemEnchantment5_1     = FieldMask{Size: 2, Offset: 0x22, Name: "FieldMaskItemEnchantment5_1"}
	FieldMaskItemEnchantment5_3     = FieldMask{Size: 1, Offset: 0x24, Name: "FieldMaskItemEnchantment5_3"}
	FieldMaskItemEnchantment6_1     = FieldMask{Size: 2, Offset: 0x25, Name: "FieldMaskItemEnchantment6_1"}
	FieldMaskItemEnchantment6_3     = FieldMask{Size: 1, Offset: 0x27, Name: "FieldMaskItemEnchantment6_3"}
	FieldMaskItemEnchantment7_1     = FieldMask{Size: 2, Offset: 0x28, Name: "FieldMaskItemEnchantment7_1"}
	FieldMaskItemEnchantment7_3     = FieldMask{Size: 1, Offset: 0x2A, Name: "FieldMaskItemEnchantment7_3"}
	FieldMaskItemEnchantment8_1     = FieldMask{Size: 2, Offset: 0x2B, Name: "FieldMaskItemEnchantment8_1"}
	FieldMaskItemEnchantment8_3     = FieldMask{Size: 1, Offset: 0x2D, Name: "FieldMaskItemEnchantment8_3"}
	FieldMaskItemEnchantment9_1     = FieldMask{Size: 2, Offset: 0x2E, Name: "FieldMaskItemEnchantment9_1"}
	FieldMaskItemEnchantment9_3     = FieldMask{Size: 1, Offset: 0x30, Name: "FieldMaskItemEnchantment9_3"}
	FieldMaskItemEnchantment10_1    = FieldMask{Size: 2, Offset: 0x31, Name: "FieldMaskItemEnchantment10_1"}
	FieldMaskItemEnchantment10_3    = FieldMask{Size: 1, Offset: 0x33, Name: "FieldMaskItemEnchantment10_3"}
	FieldMaskItemEnchantment11_1    = FieldMask{Size: 2, Offset: 0x34, Name: "FieldMaskItemEnchantment11_1"}
	FieldMaskItemEnchantment11_3    = FieldMask{Size: 1, Offset: 0x36, Name: "FieldMaskItemEnchantment11_3"}
	FieldMaskItemEnchantment12_1    = FieldMask{Size: 2, Offset: 0x37, Name: "FieldMaskItemEnchantment12_1"}
	FieldMaskItemEnchantment12_3    = FieldMask{Size: 1, Offset: 0x39, Name: "FieldMaskItemEnchantment12_3"}
	FieldMaskItemPropertySeed       = FieldMask{Size: 1, Offset: 0x3A, Name: "FieldMaskItemPropertySeed"}
	FieldMaskItemRandomPropertiesId = FieldMask{Size: 1, Offset: 0x3B, Name: "FieldMaskItemRandomPropertiesId"}
	FieldMaskItemDurability         = FieldMask{Size: 1, Offset: 0x3C, Name: "FieldMaskItemDurability"}
	FieldMaskItemMaxdurability      = FieldMask{Size: 1, Offset: 0x3D, Name: "FieldMaskItemMaxdurability"}
	FieldMaskItemCreatePlayedTime   = FieldMask{Size: 1, Offset: 0x3E, Name: "FieldMaskItemCreatePlayedTime"}

	FieldMaskContainerNumSlots = FieldMask{Size: 1, Offset: 0x40, Name: "FieldMaskContainerNumSlots"}
	FieldMaskContainerSlot1    = FieldMask{Size: 7, Offset: 0x42, Name: "FieldMaskContainerSlot1"}

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
	FieldMaskUnitPower1                            = FieldMask{Size: 1, Offset: 0x19, Name: "FieldMaskUnitPower1"}
	FieldMaskUnitPower2                            = FieldMask{Size: 1, Offset: 0x1A, Name: "FieldMaskUnitPower2"}
	FieldMaskUnitPower3                            = FieldMask{Size: 1, Offset: 0x1B, Name: "FieldMaskUnitPower3"}
	FieldMaskUnitPower4                            = FieldMask{Size: 1, Offset: 0x1C, Name: "FieldMaskUnitPower4"}
	FieldMaskUnitPower5                            = FieldMask{Size: 1, Offset: 0x1D, Name: "FieldMaskUnitPower5"}
	FieldMaskUnitPower6                            = FieldMask{Size: 1, Offset: 0x1E, Name: "FieldMaskUnitPower6"}
	FieldMaskUnitPower7                            = FieldMask{Size: 1, Offset: 0x1F, Name: "FieldMaskUnitPower7"}
	FieldMaskUnitMaxHealth                         = FieldMask{Size: 1, Offset: 0x20, Name: "FieldMaskUnitMaxHealth"}
	FieldMaskUnitMaxpower1                         = FieldMask{Size: 1, Offset: 0x21, Name: "FieldMaskUnitMaxpower1"}
	FieldMaskUnitMaxpower2                         = FieldMask{Size: 1, Offset: 0x22, Name: "FieldMaskUnitMaxpower2"}
	FieldMaskUnitMaxpower3                         = FieldMask{Size: 1, Offset: 0x23, Name: "FieldMaskUnitMaxpower3"}
	FieldMaskUnitMaxpower4                         = FieldMask{Size: 1, Offset: 0x24, Name: "FieldMaskUnitMaxpower4"}
	FieldMaskUnitMaxpower5                         = FieldMask{Size: 1, Offset: 0x25, Name: "FieldMaskUnitMaxpower5"}
	FieldMaskUnitMaxpower6                         = FieldMask{Size: 1, Offset: 0x26, Name: "FieldMaskUnitMaxpower6"}
	FieldMaskUnitMaxpower7                         = FieldMask{Size: 1, Offset: 0x27, Name: "FieldMaskUnitMaxpower7"}
	FieldMaskUnitPowerRegenFlatModifier            = FieldMask{Size: 7, Offset: 0x28, Name: "FieldMaskUnitPowerRegenFlatModifier"}
	FieldMaskUnitPowerRegenInterruptedFlatModifier = FieldMask{Size: 7, Offset: 0x2F, Name: "FieldMaskUnitPowerRegenInterruptedFlatModifier"}
	FieldMaskUnitLevel                             = FieldMask{Size: 1, Offset: 0x36, Name: "FieldMaskUnitLevel"}
	FieldMaskUnitFactionTemplate                   = FieldMask{Size: 1, Offset: 0x37, Name: "FieldMaskUnitFactionTemplate"}
	FieldMaskUnitVirtualItemSlotId                 = FieldMask{Size: 3, Offset: 0x38, Name: "FieldMaskUnitVirtualItemSlotId"}
	FieldMaskUnitFlags                             = FieldMask{Size: 1, Offset: 0x3B, Name: "FieldMaskUnitFlags"}
	FieldMaskUnitFlags2                            = FieldMask{Size: 1, Offset: 0x3C, Name: "FieldMaskUnitFlags2"}
	FieldMaskUnitAurastate                         = FieldMask{Size: 1, Offset: 0x3D, Name: "FieldMaskUnitAurastate"}
	FieldMaskUnitBaseattacktime                    = FieldMask{Size: 2, Offset: 0x3E, Name: "FieldMaskUnitBaseattacktime"}
	FieldMaskUnitRangedattacktime                  = FieldMask{Size: 1, Offset: 0x40, Name: "FieldMaskUnitRangedattacktime"}
	FieldMaskUnitBoundingradius                    = FieldMask{Size: 1, Offset: 0x41, Name: "FieldMaskUnitBoundingradius"}
	FieldMaskUnitCombatreach                       = FieldMask{Size: 1, Offset: 0x42, Name: "FieldMaskUnitCombatreach"}
	FieldMaskUnitDisplayId                         = FieldMask{Size: 1, Offset: 0x43, Name: "FieldMaskUnitDisplayId"}
	FieldMaskUnitNativeDisplayId                   = FieldMask{Size: 1, Offset: 0x44, Name: "FieldMaskUnitNativeDisplayId"}
	FieldMaskUnitMountdisplayid                    = FieldMask{Size: 1, Offset: 0x45, Name: "FieldMaskUnitMountdisplayid"}
	FieldMaskUnitMindamage                         = FieldMask{Size: 1, Offset: 0x46, Name: "FieldMaskUnitMindamage"}
	FieldMaskUnitMaxdamage                         = FieldMask{Size: 1, Offset: 0x47, Name: "FieldMaskUnitMaxdamage"}
	FieldMaskUnitMinoffhanddamage                  = FieldMask{Size: 1, Offset: 0x48, Name: "FieldMaskUnitMinoffhanddamage"}
	FieldMaskUnitMaxoffhanddamage                  = FieldMask{Size: 1, Offset: 0x49, Name: "FieldMaskUnitMaxoffhanddamage"}
	FieldMaskUnitBytes1                            = FieldMask{Size: 1, Offset: 0x4A, Name: "FieldMaskUnitBytes1"}
	FieldMaskUnitPetnumber                         = FieldMask{Size: 1, Offset: 0x4B, Name: "FieldMaskUnitPetnumber"}
	FieldMaskUnitPetNameTimestamp                  = FieldMask{Size: 1, Offset: 0x4C, Name: "FieldMaskUnitPetNameTimestamp"}
	FieldMaskUnitPetexperience                     = FieldMask{Size: 1, Offset: 0x4D, Name: "FieldMaskUnitPetexperience"}
	FieldMaskUnitPetnextlevelexp                   = FieldMask{Size: 1, Offset: 0x4E, Name: "FieldMaskUnitPetnextlevelexp"}
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
	FieldMaskUnitPosStat0                          = FieldMask{Size: 1, Offset: 0x59, Name: "FieldMaskUnitPosStat0"}
	FieldMaskUnitPosStat1                          = FieldMask{Size: 1, Offset: 0x5A, Name: "FieldMaskUnitPosStat1"}
	FieldMaskUnitPosStat2                          = FieldMask{Size: 1, Offset: 0x5B, Name: "FieldMaskUnitPosStat2"}
	FieldMaskUnitPosStat3                          = FieldMask{Size: 1, Offset: 0x5C, Name: "FieldMaskUnitPosStat3"}
	FieldMaskUnitPosStat4                          = FieldMask{Size: 1, Offset: 0x5D, Name: "FieldMaskUnitPosStat4"}
	FieldMaskUnitNegStat0                          = FieldMask{Size: 1, Offset: 0x5E, Name: "FieldMaskUnitNegStat0"}
	FieldMaskUnitNegStat1                          = FieldMask{Size: 1, Offset: 0x5F, Name: "FieldMaskUnitNegStat1"}
	FieldMaskUnitNegStat2                          = FieldMask{Size: 1, Offset: 0x60, Name: "FieldMaskUnitNegStat2"}
	FieldMaskUnitNegStat3                          = FieldMask{Size: 1, Offset: 0x61, Name: "FieldMaskUnitNegStat3"}
	FieldMaskUnitNegStat4                          = FieldMask{Size: 1, Offset: 0x62, Name: "FieldMaskUnitNegStat4"}
	FieldMaskUnitResistances                       = FieldMask{Size: 7, Offset: 0x63, Name: "FieldMaskUnitResistances"}
	FieldMaskUnitResistanceBuffModsPositive        = FieldMask{Size: 7, Offset: 0x6A, Name: "FieldMaskUnitResistanceBuffModsPositive"}
	FieldMaskUnitResistanceBuffModsNegative        = FieldMask{Size: 7, Offset: 0x71, Name: "FieldMaskUnitResistanceBuffModsNegative"}
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
	FieldMaskUnitPowerCostModifier                 = FieldMask{Size: 7, Offset: 0x83, Name: "FieldMaskUnitPowerCostModifier"}
	FieldMaskUnitPowerCostMultiplier               = FieldMask{Size: 7, Offset: 0x8A, Name: "FieldMaskUnitPowerCostMultiplier"}
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
	FieldMaskPlayerQuestLog1_1                 = FieldMask{Size: 1, Offset: 0x9E, Name: "FieldMaskPlayerQuestLog1_1"}
	FieldMaskPlayerQuestLog1_2                 = FieldMask{Size: 1, Offset: 0x9F, Name: "FieldMaskPlayerQuestLog1_2"}
	FieldMaskPlayerQuestLog1_3                 = FieldMask{Size: 2, Offset: 0xA0, Name: "FieldMaskPlayerQuestLog1_3"}
	FieldMaskPlayerQuestLog1_4                 = FieldMask{Size: 1, Offset: 0xA2, Name: "FieldMaskPlayerQuestLog1_4"}
	FieldMaskPlayerQuestLog2_1                 = FieldMask{Size: 1, Offset: 0xA3, Name: "FieldMaskPlayerQuestLog2_1"}
	FieldMaskPlayerQuestLog2_2                 = FieldMask{Size: 1, Offset: 0xA4, Name: "FieldMaskPlayerQuestLog2_2"}
	FieldMaskPlayerQuestLog2_3                 = FieldMask{Size: 2, Offset: 0xA5, Name: "FieldMaskPlayerQuestLog2_3"}
	FieldMaskPlayerQuestLog2_5                 = FieldMask{Size: 1, Offset: 0xA7, Name: "FieldMaskPlayerQuestLog2_5"}
	FieldMaskPlayerQuestLog3_1                 = FieldMask{Size: 1, Offset: 0xA8, Name: "FieldMaskPlayerQuestLog3_1"}
	FieldMaskPlayerQuestLog3_2                 = FieldMask{Size: 1, Offset: 0xA9, Name: "FieldMaskPlayerQuestLog3_2"}
	FieldMaskPlayerQuestLog3_3                 = FieldMask{Size: 2, Offset: 0xAA, Name: "FieldMaskPlayerQuestLog3_3"}
	FieldMaskPlayerQuestLog3_5                 = FieldMask{Size: 1, Offset: 0xAC, Name: "FieldMaskPlayerQuestLog3_5"}
	FieldMaskPlayerQuestLog4_1                 = FieldMask{Size: 1, Offset: 0xAD, Name: "FieldMaskPlayerQuestLog4_1"}
	FieldMaskPlayerQuestLog4_2                 = FieldMask{Size: 1, Offset: 0xAE, Name: "FieldMaskPlayerQuestLog4_2"}
	FieldMaskPlayerQuestLog4_3                 = FieldMask{Size: 2, Offset: 0xAF, Name: "FieldMaskPlayerQuestLog4_3"}
	FieldMaskPlayerQuestLog4_5                 = FieldMask{Size: 1, Offset: 0xB1, Name: "FieldMaskPlayerQuestLog4_5"}
	FieldMaskPlayerQuestLog5_1                 = FieldMask{Size: 1, Offset: 0xB2, Name: "FieldMaskPlayerQuestLog5_1"}
	FieldMaskPlayerQuestLog5_2                 = FieldMask{Size: 1, Offset: 0xB3, Name: "FieldMaskPlayerQuestLog5_2"}
	FieldMaskPlayerQuestLog5_3                 = FieldMask{Size: 2, Offset: 0xB4, Name: "FieldMaskPlayerQuestLog5_3"}
	FieldMaskPlayerQuestLog5_5                 = FieldMask{Size: 1, Offset: 0xB6, Name: "FieldMaskPlayerQuestLog5_5"}
	FieldMaskPlayerQuestLog6_1                 = FieldMask{Size: 1, Offset: 0xB7, Name: "FieldMaskPlayerQuestLog6_1"}
	FieldMaskPlayerQuestLog6_2                 = FieldMask{Size: 1, Offset: 0xB8, Name: "FieldMaskPlayerQuestLog6_2"}
	FieldMaskPlayerQuestLog6_3                 = FieldMask{Size: 2, Offset: 0xB9, Name: "FieldMaskPlayerQuestLog6_3"}
	FieldMaskPlayerQuestLog6_5                 = FieldMask{Size: 1, Offset: 0xBB, Name: "FieldMaskPlayerQuestLog6_5"}
	FieldMaskPlayerQuestLog7_1                 = FieldMask{Size: 1, Offset: 0xBC, Name: "FieldMaskPlayerQuestLog7_1"}
	FieldMaskPlayerQuestLog7_2                 = FieldMask{Size: 1, Offset: 0xBD, Name: "FieldMaskPlayerQuestLog7_2"}
	FieldMaskPlayerQuestLog7_3                 = FieldMask{Size: 2, Offset: 0xBE, Name: "FieldMaskPlayerQuestLog7_3"}
	FieldMaskPlayerQuestLog7_5                 = FieldMask{Size: 1, Offset: 0xC0, Name: "FieldMaskPlayerQuestLog7_5"}
	FieldMaskPlayerQuestLog8_1                 = FieldMask{Size: 1, Offset: 0xC1, Name: "FieldMaskPlayerQuestLog8_1"}
	FieldMaskPlayerQuestLog8_2                 = FieldMask{Size: 1, Offset: 0xC2, Name: "FieldMaskPlayerQuestLog8_2"}
	FieldMaskPlayerQuestLog8_3                 = FieldMask{Size: 2, Offset: 0xC3, Name: "FieldMaskPlayerQuestLog8_3"}
	FieldMaskPlayerQuestLog8_5                 = FieldMask{Size: 1, Offset: 0xC5, Name: "FieldMaskPlayerQuestLog8_5"}
	FieldMaskPlayerQuestLog9_1                 = FieldMask{Size: 1, Offset: 0xC6, Name: "FieldMaskPlayerQuestLog9_1"}
	FieldMaskPlayerQuestLog9_2                 = FieldMask{Size: 1, Offset: 0xC7, Name: "FieldMaskPlayerQuestLog9_2"}
	FieldMaskPlayerQuestLog9_3                 = FieldMask{Size: 2, Offset: 0xC8, Name: "FieldMaskPlayerQuestLog9_3"}
	FieldMaskPlayerQuestLog9_5                 = FieldMask{Size: 1, Offset: 0xCA, Name: "FieldMaskPlayerQuestLog9_5"}
	FieldMaskPlayerQuestLog10_1                = FieldMask{Size: 1, Offset: 0xCB, Name: "FieldMaskPlayerQuestLog10_1"}
	FieldMaskPlayerQuestLog10_2                = FieldMask{Size: 1, Offset: 0xCC, Name: "FieldMaskPlayerQuestLog10_2"}
	FieldMaskPlayerQuestLog10_3                = FieldMask{Size: 2, Offset: 0xCD, Name: "FieldMaskPlayerQuestLog10_3"}
	FieldMaskPlayerQuestLog10_5                = FieldMask{Size: 1, Offset: 0xCF, Name: "FieldMaskPlayerQuestLog10_5"}
	FieldMaskPlayerQuestLog11_1                = FieldMask{Size: 1, Offset: 0xD0, Name: "FieldMaskPlayerQuestLog11_1"}
	FieldMaskPlayerQuestLog11_2                = FieldMask{Size: 1, Offset: 0xD1, Name: "FieldMaskPlayerQuestLog11_2"}
	FieldMaskPlayerQuestLog11_3                = FieldMask{Size: 2, Offset: 0xD2, Name: "FieldMaskPlayerQuestLog11_3"}
	FieldMaskPlayerQuestLog11_5                = FieldMask{Size: 1, Offset: 0xD4, Name: "FieldMaskPlayerQuestLog11_5"}
	FieldMaskPlayerQuestLog12_1                = FieldMask{Size: 1, Offset: 0xD5, Name: "FieldMaskPlayerQuestLog12_1"}
	FieldMaskPlayerQuestLog12_2                = FieldMask{Size: 1, Offset: 0xD6, Name: "FieldMaskPlayerQuestLog12_2"}
	FieldMaskPlayerQuestLog12_3                = FieldMask{Size: 2, Offset: 0xD7, Name: "FieldMaskPlayerQuestLog12_3"}
	FieldMaskPlayerQuestLog12_5                = FieldMask{Size: 1, Offset: 0xD9, Name: "FieldMaskPlayerQuestLog12_5"}
	FieldMaskPlayerQuestLog13_1                = FieldMask{Size: 1, Offset: 0xDA, Name: "FieldMaskPlayerQuestLog13_1"}
	FieldMaskPlayerQuestLog13_2                = FieldMask{Size: 1, Offset: 0xDB, Name: "FieldMaskPlayerQuestLog13_2"}
	FieldMaskPlayerQuestLog13_3                = FieldMask{Size: 2, Offset: 0xDC, Name: "FieldMaskPlayerQuestLog13_3"}
	FieldMaskPlayerQuestLog13_5                = FieldMask{Size: 1, Offset: 0xDE, Name: "FieldMaskPlayerQuestLog13_5"}
	FieldMaskPlayerQuestLog14_1                = FieldMask{Size: 1, Offset: 0xDF, Name: "FieldMaskPlayerQuestLog14_1"}
	FieldMaskPlayerQuestLog14_2                = FieldMask{Size: 1, Offset: 0xE0, Name: "FieldMaskPlayerQuestLog14_2"}
	FieldMaskPlayerQuestLog14_3                = FieldMask{Size: 2, Offset: 0xE1, Name: "FieldMaskPlayerQuestLog14_3"}
	FieldMaskPlayerQuestLog14_5                = FieldMask{Size: 1, Offset: 0xE3, Name: "FieldMaskPlayerQuestLog14_5"}
	FieldMaskPlayerQuestLog15_1                = FieldMask{Size: 1, Offset: 0xE4, Name: "FieldMaskPlayerQuestLog15_1"}
	FieldMaskPlayerQuestLog15_2                = FieldMask{Size: 1, Offset: 0xE5, Name: "FieldMaskPlayerQuestLog15_2"}
	FieldMaskPlayerQuestLog15_3                = FieldMask{Size: 2, Offset: 0xE6, Name: "FieldMaskPlayerQuestLog15_3"}
	FieldMaskPlayerQuestLog15_5                = FieldMask{Size: 1, Offset: 0xE8, Name: "FieldMaskPlayerQuestLog15_5"}
	FieldMaskPlayerQuestLog16_1                = FieldMask{Size: 1, Offset: 0xE9, Name: "FieldMaskPlayerQuestLog16_1"}
	FieldMaskPlayerQuestLog16_2                = FieldMask{Size: 1, Offset: 0xEA, Name: "FieldMaskPlayerQuestLog16_2"}
	FieldMaskPlayerQuestLog16_3                = FieldMask{Size: 2, Offset: 0xEB, Name: "FieldMaskPlayerQuestLog16_3"}
	FieldMaskPlayerQuestLog16_5                = FieldMask{Size: 1, Offset: 0xED, Name: "FieldMaskPlayerQuestLog16_5"}
	FieldMaskPlayerQuestLog17_1                = FieldMask{Size: 1, Offset: 0xEE, Name: "FieldMaskPlayerQuestLog17_1"}
	FieldMaskPlayerQuestLog17_2                = FieldMask{Size: 1, Offset: 0xEF, Name: "FieldMaskPlayerQuestLog17_2"}
	FieldMaskPlayerQuestLog17_3                = FieldMask{Size: 2, Offset: 0xF0, Name: "FieldMaskPlayerQuestLog17_3"}
	FieldMaskPlayerQuestLog17_5                = FieldMask{Size: 1, Offset: 0xF2, Name: "FieldMaskPlayerQuestLog17_5"}
	FieldMaskPlayerQuestLog18_1                = FieldMask{Size: 1, Offset: 0xF3, Name: "FieldMaskPlayerQuestLog18_1"}
	FieldMaskPlayerQuestLog18_2                = FieldMask{Size: 1, Offset: 0xF4, Name: "FieldMaskPlayerQuestLog18_2"}
	FieldMaskPlayerQuestLog18_3                = FieldMask{Size: 2, Offset: 0xF5, Name: "FieldMaskPlayerQuestLog18_3"}
	FieldMaskPlayerQuestLog18_5                = FieldMask{Size: 1, Offset: 0xF7, Name: "FieldMaskPlayerQuestLog18_5"}
	FieldMaskPlayerQuestLog19_1                = FieldMask{Size: 1, Offset: 0xF8, Name: "FieldMaskPlayerQuestLog19_1"}
	FieldMaskPlayerQuestLog19_2                = FieldMask{Size: 1, Offset: 0xF9, Name: "FieldMaskPlayerQuestLog19_2"}
	FieldMaskPlayerQuestLog19_3                = FieldMask{Size: 2, Offset: 0xFA, Name: "FieldMaskPlayerQuestLog19_3"}
	FieldMaskPlayerQuestLog19_5                = FieldMask{Size: 1, Offset: 0xFC, Name: "FieldMaskPlayerQuestLog19_5"}
	FieldMaskPlayerQuestLog20_1                = FieldMask{Size: 1, Offset: 0xFD, Name: "FieldMaskPlayerQuestLog20_1"}
	FieldMaskPlayerQuestLog20_2                = FieldMask{Size: 1, Offset: 0xFE, Name: "FieldMaskPlayerQuestLog20_2"}
	FieldMaskPlayerQuestLog20_3                = FieldMask{Size: 2, Offset: 0xFF, Name: "FieldMaskPlayerQuestLog20_3"}
	FieldMaskPlayerQuestLog20_5                = FieldMask{Size: 1, Offset: 0x101, Name: "FieldMaskPlayerQuestLog20_5"}
	FieldMaskPlayerQuestLog21_1                = FieldMask{Size: 1, Offset: 0x102, Name: "FieldMaskPlayerQuestLog21_1"}
	FieldMaskPlayerQuestLog21_2                = FieldMask{Size: 1, Offset: 0x103, Name: "FieldMaskPlayerQuestLog21_2"}
	FieldMaskPlayerQuestLog21_3                = FieldMask{Size: 2, Offset: 0x104, Name: "FieldMaskPlayerQuestLog21_3"}
	FieldMaskPlayerQuestLog21_5                = FieldMask{Size: 1, Offset: 0x106, Name: "FieldMaskPlayerQuestLog21_5"}
	FieldMaskPlayerQuestLog22_1                = FieldMask{Size: 1, Offset: 0x107, Name: "FieldMaskPlayerQuestLog22_1"}
	FieldMaskPlayerQuestLog22_2                = FieldMask{Size: 1, Offset: 0x108, Name: "FieldMaskPlayerQuestLog22_2"}
	FieldMaskPlayerQuestLog22_3                = FieldMask{Size: 2, Offset: 0x109, Name: "FieldMaskPlayerQuestLog22_3"}
	FieldMaskPlayerQuestLog22_5                = FieldMask{Size: 1, Offset: 0x10B, Name: "FieldMaskPlayerQuestLog22_5"}
	FieldMaskPlayerQuestLog23_1                = FieldMask{Size: 1, Offset: 0x10C, Name: "FieldMaskPlayerQuestLog23_1"}
	FieldMaskPlayerQuestLog23_2                = FieldMask{Size: 1, Offset: 0x10D, Name: "FieldMaskPlayerQuestLog23_2"}
	FieldMaskPlayerQuestLog23_3                = FieldMask{Size: 2, Offset: 0x10E, Name: "FieldMaskPlayerQuestLog23_3"}
	FieldMaskPlayerQuestLog23_5                = FieldMask{Size: 1, Offset: 0x110, Name: "FieldMaskPlayerQuestLog23_5"}
	FieldMaskPlayerQuestLog24_1                = FieldMask{Size: 1, Offset: 0x111, Name: "FieldMaskPlayerQuestLog24_1"}
	FieldMaskPlayerQuestLog24_2                = FieldMask{Size: 1, Offset: 0x112, Name: "FieldMaskPlayerQuestLog24_2"}
	FieldMaskPlayerQuestLog24_3                = FieldMask{Size: 2, Offset: 0x113, Name: "FieldMaskPlayerQuestLog24_3"}
	FieldMaskPlayerQuestLog24_5                = FieldMask{Size: 1, Offset: 0x115, Name: "FieldMaskPlayerQuestLog24_5"}
	FieldMaskPlayerQuestLog25_1                = FieldMask{Size: 1, Offset: 0x116, Name: "FieldMaskPlayerQuestLog25_1"}
	FieldMaskPlayerQuestLog25_2                = FieldMask{Size: 1, Offset: 0x117, Name: "FieldMaskPlayerQuestLog25_2"}
	FieldMaskPlayerQuestLog25_3                = FieldMask{Size: 2, Offset: 0x118, Name: "FieldMaskPlayerQuestLog25_3"}
	FieldMaskPlayerQuestLog25_5                = FieldMask{Size: 1, Offset: 0x11A, Name: "FieldMaskPlayerQuestLog25_5"}
	FieldMaskPlayerVisibleItem                 = FieldMask{Size: 3, Offset: 0x11B, Name: "FieldMaskPlayerVisibleItem"}
	FieldMaskPlayerChosenTitle                 = FieldMask{Size: 1, Offset: 0x141, Name: "FieldMaskPlayerChosenTitle"}
	FieldMaskPlayerFakeInebriation             = FieldMask{Size: 1, Offset: 0x142, Name: "FieldMaskPlayerFakeInebriation"}
	FieldMaskPlayerFieldInv                    = FieldMask{Size: 3, Offset: 0x144, Name: "FieldMaskPlayerFieldInv"}
	FieldMaskPlayerFarsight                    = FieldMask{Size: 2, Offset: 0x270, Name: "FieldMaskPlayerFarsight"}
	FieldMaskPlayerKnownTitles                 = FieldMask{Size: 2, Offset: 0x272, Name: "FieldMaskPlayerKnownTitles"}
	FieldMaskPlayerKnownTitles1                = FieldMask{Size: 2, Offset: 0x274, Name: "FieldMaskPlayerKnownTitles1"}
	FieldMaskPlayerKnownTitles2                = FieldMask{Size: 2, Offset: 0x276, Name: "FieldMaskPlayerKnownTitles2"}
	FieldMaskPlayerKnownCurrencies             = FieldMask{Size: 2, Offset: 0x278, Name: "FieldMaskPlayerKnownCurrencies"}
	FieldMaskPlayerXp                          = FieldMask{Size: 1, Offset: 0x27A, Name: "FieldMaskPlayerXp"}
	FieldMaskPlayerNextLevelXp                 = FieldMask{Size: 1, Offset: 0x27B, Name: "FieldMaskPlayerNextLevelXp"}
	FieldMaskPlayerSkillInfo                   = FieldMask{Size: 3, Offset: 0x27C, Name: "FieldMaskPlayerSkillInfo"}
	FieldMaskPlayerCharacterPoints1            = FieldMask{Size: 1, Offset: 0x3FC, Name: "FieldMaskPlayerCharacterPoints1"}
	FieldMaskPlayerCharacterPoints2            = FieldMask{Size: 1, Offset: 0x3FD, Name: "FieldMaskPlayerCharacterPoints2"}
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
	FieldMaskPlayerSpellCritPercentage1        = FieldMask{Size: 7, Offset: 0x408, Name: "FieldMaskPlayerSpellCritPercentage1"}
	FieldMaskPlayerShieldBlock                 = FieldMask{Size: 1, Offset: 0x40F, Name: "FieldMaskPlayerShieldBlock"}
	FieldMaskPlayerShieldBlockCritPercentage   = FieldMask{Size: 1, Offset: 0x410, Name: "FieldMaskPlayerShieldBlockCritPercentage"}
	FieldMaskPlayerExploredZones1              = FieldMask{Size: 1, Offset: 0x411, Name: "FieldMaskPlayerExploredZones1"}
	FieldMaskPlayerRestStateExperience         = FieldMask{Size: 1, Offset: 0x491, Name: "FieldMaskPlayerRestStateExperience"}
	FieldMaskPlayerCoinage                     = FieldMask{Size: 1, Offset: 0x492, Name: "FieldMaskPlayerCoinage"}
	FieldMaskPlayerModDamageDonePos            = FieldMask{Size: 7, Offset: 0x493, Name: "FieldMaskPlayerModDamageDonePos"}
	FieldMaskPlayerModDamageDoneNeg            = FieldMask{Size: 7, Offset: 0x49A, Name: "FieldMaskPlayerModDamageDoneNeg"}
	FieldMaskPlayerModDamageDonePct            = FieldMask{Size: 7, Offset: 0x4A1, Name: "FieldMaskPlayerModDamageDonePct"}
	FieldMaskPlayerModHealingDonePos           = FieldMask{Size: 1, Offset: 0x4A8, Name: "FieldMaskPlayerModHealingDonePos"}
	FieldMaskPlayerModHealingPct               = FieldMask{Size: 1, Offset: 0x4A9, Name: "FieldMaskPlayerModHealingPct"}
	FieldMaskPlayerModHealingDonePct           = FieldMask{Size: 1, Offset: 0x4AA, Name: "FieldMaskPlayerModHealingDonePct"}
	FieldMaskPlayerModTargetResistance         = FieldMask{Size: 1, Offset: 0x4AB, Name: "FieldMaskPlayerModTargetResistance"}
	FieldMaskPlayerModTargetPhysicalResistance = FieldMask{Size: 1, Offset: 0x4AC, Name: "FieldMaskPlayerModTargetPhysicalResistance"}
	FieldMaskPlayerFeatures                    = FieldMask{Size: 1, Offset: 0x4AD, Name: "FieldMaskPlayerFeatures"}
	FieldMaskPlayerAmmoId                      = FieldMask{Size: 1, Offset: 0x4AE, Name: "FieldMaskPlayerAmmoId"}
	FieldMaskPlayerSelfResSpell                = FieldMask{Size: 1, Offset: 0x4AF, Name: "FieldMaskPlayerSelfResSpell"}
	FieldMaskPlayerPvpMedals                   = FieldMask{Size: 1, Offset: 0x4B0, Name: "FieldMaskPlayerPvpMedals"}
	FieldMaskPlayerBuybackPrice1               = FieldMask{Size: 1, Offset: 0x4B1, Name: "FieldMaskPlayerBuybackPrice1"}
	FieldMaskPlayerBuybackTimestamp1           = FieldMask{Size: 1, Offset: 0x4BD, Name: "FieldMaskPlayerBuybackTimestamp1"}
	FieldMaskPlayerKills                       = FieldMask{Size: 1, Offset: 0x4C9, Name: "FieldMaskPlayerKills"}
	FieldMaskPlayerTodayContribution           = FieldMask{Size: 1, Offset: 0x4CA, Name: "FieldMaskPlayerTodayContribution"}
	FieldMaskPlayerYesterdayContribution       = FieldMask{Size: 1, Offset: 0x4CB, Name: "FieldMaskPlayerYesterdayContribution"}
	FieldMaskPlayerLifetimeHonorableKills      = FieldMask{Size: 1, Offset: 0x4CC, Name: "FieldMaskPlayerLifetimeHonorableKills"}
	FieldMaskPlayerBytes4                      = FieldMask{Size: 1, Offset: 0x4CD, Name: "FieldMaskPlayerBytes4"}
	FieldMaskPlayerWatchedFactionIndex         = FieldMask{Size: 1, Offset: 0x4CE, Name: "FieldMaskPlayerWatchedFactionIndex"}
	FieldMaskPlayerCombatRating1               = FieldMask{Size: 2, Offset: 0x4CF, Name: "FieldMaskPlayerCombatRating1"}
	FieldMaskPlayerArenaTeamInfo11             = FieldMask{Size: 2, Offset: 0x4E8, Name: "FieldMaskPlayerArenaTeamInfo11"}
	FieldMaskPlayerHonorCurrency               = FieldMask{Size: 1, Offset: 0x4FD, Name: "FieldMaskPlayerHonorCurrency"}
	FieldMaskPlayerArenaCurrency               = FieldMask{Size: 1, Offset: 0x4FE, Name: "FieldMaskPlayerArenaCurrency"}
	FieldMaskPlayerMaxLevel                    = FieldMask{Size: 1, Offset: 0x4FF, Name: "FieldMaskPlayerMaxLevel"}
	FieldMaskPlayerDailyQuests1                = FieldMask{Size: 2, Offset: 0x500, Name: "FieldMaskPlayerDailyQuests1"}
	FieldMaskPlayerRuneRegen1                  = FieldMask{Size: 4, Offset: 0x519, Name: "FieldMaskPlayerRuneRegen1"}
	FieldMaskPlayerNoReagentCost1              = FieldMask{Size: 3, Offset: 0x51D, Name: "FieldMaskPlayerNoReagentCost1"}
	FieldMaskPlayerGlyphSlots1                 = FieldMask{Size: 6, Offset: 0x520, Name: "FieldMaskPlayerGlyphSlots1"}
	FieldMaskPlayerGlyphs1                     = FieldMask{Size: 6, Offset: 0x526, Name: "FieldMaskPlayerGlyphs1"}
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
