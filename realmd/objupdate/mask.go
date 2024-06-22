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
	name   string
}

func (fm *FieldMask) String() string {
	return fm.name
}

var (
	FieldMaskObjectGuid    = FieldMask{Size: 2, Offset: 0x0, name: "FieldMaskObjectGuid"}
	FieldMaskObjectType    = FieldMask{Size: 1, Offset: 0x2, name: "FieldMaskObjectType"}
	FieldMaskObjectEntry   = FieldMask{Size: 1, Offset: 0x3, name: "FieldMaskObjectEntry"}
	FieldMaskObjectScaleX  = FieldMask{Size: 1, Offset: 0x4, name: "FieldMaskObjectScaleX"}
	FieldMaskObjectPadding = FieldMask{Size: 1, Offset: 0x5, name: "FieldMaskObjectPadding"}

	FieldMaskItemOwner              = FieldMask{Size: 2, Offset: 0x6, name: "FieldMaskItemOwner"}
	FieldMaskItemContained          = FieldMask{Size: 2, Offset: 0x8, name: "FieldMaskItemContained"}
	FieldMaskItemCreator            = FieldMask{Size: 2, Offset: 0xA, name: "FieldMaskItemCreator"}
	FieldMaskItemGiftcreator        = FieldMask{Size: 2, Offset: 0xC, name: "FieldMaskItemGiftcreator"}
	FieldMaskItemStackCount         = FieldMask{Size: 1, Offset: 0xE, name: "FieldMaskItemStackCount"}
	FieldMaskItemDuration           = FieldMask{Size: 1, Offset: 0xF, name: "FieldMaskItemDuration"}
	FieldMaskItemSpellCharges       = FieldMask{Size: 5, Offset: 0x10, name: "FieldMaskItemSpellCharges"}
	FieldMaskItemFlags              = FieldMask{Size: 1, Offset: 0x15, name: "FieldMaskItemFlags"}
	FieldMaskItemEnchantment1_1     = FieldMask{Size: 2, Offset: 0x16, name: "FieldMaskItemEnchantment1_1"}
	FieldMaskItemEnchantment1_3     = FieldMask{Size: 1, Offset: 0x18, name: "FieldMaskItemEnchantment1_3"}
	FieldMaskItemEnchantment2_1     = FieldMask{Size: 2, Offset: 0x19, name: "FieldMaskItemEnchantment2_1"}
	FieldMaskItemEnchantment2_3     = FieldMask{Size: 1, Offset: 0x1B, name: "FieldMaskItemEnchantment2_3"}
	FieldMaskItemEnchantment3_1     = FieldMask{Size: 2, Offset: 0x1C, name: "FieldMaskItemEnchantment3_1"}
	FieldMaskItemEnchantment3_3     = FieldMask{Size: 1, Offset: 0x1E, name: "FieldMaskItemEnchantment3_3"}
	FieldMaskItemEnchantment4_1     = FieldMask{Size: 2, Offset: 0x1F, name: "FieldMaskItemEnchantment4_1"}
	FieldMaskItemEnchantment4_3     = FieldMask{Size: 1, Offset: 0x21, name: "FieldMaskItemEnchantment4_3"}
	FieldMaskItemEnchantment5_1     = FieldMask{Size: 2, Offset: 0x22, name: "FieldMaskItemEnchantment5_1"}
	FieldMaskItemEnchantment5_3     = FieldMask{Size: 1, Offset: 0x24, name: "FieldMaskItemEnchantment5_3"}
	FieldMaskItemEnchantment6_1     = FieldMask{Size: 2, Offset: 0x25, name: "FieldMaskItemEnchantment6_1"}
	FieldMaskItemEnchantment6_3     = FieldMask{Size: 1, Offset: 0x27, name: "FieldMaskItemEnchantment6_3"}
	FieldMaskItemEnchantment7_1     = FieldMask{Size: 2, Offset: 0x28, name: "FieldMaskItemEnchantment7_1"}
	FieldMaskItemEnchantment7_3     = FieldMask{Size: 1, Offset: 0x2A, name: "FieldMaskItemEnchantment7_3"}
	FieldMaskItemEnchantment8_1     = FieldMask{Size: 2, Offset: 0x2B, name: "FieldMaskItemEnchantment8_1"}
	FieldMaskItemEnchantment8_3     = FieldMask{Size: 1, Offset: 0x2D, name: "FieldMaskItemEnchantment8_3"}
	FieldMaskItemEnchantment9_1     = FieldMask{Size: 2, Offset: 0x2E, name: "FieldMaskItemEnchantment9_1"}
	FieldMaskItemEnchantment9_3     = FieldMask{Size: 1, Offset: 0x30, name: "FieldMaskItemEnchantment9_3"}
	FieldMaskItemEnchantment10_1    = FieldMask{Size: 2, Offset: 0x31, name: "FieldMaskItemEnchantment10_1"}
	FieldMaskItemEnchantment10_3    = FieldMask{Size: 1, Offset: 0x33, name: "FieldMaskItemEnchantment10_3"}
	FieldMaskItemEnchantment11_1    = FieldMask{Size: 2, Offset: 0x34, name: "FieldMaskItemEnchantment11_1"}
	FieldMaskItemEnchantment11_3    = FieldMask{Size: 1, Offset: 0x36, name: "FieldMaskItemEnchantment11_3"}
	FieldMaskItemEnchantment12_1    = FieldMask{Size: 2, Offset: 0x37, name: "FieldMaskItemEnchantment12_1"}
	FieldMaskItemEnchantment12_3    = FieldMask{Size: 1, Offset: 0x39, name: "FieldMaskItemEnchantment12_3"}
	FieldMaskItemPropertySeed       = FieldMask{Size: 1, Offset: 0x3A, name: "FieldMaskItemPropertySeed"}
	FieldMaskItemRandomPropertiesId = FieldMask{Size: 1, Offset: 0x3B, name: "FieldMaskItemRandomPropertiesId"}
	FieldMaskItemDurability         = FieldMask{Size: 1, Offset: 0x3C, name: "FieldMaskItemDurability"}
	FieldMaskItemMaxdurability      = FieldMask{Size: 1, Offset: 0x3D, name: "FieldMaskItemMaxdurability"}
	FieldMaskItemCreatePlayedTime   = FieldMask{Size: 1, Offset: 0x3E, name: "FieldMaskItemCreatePlayedTime"}

	FieldMaskContainerNumSlots = FieldMask{Size: 1, Offset: 0x40, name: "FieldMaskContainerNumSlots"}
	FieldMaskContainerSlot1    = FieldMask{Size: 7, Offset: 0x42, name: "FieldMaskContainerSlot1"}

	FieldMaskUnitCharm                             = FieldMask{Size: 2, Offset: 0x6, name: "FieldMaskUnitCharm"}
	FieldMaskUnitSummon                            = FieldMask{Size: 2, Offset: 0x8, name: "FieldMaskUnitSummon"}
	FieldMaskUnitCritter                           = FieldMask{Size: 2, Offset: 0xA, name: "FieldMaskUnitCritter"}
	FieldMaskUnitCharmedby                         = FieldMask{Size: 2, Offset: 0xC, name: "FieldMaskUnitCharmedby"}
	FieldMaskUnitSummonedby                        = FieldMask{Size: 2, Offset: 0xE, name: "FieldMaskUnitSummonedby"}
	FieldMaskUnitCreatedby                         = FieldMask{Size: 2, Offset: 0x10, name: "FieldMaskUnitCreatedby"}
	FieldMaskUnitTarget                            = FieldMask{Size: 2, Offset: 0x12, name: "FieldMaskUnitTarget"}
	FieldMaskUnitChannelObject                     = FieldMask{Size: 2, Offset: 0x14, name: "FieldMaskUnitChannelObject"}
	FieldMaskUnitChannelSpell                      = FieldMask{Size: 1, Offset: 0x16, name: "FieldMaskUnitChannelSpell"}
	FieldMaskUnitRaceClassGenderPower              = FieldMask{Size: 1, Offset: 0x17, name: "FieldMaskUnitRaceClassGenderPower"} // Bytes0
	FieldMaskUnitHealth                            = FieldMask{Size: 1, Offset: 0x18, name: "FieldMaskUnitHealth"}
	FieldMaskUnitPower1                            = FieldMask{Size: 1, Offset: 0x19, name: "FieldMaskUnitPower1"}
	FieldMaskUnitPower2                            = FieldMask{Size: 1, Offset: 0x1A, name: "FieldMaskUnitPower2"}
	FieldMaskUnitPower3                            = FieldMask{Size: 1, Offset: 0x1B, name: "FieldMaskUnitPower3"}
	FieldMaskUnitPower4                            = FieldMask{Size: 1, Offset: 0x1C, name: "FieldMaskUnitPower4"}
	FieldMaskUnitPower5                            = FieldMask{Size: 1, Offset: 0x1D, name: "FieldMaskUnitPower5"}
	FieldMaskUnitPower6                            = FieldMask{Size: 1, Offset: 0x1E, name: "FieldMaskUnitPower6"}
	FieldMaskUnitPower7                            = FieldMask{Size: 1, Offset: 0x1F, name: "FieldMaskUnitPower7"}
	FieldMaskUnitMaxHealth                         = FieldMask{Size: 1, Offset: 0x20, name: "FieldMaskUnitMaxHealth"}
	FieldMaskUnitMaxpower1                         = FieldMask{Size: 1, Offset: 0x21, name: "FieldMaskUnitMaxpower1"}
	FieldMaskUnitMaxpower2                         = FieldMask{Size: 1, Offset: 0x22, name: "FieldMaskUnitMaxpower2"}
	FieldMaskUnitMaxpower3                         = FieldMask{Size: 1, Offset: 0x23, name: "FieldMaskUnitMaxpower3"}
	FieldMaskUnitMaxpower4                         = FieldMask{Size: 1, Offset: 0x24, name: "FieldMaskUnitMaxpower4"}
	FieldMaskUnitMaxpower5                         = FieldMask{Size: 1, Offset: 0x25, name: "FieldMaskUnitMaxpower5"}
	FieldMaskUnitMaxpower6                         = FieldMask{Size: 1, Offset: 0x26, name: "FieldMaskUnitMaxpower6"}
	FieldMaskUnitMaxpower7                         = FieldMask{Size: 1, Offset: 0x27, name: "FieldMaskUnitMaxpower7"}
	FieldMaskUnitPowerRegenFlatModifier            = FieldMask{Size: 7, Offset: 0x28, name: "FieldMaskUnitPowerRegenFlatModifier"}
	FieldMaskUnitPowerRegenInterruptedFlatModifier = FieldMask{Size: 7, Offset: 0x2F, name: "FieldMaskUnitPowerRegenInterruptedFlatModifier"}
	FieldMaskUnitLevel                             = FieldMask{Size: 1, Offset: 0x36, name: "FieldMaskUnitLevel"}
	FieldMaskUnitFactionTemplate                   = FieldMask{Size: 1, Offset: 0x37, name: "FieldMaskUnitFactionTemplate"}
	FieldMaskUnitVirtualItemSlotId                 = FieldMask{Size: 3, Offset: 0x38, name: "FieldMaskUnitVirtualItemSlotId"}
	FieldMaskUnitFlags                             = FieldMask{Size: 1, Offset: 0x3B, name: "FieldMaskUnitFlags"}
	FieldMaskUnitFlags2                            = FieldMask{Size: 1, Offset: 0x3C, name: "FieldMaskUnitFlags2"}
	FieldMaskUnitAurastate                         = FieldMask{Size: 1, Offset: 0x3D, name: "FieldMaskUnitAurastate"}
	FieldMaskUnitBaseattacktime                    = FieldMask{Size: 2, Offset: 0x3E, name: "FieldMaskUnitBaseattacktime"}
	FieldMaskUnitRangedattacktime                  = FieldMask{Size: 1, Offset: 0x40, name: "FieldMaskUnitRangedattacktime"}
	FieldMaskUnitBoundingradius                    = FieldMask{Size: 1, Offset: 0x41, name: "FieldMaskUnitBoundingradius"}
	FieldMaskUnitCombatreach                       = FieldMask{Size: 1, Offset: 0x42, name: "FieldMaskUnitCombatreach"}
	FieldMaskUnitDisplayId                         = FieldMask{Size: 1, Offset: 0x43, name: "FieldMaskUnitDisplayId"}
	FieldMaskUnitNativeDisplayId                   = FieldMask{Size: 1, Offset: 0x44, name: "FieldMaskUnitNativeDisplayId"}
	FieldMaskUnitMountdisplayid                    = FieldMask{Size: 1, Offset: 0x45, name: "FieldMaskUnitMountdisplayid"}
	FieldMaskUnitMindamage                         = FieldMask{Size: 1, Offset: 0x46, name: "FieldMaskUnitMindamage"}
	FieldMaskUnitMaxdamage                         = FieldMask{Size: 1, Offset: 0x47, name: "FieldMaskUnitMaxdamage"}
	FieldMaskUnitMinoffhanddamage                  = FieldMask{Size: 1, Offset: 0x48, name: "FieldMaskUnitMinoffhanddamage"}
	FieldMaskUnitMaxoffhanddamage                  = FieldMask{Size: 1, Offset: 0x49, name: "FieldMaskUnitMaxoffhanddamage"}
	FieldMaskUnitBytes1                            = FieldMask{Size: 1, Offset: 0x4A, name: "FieldMaskUnitBytes1"}
	FieldMaskUnitPetnumber                         = FieldMask{Size: 1, Offset: 0x4B, name: "FieldMaskUnitPetnumber"}
	FieldMaskUnitPetNameTimestamp                  = FieldMask{Size: 1, Offset: 0x4C, name: "FieldMaskUnitPetNameTimestamp"}
	FieldMaskUnitPetexperience                     = FieldMask{Size: 1, Offset: 0x4D, name: "FieldMaskUnitPetexperience"}
	FieldMaskUnitPetnextlevelexp                   = FieldMask{Size: 1, Offset: 0x4E, name: "FieldMaskUnitPetnextlevelexp"}
	FieldMaskUnitDynamicFlags                      = FieldMask{Size: 1, Offset: 0x4F, name: "FieldMaskUnitDynamicFlags"}
	FieldMaskUnitModCastSpeed                      = FieldMask{Size: 1, Offset: 0x50, name: "FieldMaskUnitModCastSpeed"}
	FieldMaskUnitCreatedBySpell                    = FieldMask{Size: 1, Offset: 0x51, name: "FieldMaskUnitCreatedBySpell"}
	FieldMaskUnitNpcFlags                          = FieldMask{Size: 1, Offset: 0x52, name: "FieldMaskUnitNpcFlags"}
	FieldMaskUnitNpcEmotestate                     = FieldMask{Size: 1, Offset: 0x53, name: "FieldMaskUnitNpcEmotestate"}
	FieldMaskUnitStrength                          = FieldMask{Size: 1, Offset: 0x54, name: "FieldMaskUnitStrength"}
	FieldMaskUnitAgility                           = FieldMask{Size: 1, Offset: 0x55, name: "FieldMaskUnitAgility"}
	FieldMaskUnitStamina                           = FieldMask{Size: 1, Offset: 0x56, name: "FieldMaskUnitStamina"}
	FieldMaskUnitIntellect                         = FieldMask{Size: 1, Offset: 0x57, name: "FieldMaskUnitIntellect"}
	FieldMaskUnitSpirit                            = FieldMask{Size: 1, Offset: 0x58, name: "FieldMaskUnitSpirit"}
	FieldMaskUnitPosStat0                          = FieldMask{Size: 1, Offset: 0x59, name: "FieldMaskUnitPosStat0"}
	FieldMaskUnitPosStat1                          = FieldMask{Size: 1, Offset: 0x5A, name: "FieldMaskUnitPosStat1"}
	FieldMaskUnitPosStat2                          = FieldMask{Size: 1, Offset: 0x5B, name: "FieldMaskUnitPosStat2"}
	FieldMaskUnitPosStat3                          = FieldMask{Size: 1, Offset: 0x5C, name: "FieldMaskUnitPosStat3"}
	FieldMaskUnitPosStat4                          = FieldMask{Size: 1, Offset: 0x5D, name: "FieldMaskUnitPosStat4"}
	FieldMaskUnitNegStat0                          = FieldMask{Size: 1, Offset: 0x5E, name: "FieldMaskUnitNegStat0"}
	FieldMaskUnitNegStat1                          = FieldMask{Size: 1, Offset: 0x5F, name: "FieldMaskUnitNegStat1"}
	FieldMaskUnitNegStat2                          = FieldMask{Size: 1, Offset: 0x60, name: "FieldMaskUnitNegStat2"}
	FieldMaskUnitNegStat3                          = FieldMask{Size: 1, Offset: 0x61, name: "FieldMaskUnitNegStat3"}
	FieldMaskUnitNegStat4                          = FieldMask{Size: 1, Offset: 0x62, name: "FieldMaskUnitNegStat4"}
	FieldMaskUnitResistances                       = FieldMask{Size: 7, Offset: 0x63, name: "FieldMaskUnitResistances"}
	FieldMaskUnitResistanceBuffModsPositive        = FieldMask{Size: 7, Offset: 0x6A, name: "FieldMaskUnitResistanceBuffModsPositive"}
	FieldMaskUnitResistanceBuffModsNegative        = FieldMask{Size: 7, Offset: 0x71, name: "FieldMaskUnitResistanceBuffModsNegative"}
	FieldMaskUnitBaseMana                          = FieldMask{Size: 1, Offset: 0x78, name: "FieldMaskUnitBaseMana"}
	FieldMaskUnitBaseHealth                        = FieldMask{Size: 1, Offset: 0x79, name: "FieldMaskUnitBaseHealth"}
	FieldMaskUnitBytes2                            = FieldMask{Size: 1, Offset: 0x7A, name: "FieldMaskUnitBytes2"}
	FieldMaskUnitAttackPower                       = FieldMask{Size: 1, Offset: 0x7B, name: "FieldMaskUnitAttackPower"}
	FieldMaskUnitAttackPowerMods                   = FieldMask{Size: 1, Offset: 0x7C, name: "FieldMaskUnitAttackPowerMods"}
	FieldMaskUnitAttackPowerMultiplier             = FieldMask{Size: 1, Offset: 0x7D, name: "FieldMaskUnitAttackPowerMultiplier"}
	FieldMaskUnitRangedAttackPower                 = FieldMask{Size: 1, Offset: 0x7E, name: "FieldMaskUnitRangedAttackPower"}
	FieldMaskUnitRangedAttackPowerMods             = FieldMask{Size: 1, Offset: 0x7F, name: "FieldMaskUnitRangedAttackPowerMods"}
	FieldMaskUnitRangedAttackPowerMultiplier       = FieldMask{Size: 1, Offset: 0x80, name: "FieldMaskUnitRangedAttackPowerMultiplier"}
	FieldMaskUnitMinRangedDamage                   = FieldMask{Size: 1, Offset: 0x81, name: "FieldMaskUnitMinRangedDamage"}
	FieldMaskUnitMaxRangedDamage                   = FieldMask{Size: 1, Offset: 0x82, name: "FieldMaskUnitMaxRangedDamage"}
	FieldMaskUnitPowerCostModifier                 = FieldMask{Size: 7, Offset: 0x83, name: "FieldMaskUnitPowerCostModifier"}
	FieldMaskUnitPowerCostMultiplier               = FieldMask{Size: 7, Offset: 0x8A, name: "FieldMaskUnitPowerCostMultiplier"}
	FieldMaskUnitMaxHealthModifier                 = FieldMask{Size: 1, Offset: 0x91, name: "FieldMaskUnitMaxHealthModifier"}
	FieldMaskUnitHoverHeight                       = FieldMask{Size: 1, Offset: 0x92, name: "FieldMaskUnitHoverHeight"}

	FieldMaskPlayerDuelArbiter                 = FieldMask{Size: 2, Offset: 0x94, name: "FieldMaskPlayerDuelArbiter"}
	FieldMaskPlayerFlags                       = FieldMask{Size: 1, Offset: 0x96, name: "FieldMaskPlayerFlags"}
	FieldMaskPlayerGuildid                     = FieldMask{Size: 1, Offset: 0x97, name: "FieldMaskPlayerGuildid"}
	FieldMaskPlayerGuildrank                   = FieldMask{Size: 1, Offset: 0x98, name: "FieldMaskPlayerGuildrank"}
	FieldMaskPlayerFieldBytes                  = FieldMask{Size: 1, Offset: 0x99, name: "FieldMaskPlayerFieldBytes"}
	FieldMaskPlayerBytes2                      = FieldMask{Size: 1, Offset: 0x9A, name: "FieldMaskPlayerBytes2"}
	FieldMaskPlayerBytes3                      = FieldMask{Size: 1, Offset: 0x9B, name: "FieldMaskPlayerBytes3"}
	FieldMaskPlayerDuelTeam                    = FieldMask{Size: 1, Offset: 0x9C, name: "FieldMaskPlayerDuelTeam"}
	FieldMaskPlayerGuildTimestamp              = FieldMask{Size: 1, Offset: 0x9D, name: "FieldMaskPlayerGuildTimestamp"}
	FieldMaskPlayerQuestLog1_1                 = FieldMask{Size: 1, Offset: 0x9E, name: "FieldMaskPlayerQuestLog1_1"}
	FieldMaskPlayerQuestLog1_2                 = FieldMask{Size: 1, Offset: 0x9F, name: "FieldMaskPlayerQuestLog1_2"}
	FieldMaskPlayerQuestLog1_3                 = FieldMask{Size: 2, Offset: 0xA0, name: "FieldMaskPlayerQuestLog1_3"}
	FieldMaskPlayerQuestLog1_4                 = FieldMask{Size: 1, Offset: 0xA2, name: "FieldMaskPlayerQuestLog1_4"}
	FieldMaskPlayerQuestLog2_1                 = FieldMask{Size: 1, Offset: 0xA3, name: "FieldMaskPlayerQuestLog2_1"}
	FieldMaskPlayerQuestLog2_2                 = FieldMask{Size: 1, Offset: 0xA4, name: "FieldMaskPlayerQuestLog2_2"}
	FieldMaskPlayerQuestLog2_3                 = FieldMask{Size: 2, Offset: 0xA5, name: "FieldMaskPlayerQuestLog2_3"}
	FieldMaskPlayerQuestLog2_5                 = FieldMask{Size: 1, Offset: 0xA7, name: "FieldMaskPlayerQuestLog2_5"}
	FieldMaskPlayerQuestLog3_1                 = FieldMask{Size: 1, Offset: 0xA8, name: "FieldMaskPlayerQuestLog3_1"}
	FieldMaskPlayerQuestLog3_2                 = FieldMask{Size: 1, Offset: 0xA9, name: "FieldMaskPlayerQuestLog3_2"}
	FieldMaskPlayerQuestLog3_3                 = FieldMask{Size: 2, Offset: 0xAA, name: "FieldMaskPlayerQuestLog3_3"}
	FieldMaskPlayerQuestLog3_5                 = FieldMask{Size: 1, Offset: 0xAC, name: "FieldMaskPlayerQuestLog3_5"}
	FieldMaskPlayerQuestLog4_1                 = FieldMask{Size: 1, Offset: 0xAD, name: "FieldMaskPlayerQuestLog4_1"}
	FieldMaskPlayerQuestLog4_2                 = FieldMask{Size: 1, Offset: 0xAE, name: "FieldMaskPlayerQuestLog4_2"}
	FieldMaskPlayerQuestLog4_3                 = FieldMask{Size: 2, Offset: 0xAF, name: "FieldMaskPlayerQuestLog4_3"}
	FieldMaskPlayerQuestLog4_5                 = FieldMask{Size: 1, Offset: 0xB1, name: "FieldMaskPlayerQuestLog4_5"}
	FieldMaskPlayerQuestLog5_1                 = FieldMask{Size: 1, Offset: 0xB2, name: "FieldMaskPlayerQuestLog5_1"}
	FieldMaskPlayerQuestLog5_2                 = FieldMask{Size: 1, Offset: 0xB3, name: "FieldMaskPlayerQuestLog5_2"}
	FieldMaskPlayerQuestLog5_3                 = FieldMask{Size: 2, Offset: 0xB4, name: "FieldMaskPlayerQuestLog5_3"}
	FieldMaskPlayerQuestLog5_5                 = FieldMask{Size: 1, Offset: 0xB6, name: "FieldMaskPlayerQuestLog5_5"}
	FieldMaskPlayerQuestLog6_1                 = FieldMask{Size: 1, Offset: 0xB7, name: "FieldMaskPlayerQuestLog6_1"}
	FieldMaskPlayerQuestLog6_2                 = FieldMask{Size: 1, Offset: 0xB8, name: "FieldMaskPlayerQuestLog6_2"}
	FieldMaskPlayerQuestLog6_3                 = FieldMask{Size: 2, Offset: 0xB9, name: "FieldMaskPlayerQuestLog6_3"}
	FieldMaskPlayerQuestLog6_5                 = FieldMask{Size: 1, Offset: 0xBB, name: "FieldMaskPlayerQuestLog6_5"}
	FieldMaskPlayerQuestLog7_1                 = FieldMask{Size: 1, Offset: 0xBC, name: "FieldMaskPlayerQuestLog7_1"}
	FieldMaskPlayerQuestLog7_2                 = FieldMask{Size: 1, Offset: 0xBD, name: "FieldMaskPlayerQuestLog7_2"}
	FieldMaskPlayerQuestLog7_3                 = FieldMask{Size: 2, Offset: 0xBE, name: "FieldMaskPlayerQuestLog7_3"}
	FieldMaskPlayerQuestLog7_5                 = FieldMask{Size: 1, Offset: 0xC0, name: "FieldMaskPlayerQuestLog7_5"}
	FieldMaskPlayerQuestLog8_1                 = FieldMask{Size: 1, Offset: 0xC1, name: "FieldMaskPlayerQuestLog8_1"}
	FieldMaskPlayerQuestLog8_2                 = FieldMask{Size: 1, Offset: 0xC2, name: "FieldMaskPlayerQuestLog8_2"}
	FieldMaskPlayerQuestLog8_3                 = FieldMask{Size: 2, Offset: 0xC3, name: "FieldMaskPlayerQuestLog8_3"}
	FieldMaskPlayerQuestLog8_5                 = FieldMask{Size: 1, Offset: 0xC5, name: "FieldMaskPlayerQuestLog8_5"}
	FieldMaskPlayerQuestLog9_1                 = FieldMask{Size: 1, Offset: 0xC6, name: "FieldMaskPlayerQuestLog9_1"}
	FieldMaskPlayerQuestLog9_2                 = FieldMask{Size: 1, Offset: 0xC7, name: "FieldMaskPlayerQuestLog9_2"}
	FieldMaskPlayerQuestLog9_3                 = FieldMask{Size: 2, Offset: 0xC8, name: "FieldMaskPlayerQuestLog9_3"}
	FieldMaskPlayerQuestLog9_5                 = FieldMask{Size: 1, Offset: 0xCA, name: "FieldMaskPlayerQuestLog9_5"}
	FieldMaskPlayerQuestLog10_1                = FieldMask{Size: 1, Offset: 0xCB, name: "FieldMaskPlayerQuestLog10_1"}
	FieldMaskPlayerQuestLog10_2                = FieldMask{Size: 1, Offset: 0xCC, name: "FieldMaskPlayerQuestLog10_2"}
	FieldMaskPlayerQuestLog10_3                = FieldMask{Size: 2, Offset: 0xCD, name: "FieldMaskPlayerQuestLog10_3"}
	FieldMaskPlayerQuestLog10_5                = FieldMask{Size: 1, Offset: 0xCF, name: "FieldMaskPlayerQuestLog10_5"}
	FieldMaskPlayerQuestLog11_1                = FieldMask{Size: 1, Offset: 0xD0, name: "FieldMaskPlayerQuestLog11_1"}
	FieldMaskPlayerQuestLog11_2                = FieldMask{Size: 1, Offset: 0xD1, name: "FieldMaskPlayerQuestLog11_2"}
	FieldMaskPlayerQuestLog11_3                = FieldMask{Size: 2, Offset: 0xD2, name: "FieldMaskPlayerQuestLog11_3"}
	FieldMaskPlayerQuestLog11_5                = FieldMask{Size: 1, Offset: 0xD4, name: "FieldMaskPlayerQuestLog11_5"}
	FieldMaskPlayerQuestLog12_1                = FieldMask{Size: 1, Offset: 0xD5, name: "FieldMaskPlayerQuestLog12_1"}
	FieldMaskPlayerQuestLog12_2                = FieldMask{Size: 1, Offset: 0xD6, name: "FieldMaskPlayerQuestLog12_2"}
	FieldMaskPlayerQuestLog12_3                = FieldMask{Size: 2, Offset: 0xD7, name: "FieldMaskPlayerQuestLog12_3"}
	FieldMaskPlayerQuestLog12_5                = FieldMask{Size: 1, Offset: 0xD9, name: "FieldMaskPlayerQuestLog12_5"}
	FieldMaskPlayerQuestLog13_1                = FieldMask{Size: 1, Offset: 0xDA, name: "FieldMaskPlayerQuestLog13_1"}
	FieldMaskPlayerQuestLog13_2                = FieldMask{Size: 1, Offset: 0xDB, name: "FieldMaskPlayerQuestLog13_2"}
	FieldMaskPlayerQuestLog13_3                = FieldMask{Size: 2, Offset: 0xDC, name: "FieldMaskPlayerQuestLog13_3"}
	FieldMaskPlayerQuestLog13_5                = FieldMask{Size: 1, Offset: 0xDE, name: "FieldMaskPlayerQuestLog13_5"}
	FieldMaskPlayerQuestLog14_1                = FieldMask{Size: 1, Offset: 0xDF, name: "FieldMaskPlayerQuestLog14_1"}
	FieldMaskPlayerQuestLog14_2                = FieldMask{Size: 1, Offset: 0xE0, name: "FieldMaskPlayerQuestLog14_2"}
	FieldMaskPlayerQuestLog14_3                = FieldMask{Size: 2, Offset: 0xE1, name: "FieldMaskPlayerQuestLog14_3"}
	FieldMaskPlayerQuestLog14_5                = FieldMask{Size: 1, Offset: 0xE3, name: "FieldMaskPlayerQuestLog14_5"}
	FieldMaskPlayerQuestLog15_1                = FieldMask{Size: 1, Offset: 0xE4, name: "FieldMaskPlayerQuestLog15_1"}
	FieldMaskPlayerQuestLog15_2                = FieldMask{Size: 1, Offset: 0xE5, name: "FieldMaskPlayerQuestLog15_2"}
	FieldMaskPlayerQuestLog15_3                = FieldMask{Size: 2, Offset: 0xE6, name: "FieldMaskPlayerQuestLog15_3"}
	FieldMaskPlayerQuestLog15_5                = FieldMask{Size: 1, Offset: 0xE8, name: "FieldMaskPlayerQuestLog15_5"}
	FieldMaskPlayerQuestLog16_1                = FieldMask{Size: 1, Offset: 0xE9, name: "FieldMaskPlayerQuestLog16_1"}
	FieldMaskPlayerQuestLog16_2                = FieldMask{Size: 1, Offset: 0xEA, name: "FieldMaskPlayerQuestLog16_2"}
	FieldMaskPlayerQuestLog16_3                = FieldMask{Size: 2, Offset: 0xEB, name: "FieldMaskPlayerQuestLog16_3"}
	FieldMaskPlayerQuestLog16_5                = FieldMask{Size: 1, Offset: 0xED, name: "FieldMaskPlayerQuestLog16_5"}
	FieldMaskPlayerQuestLog17_1                = FieldMask{Size: 1, Offset: 0xEE, name: "FieldMaskPlayerQuestLog17_1"}
	FieldMaskPlayerQuestLog17_2                = FieldMask{Size: 1, Offset: 0xEF, name: "FieldMaskPlayerQuestLog17_2"}
	FieldMaskPlayerQuestLog17_3                = FieldMask{Size: 2, Offset: 0xF0, name: "FieldMaskPlayerQuestLog17_3"}
	FieldMaskPlayerQuestLog17_5                = FieldMask{Size: 1, Offset: 0xF2, name: "FieldMaskPlayerQuestLog17_5"}
	FieldMaskPlayerQuestLog18_1                = FieldMask{Size: 1, Offset: 0xF3, name: "FieldMaskPlayerQuestLog18_1"}
	FieldMaskPlayerQuestLog18_2                = FieldMask{Size: 1, Offset: 0xF4, name: "FieldMaskPlayerQuestLog18_2"}
	FieldMaskPlayerQuestLog18_3                = FieldMask{Size: 2, Offset: 0xF5, name: "FieldMaskPlayerQuestLog18_3"}
	FieldMaskPlayerQuestLog18_5                = FieldMask{Size: 1, Offset: 0xF7, name: "FieldMaskPlayerQuestLog18_5"}
	FieldMaskPlayerQuestLog19_1                = FieldMask{Size: 1, Offset: 0xF8, name: "FieldMaskPlayerQuestLog19_1"}
	FieldMaskPlayerQuestLog19_2                = FieldMask{Size: 1, Offset: 0xF9, name: "FieldMaskPlayerQuestLog19_2"}
	FieldMaskPlayerQuestLog19_3                = FieldMask{Size: 2, Offset: 0xFA, name: "FieldMaskPlayerQuestLog19_3"}
	FieldMaskPlayerQuestLog19_5                = FieldMask{Size: 1, Offset: 0xFC, name: "FieldMaskPlayerQuestLog19_5"}
	FieldMaskPlayerQuestLog20_1                = FieldMask{Size: 1, Offset: 0xFD, name: "FieldMaskPlayerQuestLog20_1"}
	FieldMaskPlayerQuestLog20_2                = FieldMask{Size: 1, Offset: 0xFE, name: "FieldMaskPlayerQuestLog20_2"}
	FieldMaskPlayerQuestLog20_3                = FieldMask{Size: 2, Offset: 0xFF, name: "FieldMaskPlayerQuestLog20_3"}
	FieldMaskPlayerQuestLog20_5                = FieldMask{Size: 1, Offset: 0x101, name: "FieldMaskPlayerQuestLog20_5"}
	FieldMaskPlayerQuestLog21_1                = FieldMask{Size: 1, Offset: 0x102, name: "FieldMaskPlayerQuestLog21_1"}
	FieldMaskPlayerQuestLog21_2                = FieldMask{Size: 1, Offset: 0x103, name: "FieldMaskPlayerQuestLog21_2"}
	FieldMaskPlayerQuestLog21_3                = FieldMask{Size: 2, Offset: 0x104, name: "FieldMaskPlayerQuestLog21_3"}
	FieldMaskPlayerQuestLog21_5                = FieldMask{Size: 1, Offset: 0x106, name: "FieldMaskPlayerQuestLog21_5"}
	FieldMaskPlayerQuestLog22_1                = FieldMask{Size: 1, Offset: 0x107, name: "FieldMaskPlayerQuestLog22_1"}
	FieldMaskPlayerQuestLog22_2                = FieldMask{Size: 1, Offset: 0x108, name: "FieldMaskPlayerQuestLog22_2"}
	FieldMaskPlayerQuestLog22_3                = FieldMask{Size: 2, Offset: 0x109, name: "FieldMaskPlayerQuestLog22_3"}
	FieldMaskPlayerQuestLog22_5                = FieldMask{Size: 1, Offset: 0x10B, name: "FieldMaskPlayerQuestLog22_5"}
	FieldMaskPlayerQuestLog23_1                = FieldMask{Size: 1, Offset: 0x10C, name: "FieldMaskPlayerQuestLog23_1"}
	FieldMaskPlayerQuestLog23_2                = FieldMask{Size: 1, Offset: 0x10D, name: "FieldMaskPlayerQuestLog23_2"}
	FieldMaskPlayerQuestLog23_3                = FieldMask{Size: 2, Offset: 0x10E, name: "FieldMaskPlayerQuestLog23_3"}
	FieldMaskPlayerQuestLog23_5                = FieldMask{Size: 1, Offset: 0x110, name: "FieldMaskPlayerQuestLog23_5"}
	FieldMaskPlayerQuestLog24_1                = FieldMask{Size: 1, Offset: 0x111, name: "FieldMaskPlayerQuestLog24_1"}
	FieldMaskPlayerQuestLog24_2                = FieldMask{Size: 1, Offset: 0x112, name: "FieldMaskPlayerQuestLog24_2"}
	FieldMaskPlayerQuestLog24_3                = FieldMask{Size: 2, Offset: 0x113, name: "FieldMaskPlayerQuestLog24_3"}
	FieldMaskPlayerQuestLog24_5                = FieldMask{Size: 1, Offset: 0x115, name: "FieldMaskPlayerQuestLog24_5"}
	FieldMaskPlayerQuestLog25_1                = FieldMask{Size: 1, Offset: 0x116, name: "FieldMaskPlayerQuestLog25_1"}
	FieldMaskPlayerQuestLog25_2                = FieldMask{Size: 1, Offset: 0x117, name: "FieldMaskPlayerQuestLog25_2"}
	FieldMaskPlayerQuestLog25_3                = FieldMask{Size: 2, Offset: 0x118, name: "FieldMaskPlayerQuestLog25_3"}
	FieldMaskPlayerQuestLog25_5                = FieldMask{Size: 1, Offset: 0x11A, name: "FieldMaskPlayerQuestLog25_5"}
	FieldMaskPlayerVisibleItem                 = FieldMask{Size: 3, Offset: 0x11B, name: "FieldMaskPlayerVisibleItem"}
	FieldMaskPlayerChosenTitle                 = FieldMask{Size: 1, Offset: 0x141, name: "FieldMaskPlayerChosenTitle"}
	FieldMaskPlayerFakeInebriation             = FieldMask{Size: 1, Offset: 0x142, name: "FieldMaskPlayerFakeInebriation"}
	FieldMaskPlayerFieldInv                    = FieldMask{Size: 3, Offset: 0x144, name: "FieldMaskPlayerFieldInv"}
	FieldMaskPlayerFarsight                    = FieldMask{Size: 2, Offset: 0x270, name: "FieldMaskPlayerFarsight"}
	FieldMaskPlayerKnownTitles                 = FieldMask{Size: 2, Offset: 0x272, name: "FieldMaskPlayerKnownTitles"}
	FieldMaskPlayerKnownTitles1                = FieldMask{Size: 2, Offset: 0x274, name: "FieldMaskPlayerKnownTitles1"}
	FieldMaskPlayerKnownTitles2                = FieldMask{Size: 2, Offset: 0x276, name: "FieldMaskPlayerKnownTitles2"}
	FieldMaskPlayerKnownCurrencies             = FieldMask{Size: 2, Offset: 0x278, name: "FieldMaskPlayerKnownCurrencies"}
	FieldMaskPlayerXp                          = FieldMask{Size: 1, Offset: 0x27A, name: "FieldMaskPlayerXp"}
	FieldMaskPlayerNextLevelXp                 = FieldMask{Size: 1, Offset: 0x27B, name: "FieldMaskPlayerNextLevelXp"}
	FieldMaskPlayerSkillInfo                   = FieldMask{Size: 3, Offset: 0x27C, name: "FieldMaskPlayerSkillInfo"}
	FieldMaskPlayerCharacterPoints1            = FieldMask{Size: 1, Offset: 0x3FC, name: "FieldMaskPlayerCharacterPoints1"}
	FieldMaskPlayerCharacterPoints2            = FieldMask{Size: 1, Offset: 0x3FD, name: "FieldMaskPlayerCharacterPoints2"}
	FieldMaskPlayerTrackCreatures              = FieldMask{Size: 1, Offset: 0x3FE, name: "FieldMaskPlayerTrackCreatures"}
	FieldMaskPlayerTrackResources              = FieldMask{Size: 1, Offset: 0x3FF, name: "FieldMaskPlayerTrackResources"}
	FieldMaskPlayerBlockPercentage             = FieldMask{Size: 1, Offset: 0x400, name: "FieldMaskPlayerBlockPercentage"}
	FieldMaskPlayerDodgePercentage             = FieldMask{Size: 1, Offset: 0x401, name: "FieldMaskPlayerDodgePercentage"}
	FieldMaskPlayerParryPercentage             = FieldMask{Size: 1, Offset: 0x402, name: "FieldMaskPlayerParryPercentage"}
	FieldMaskPlayerExpertise                   = FieldMask{Size: 1, Offset: 0x403, name: "FieldMaskPlayerExpertise"}
	FieldMaskPlayerOffhandExpertise            = FieldMask{Size: 1, Offset: 0x404, name: "FieldMaskPlayerOffhandExpertise"}
	FieldMaskPlayerCritPercentage              = FieldMask{Size: 1, Offset: 0x405, name: "FieldMaskPlayerCritPercentage"}
	FieldMaskPlayerRangedCritPercentage        = FieldMask{Size: 1, Offset: 0x406, name: "FieldMaskPlayerRangedCritPercentage"}
	FieldMaskPlayerOffhandCritPercentage       = FieldMask{Size: 1, Offset: 0x407, name: "FieldMaskPlayerOffhandCritPercentage"}
	FieldMaskPlayerSpellCritPercentage1        = FieldMask{Size: 7, Offset: 0x408, name: "FieldMaskPlayerSpellCritPercentage1"}
	FieldMaskPlayerShieldBlock                 = FieldMask{Size: 1, Offset: 0x40F, name: "FieldMaskPlayerShieldBlock"}
	FieldMaskPlayerShieldBlockCritPercentage   = FieldMask{Size: 1, Offset: 0x410, name: "FieldMaskPlayerShieldBlockCritPercentage"}
	FieldMaskPlayerExploredZones1              = FieldMask{Size: 1, Offset: 0x411, name: "FieldMaskPlayerExploredZones1"}
	FieldMaskPlayerRestStateExperience         = FieldMask{Size: 1, Offset: 0x491, name: "FieldMaskPlayerRestStateExperience"}
	FieldMaskPlayerCoinage                     = FieldMask{Size: 1, Offset: 0x492, name: "FieldMaskPlayerCoinage"}
	FieldMaskPlayerModDamageDonePos            = FieldMask{Size: 7, Offset: 0x493, name: "FieldMaskPlayerModDamageDonePos"}
	FieldMaskPlayerModDamageDoneNeg            = FieldMask{Size: 7, Offset: 0x49A, name: "FieldMaskPlayerModDamageDoneNeg"}
	FieldMaskPlayerModDamageDonePct            = FieldMask{Size: 7, Offset: 0x4A1, name: "FieldMaskPlayerModDamageDonePct"}
	FieldMaskPlayerModHealingDonePos           = FieldMask{Size: 1, Offset: 0x4A8, name: "FieldMaskPlayerModHealingDonePos"}
	FieldMaskPlayerModHealingPct               = FieldMask{Size: 1, Offset: 0x4A9, name: "FieldMaskPlayerModHealingPct"}
	FieldMaskPlayerModHealingDonePct           = FieldMask{Size: 1, Offset: 0x4AA, name: "FieldMaskPlayerModHealingDonePct"}
	FieldMaskPlayerModTargetResistance         = FieldMask{Size: 1, Offset: 0x4AB, name: "FieldMaskPlayerModTargetResistance"}
	FieldMaskPlayerModTargetPhysicalResistance = FieldMask{Size: 1, Offset: 0x4AC, name: "FieldMaskPlayerModTargetPhysicalResistance"}
	FieldMaskPlayerFeatures                    = FieldMask{Size: 1, Offset: 0x4AD, name: "FieldMaskPlayerFeatures"}
	FieldMaskPlayerAmmoId                      = FieldMask{Size: 1, Offset: 0x4AE, name: "FieldMaskPlayerAmmoId"}
	FieldMaskPlayerSelfResSpell                = FieldMask{Size: 1, Offset: 0x4AF, name: "FieldMaskPlayerSelfResSpell"}
	FieldMaskPlayerPvpMedals                   = FieldMask{Size: 1, Offset: 0x4B0, name: "FieldMaskPlayerPvpMedals"}
	FieldMaskPlayerBuybackPrice1               = FieldMask{Size: 1, Offset: 0x4B1, name: "FieldMaskPlayerBuybackPrice1"}
	FieldMaskPlayerBuybackTimestamp1           = FieldMask{Size: 1, Offset: 0x4BD, name: "FieldMaskPlayerBuybackTimestamp1"}
	FieldMaskPlayerKills                       = FieldMask{Size: 1, Offset: 0x4C9, name: "FieldMaskPlayerKills"}
	FieldMaskPlayerTodayContribution           = FieldMask{Size: 1, Offset: 0x4CA, name: "FieldMaskPlayerTodayContribution"}
	FieldMaskPlayerYesterdayContribution       = FieldMask{Size: 1, Offset: 0x4CB, name: "FieldMaskPlayerYesterdayContribution"}
	FieldMaskPlayerLifetimeHonorableKills      = FieldMask{Size: 1, Offset: 0x4CC, name: "FieldMaskPlayerLifetimeHonorableKills"}
	FieldMaskPlayerBytes4                      = FieldMask{Size: 1, Offset: 0x4CD, name: "FieldMaskPlayerBytes4"}
	FieldMaskPlayerWatchedFactionIndex         = FieldMask{Size: 1, Offset: 0x4CE, name: "FieldMaskPlayerWatchedFactionIndex"}
	FieldMaskPlayerCombatRating1               = FieldMask{Size: 2, Offset: 0x4CF, name: "FieldMaskPlayerCombatRating1"}
	FieldMaskPlayerArenaTeamInfo11             = FieldMask{Size: 2, Offset: 0x4E8, name: "FieldMaskPlayerArenaTeamInfo11"}
	FieldMaskPlayerHonorCurrency               = FieldMask{Size: 1, Offset: 0x4FD, name: "FieldMaskPlayerHonorCurrency"}
	FieldMaskPlayerArenaCurrency               = FieldMask{Size: 1, Offset: 0x4FE, name: "FieldMaskPlayerArenaCurrency"}
	FieldMaskPlayerMaxLevel                    = FieldMask{Size: 1, Offset: 0x4FF, name: "FieldMaskPlayerMaxLevel"}
	FieldMaskPlayerDailyQuests1                = FieldMask{Size: 2, Offset: 0x500, name: "FieldMaskPlayerDailyQuests1"}
	FieldMaskPlayerRuneRegen1                  = FieldMask{Size: 4, Offset: 0x519, name: "FieldMaskPlayerRuneRegen1"}
	FieldMaskPlayerNoReagentCost1              = FieldMask{Size: 3, Offset: 0x51D, name: "FieldMaskPlayerNoReagentCost1"}
	FieldMaskPlayerGlyphSlots1                 = FieldMask{Size: 6, Offset: 0x520, name: "FieldMaskPlayerGlyphSlots1"}
	FieldMaskPlayerGlyphs1                     = FieldMask{Size: 6, Offset: 0x526, name: "FieldMaskPlayerGlyphs1"}
	FieldMaskPlayerGlyphsEnabled               = FieldMask{Size: 1, Offset: 0x52C, name: "FieldMaskPlayerGlyphsEnabled"}
	FieldMaskPlayerPetSpellPower               = FieldMask{Size: 1, Offset: 0x52D, name: "FieldMaskPlayerPetSpellPower"}

	FieldMaskGameObjectDisplayid      = FieldMask{Size: 1, Offset: 0x8, name: "FieldMaskGameObjectDisplayid"}
	FieldMaskGameObjectFlags          = FieldMask{Size: 1, Offset: 0x9, name: "FieldMaskGameObjectFlags"}
	FieldMaskGameObjectParentrotation = FieldMask{Size: 4, Offset: 0xA, name: "FieldMaskGameObjectParentrotation"}
	FieldMaskGameObjectDynamic        = FieldMask{Size: 1, Offset: 0xE, name: "FieldMaskGameObjectDynamic"}
	FieldMaskGameObjectFaction        = FieldMask{Size: 1, Offset: 0xF, name: "FieldMaskGameObjectFaction"}
	FieldMaskGameObjectLevel          = FieldMask{Size: 1, Offset: 0x10, name: "FieldMaskGameObjectLevel"}
	FieldMaskGameObjectBytes1         = FieldMask{Size: 1, Offset: 0x11, name: "FieldMaskGameObjectBytes1"}

	FieldMaskDynamicObjectCaster   = FieldMask{Size: 2, Offset: 0x6, name: "FieldMaskDynamicObjectCaster"}
	FieldMaskDynamicObjectBytes    = FieldMask{Size: 1, Offset: 0x8, name: "FieldMaskDynamicObjectBytes"}
	FieldMaskDynamicObjectSpellid  = FieldMask{Size: 1, Offset: 0x9, name: "FieldMaskDynamicObjectSpellid"}
	FieldMaskDynamicObjectRadius   = FieldMask{Size: 1, Offset: 0xA, name: "FieldMaskDynamicObjectRadius"}
	FieldMaskDynamicObjectCasttime = FieldMask{Size: 1, Offset: 0xB, name: "FieldMaskDynamicObjectCasttime"}

	FieldMaskCorpseOwner        = FieldMask{Size: 2, Offset: 0x6, name: "FieldMaskCorpseOwner"}
	FieldMaskCorpseParty        = FieldMask{Size: 2, Offset: 0x8, name: "FieldMaskCorpseParty"}
	FieldMaskCorpseDisplayId    = FieldMask{Size: 1, Offset: 0xA, name: "FieldMaskCorpseDisplayId"}
	FieldMaskCorpseItem         = FieldMask{Size: 1, Offset: 0xB, name: "FieldMaskCorpseItem"}
	FieldMaskCorpseBytes1       = FieldMask{Size: 1, Offset: 0x1E, name: "FieldMaskCorpseBytes1"}
	FieldMaskCorpseBytes2       = FieldMask{Size: 1, Offset: 0x1F, name: "FieldMaskCorpseBytes2"}
	FieldMaskCorpseGuild        = FieldMask{Size: 1, Offset: 0x20, name: "FieldMaskCorpseGuild"}
	FieldMaskCorpseFlags        = FieldMask{Size: 1, Offset: 0x21, name: "FieldMaskCorpseFlags"}
	FieldMaskCorpseDynamicFlags = FieldMask{Size: 1, Offset: 0x22, name: "FieldMaskCorpseDynamicFlags"}
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
