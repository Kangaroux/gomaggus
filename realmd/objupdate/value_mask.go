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
}

var (
	FieldMaskObjectGuid    = FieldMask{Size: 2, Offset: 0x0}
	FieldMaskObjectType    = FieldMask{Size: 1, Offset: 0x2}
	FieldMaskObjectEntry   = FieldMask{Size: 1, Offset: 0x3}
	FieldMaskObjectScaleX  = FieldMask{Size: 1, Offset: 0x4}
	FieldMaskObjectPadding = FieldMask{Size: 1, Offset: 0x5}

	FieldMaskItemOwner              = FieldMask{Size: 2, Offset: 0x6}
	FieldMaskItemContained          = FieldMask{Size: 2, Offset: 0x8}
	FieldMaskItemCreator            = FieldMask{Size: 2, Offset: 0xA}
	FieldMaskItemGiftcreator        = FieldMask{Size: 2, Offset: 0xC}
	FieldMaskItemStackCount         = FieldMask{Size: 1, Offset: 0xE}
	FieldMaskItemDuration           = FieldMask{Size: 1, Offset: 0xF}
	FieldMaskItemSpellCharges       = FieldMask{Size: 5, Offset: 0x10}
	FieldMaskItemFlags              = FieldMask{Size: 1, Offset: 0x15}
	FieldMaskItemEnchantment1_1     = FieldMask{Size: 2, Offset: 0x16}
	FieldMaskItemEnchantment1_3     = FieldMask{Size: 1, Offset: 0x18}
	FieldMaskItemEnchantment2_1     = FieldMask{Size: 2, Offset: 0x19}
	FieldMaskItemEnchantment2_3     = FieldMask{Size: 1, Offset: 0x1B}
	FieldMaskItemEnchantment3_1     = FieldMask{Size: 2, Offset: 0x1C}
	FieldMaskItemEnchantment3_3     = FieldMask{Size: 1, Offset: 0x1E}
	FieldMaskItemEnchantment4_1     = FieldMask{Size: 2, Offset: 0x1F}
	FieldMaskItemEnchantment4_3     = FieldMask{Size: 1, Offset: 0x21}
	FieldMaskItemEnchantment5_1     = FieldMask{Size: 2, Offset: 0x22}
	FieldMaskItemEnchantment5_3     = FieldMask{Size: 1, Offset: 0x24}
	FieldMaskItemEnchantment6_1     = FieldMask{Size: 2, Offset: 0x25}
	FieldMaskItemEnchantment6_3     = FieldMask{Size: 1, Offset: 0x27}
	FieldMaskItemEnchantment7_1     = FieldMask{Size: 2, Offset: 0x28}
	FieldMaskItemEnchantment7_3     = FieldMask{Size: 1, Offset: 0x2A}
	FieldMaskItemEnchantment8_1     = FieldMask{Size: 2, Offset: 0x2B}
	FieldMaskItemEnchantment8_3     = FieldMask{Size: 1, Offset: 0x2D}
	FieldMaskItemEnchantment9_1     = FieldMask{Size: 2, Offset: 0x2E}
	FieldMaskItemEnchantment9_3     = FieldMask{Size: 1, Offset: 0x30}
	FieldMaskItemEnchantment10_1    = FieldMask{Size: 2, Offset: 0x31}
	FieldMaskItemEnchantment10_3    = FieldMask{Size: 1, Offset: 0x33}
	FieldMaskItemEnchantment11_1    = FieldMask{Size: 2, Offset: 0x34}
	FieldMaskItemEnchantment11_3    = FieldMask{Size: 1, Offset: 0x36}
	FieldMaskItemEnchantment12_1    = FieldMask{Size: 2, Offset: 0x37}
	FieldMaskItemEnchantment12_3    = FieldMask{Size: 1, Offset: 0x39}
	FieldMaskItemPropertySeed       = FieldMask{Size: 1, Offset: 0x3A}
	FieldMaskItemRandomPropertiesId = FieldMask{Size: 1, Offset: 0x3B}
	FieldMaskItemDurability         = FieldMask{Size: 1, Offset: 0x3C}
	FieldMaskItemMaxdurability      = FieldMask{Size: 1, Offset: 0x3D}
	FieldMaskItemCreatePlayedTime   = FieldMask{Size: 1, Offset: 0x3E}

	FieldMaskContainerNumSlots = FieldMask{Size: 1, Offset: 0x40}
	FieldMaskContainerSlot1    = FieldMask{Size: 7, Offset: 0x42}

	FieldMaskUnitCharm                             = FieldMask{Size: 2, Offset: 0x6}
	FieldMaskUnitSummon                            = FieldMask{Size: 2, Offset: 0x8}
	FieldMaskUnitCritter                           = FieldMask{Size: 2, Offset: 0xA}
	FieldMaskUnitCharmedby                         = FieldMask{Size: 2, Offset: 0xC}
	FieldMaskUnitSummonedby                        = FieldMask{Size: 2, Offset: 0xE}
	FieldMaskUnitCreatedby                         = FieldMask{Size: 2, Offset: 0x10}
	FieldMaskUnitTarget                            = FieldMask{Size: 2, Offset: 0x12}
	FieldMaskUnitChannelObject                     = FieldMask{Size: 2, Offset: 0x14}
	FieldMaskUnitChannelSpell                      = FieldMask{Size: 1, Offset: 0x16}
	FieldMaskUnitBytes0                            = FieldMask{Size: 1, Offset: 0x17}
	FieldMaskUnitHealth                            = FieldMask{Size: 1, Offset: 0x18}
	FieldMaskUnitPower1                            = FieldMask{Size: 1, Offset: 0x19}
	FieldMaskUnitPower2                            = FieldMask{Size: 1, Offset: 0x1A}
	FieldMaskUnitPower3                            = FieldMask{Size: 1, Offset: 0x1B}
	FieldMaskUnitPower4                            = FieldMask{Size: 1, Offset: 0x1C}
	FieldMaskUnitPower5                            = FieldMask{Size: 1, Offset: 0x1D}
	FieldMaskUnitPower6                            = FieldMask{Size: 1, Offset: 0x1E}
	FieldMaskUnitPower7                            = FieldMask{Size: 1, Offset: 0x1F}
	FieldMaskUnitMaxHealth                         = FieldMask{Size: 1, Offset: 0x20}
	FieldMaskUnitMaxpower1                         = FieldMask{Size: 1, Offset: 0x21}
	FieldMaskUnitMaxpower2                         = FieldMask{Size: 1, Offset: 0x22}
	FieldMaskUnitMaxpower3                         = FieldMask{Size: 1, Offset: 0x23}
	FieldMaskUnitMaxpower4                         = FieldMask{Size: 1, Offset: 0x24}
	FieldMaskUnitMaxpower5                         = FieldMask{Size: 1, Offset: 0x25}
	FieldMaskUnitMaxpower6                         = FieldMask{Size: 1, Offset: 0x26}
	FieldMaskUnitMaxpower7                         = FieldMask{Size: 1, Offset: 0x27}
	FieldMaskUnitPowerRegenFlatModifier            = FieldMask{Size: 7, Offset: 0x28}
	FieldMaskUnitPowerRegenInterruptedFlatModifier = FieldMask{Size: 7, Offset: 0x2F}
	FieldMaskUnitLevel                             = FieldMask{Size: 1, Offset: 0x36}
	FieldMaskUnitFactionTemplate                   = FieldMask{Size: 1, Offset: 0x37}
	FieldMaskUnitVirtualItemSlotId                 = FieldMask{Size: 3, Offset: 0x38}
	FieldMaskUnitFlags                             = FieldMask{Size: 1, Offset: 0x3B}
	FieldMaskUnitFlags2                            = FieldMask{Size: 1, Offset: 0x3C}
	FieldMaskUnitAurastate                         = FieldMask{Size: 1, Offset: 0x3D}
	FieldMaskUnitBaseattacktime                    = FieldMask{Size: 2, Offset: 0x3E}
	FieldMaskUnitRangedattacktime                  = FieldMask{Size: 1, Offset: 0x40}
	FieldMaskUnitBoundingradius                    = FieldMask{Size: 1, Offset: 0x41}
	FieldMaskUnitCombatreach                       = FieldMask{Size: 1, Offset: 0x42}
	FieldMaskUnitDisplayId                         = FieldMask{Size: 1, Offset: 0x43}
	FieldMaskUnitNativeDisplayId                   = FieldMask{Size: 1, Offset: 0x44}
	FieldMaskUnitMountdisplayid                    = FieldMask{Size: 1, Offset: 0x45}
	FieldMaskUnitMindamage                         = FieldMask{Size: 1, Offset: 0x46}
	FieldMaskUnitMaxdamage                         = FieldMask{Size: 1, Offset: 0x47}
	FieldMaskUnitMinoffhanddamage                  = FieldMask{Size: 1, Offset: 0x48}
	FieldMaskUnitMaxoffhanddamage                  = FieldMask{Size: 1, Offset: 0x49}
	FieldMaskUnitBytes1                            = FieldMask{Size: 1, Offset: 0x4A}
	FieldMaskUnitPetnumber                         = FieldMask{Size: 1, Offset: 0x4B}
	FieldMaskUnitPetNameTimestamp                  = FieldMask{Size: 1, Offset: 0x4C}
	FieldMaskUnitPetexperience                     = FieldMask{Size: 1, Offset: 0x4D}
	FieldMaskUnitPetnextlevelexp                   = FieldMask{Size: 1, Offset: 0x4E}
	FieldMaskUnitDynamicFlags                      = FieldMask{Size: 1, Offset: 0x4F}
	FieldMaskUnitModCastSpeed                      = FieldMask{Size: 1, Offset: 0x50}
	FieldMaskUnitCreatedBySpell                    = FieldMask{Size: 1, Offset: 0x51}
	FieldMaskUnitNpcFlags                          = FieldMask{Size: 1, Offset: 0x52}
	FieldMaskUnitNpcEmotestate                     = FieldMask{Size: 1, Offset: 0x53}
	FieldMaskUnitStrength                          = FieldMask{Size: 1, Offset: 0x54}
	FieldMaskUnitAgility                           = FieldMask{Size: 1, Offset: 0x55}
	FieldMaskUnitStamina                           = FieldMask{Size: 1, Offset: 0x56}
	FieldMaskUnitIntellect                         = FieldMask{Size: 1, Offset: 0x57}
	FieldMaskUnitSpirit                            = FieldMask{Size: 1, Offset: 0x58}
	FieldMaskUnitPosStat0                          = FieldMask{Size: 1, Offset: 0x59}
	FieldMaskUnitPosStat1                          = FieldMask{Size: 1, Offset: 0x5A}
	FieldMaskUnitPosStat2                          = FieldMask{Size: 1, Offset: 0x5B}
	FieldMaskUnitPosStat3                          = FieldMask{Size: 1, Offset: 0x5C}
	FieldMaskUnitPosStat4                          = FieldMask{Size: 1, Offset: 0x5D}
	FieldMaskUnitNegStat0                          = FieldMask{Size: 1, Offset: 0x5E}
	FieldMaskUnitNegStat1                          = FieldMask{Size: 1, Offset: 0x5F}
	FieldMaskUnitNegStat2                          = FieldMask{Size: 1, Offset: 0x60}
	FieldMaskUnitNegStat3                          = FieldMask{Size: 1, Offset: 0x61}
	FieldMaskUnitNegStat4                          = FieldMask{Size: 1, Offset: 0x62}
	FieldMaskUnitResistances                       = FieldMask{Size: 7, Offset: 0x63}
	FieldMaskUnitResistanceBuffModsPositive        = FieldMask{Size: 7, Offset: 0x6A}
	FieldMaskUnitResistanceBuffModsNegative        = FieldMask{Size: 7, Offset: 0x71}
	FieldMaskUnitBaseMana                          = FieldMask{Size: 1, Offset: 0x78}
	FieldMaskUnitBaseHealth                        = FieldMask{Size: 1, Offset: 0x79}
	FieldMaskUnitBytes2                            = FieldMask{Size: 1, Offset: 0x7A}
	FieldMaskUnitAttackPower                       = FieldMask{Size: 1, Offset: 0x7B}
	FieldMaskUnitAttackPowerMods                   = FieldMask{Size: 1, Offset: 0x7C}
	FieldMaskUnitAttackPowerMultiplier             = FieldMask{Size: 1, Offset: 0x7D}
	FieldMaskUnitRangedAttackPower                 = FieldMask{Size: 1, Offset: 0x7E}
	FieldMaskUnitRangedAttackPowerMods             = FieldMask{Size: 1, Offset: 0x7F}
	FieldMaskUnitRangedAttackPowerMultiplier       = FieldMask{Size: 1, Offset: 0x80}
	FieldMaskUnitMinRangedDamage                   = FieldMask{Size: 1, Offset: 0x81}
	FieldMaskUnitMaxRangedDamage                   = FieldMask{Size: 1, Offset: 0x82}
	FieldMaskUnitPowerCostModifier                 = FieldMask{Size: 7, Offset: 0x83}
	FieldMaskUnitPowerCostMultiplier               = FieldMask{Size: 7, Offset: 0x8A}
	FieldMaskUnitMaxHealthModifier                 = FieldMask{Size: 1, Offset: 0x91}
	FieldMaskUnitHoverHeight                       = FieldMask{Size: 1, Offset: 0x92}

	FieldMaskPlayerDuelArbiter                 = FieldMask{Size: 2, Offset: 0x94}
	FieldMaskPlayerFlags                       = FieldMask{Size: 1, Offset: 0x96}
	FieldMaskPlayerGuildid                     = FieldMask{Size: 1, Offset: 0x97}
	FieldMaskPlayerGuildrank                   = FieldMask{Size: 1, Offset: 0x98}
	FieldMaskPlayerFieldBytes                  = FieldMask{Size: 1, Offset: 0x99}
	FieldMaskPlayerBytes2                      = FieldMask{Size: 1, Offset: 0x9A}
	FieldMaskPlayerBytes3                      = FieldMask{Size: 1, Offset: 0x9B}
	FieldMaskPlayerDuelTeam                    = FieldMask{Size: 1, Offset: 0x9C}
	FieldMaskPlayerGuildTimestamp              = FieldMask{Size: 1, Offset: 0x9D}
	FieldMaskPlayerQuestLog1_1                 = FieldMask{Size: 1, Offset: 0x9E}
	FieldMaskPlayerQuestLog1_2                 = FieldMask{Size: 1, Offset: 0x9F}
	FieldMaskPlayerQuestLog1_3                 = FieldMask{Size: 2, Offset: 0xA0}
	FieldMaskPlayerQuestLog1_4                 = FieldMask{Size: 1, Offset: 0xA2}
	FieldMaskPlayerQuestLog2_1                 = FieldMask{Size: 1, Offset: 0xA3}
	FieldMaskPlayerQuestLog2_2                 = FieldMask{Size: 1, Offset: 0xA4}
	FieldMaskPlayerQuestLog2_3                 = FieldMask{Size: 2, Offset: 0xA5}
	FieldMaskPlayerQuestLog2_5                 = FieldMask{Size: 1, Offset: 0xA7}
	FieldMaskPlayerQuestLog3_1                 = FieldMask{Size: 1, Offset: 0xA8}
	FieldMaskPlayerQuestLog3_2                 = FieldMask{Size: 1, Offset: 0xA9}
	FieldMaskPlayerQuestLog3_3                 = FieldMask{Size: 2, Offset: 0xAA}
	FieldMaskPlayerQuestLog3_5                 = FieldMask{Size: 1, Offset: 0xAC}
	FieldMaskPlayerQuestLog4_1                 = FieldMask{Size: 1, Offset: 0xAD}
	FieldMaskPlayerQuestLog4_2                 = FieldMask{Size: 1, Offset: 0xAE}
	FieldMaskPlayerQuestLog4_3                 = FieldMask{Size: 2, Offset: 0xAF}
	FieldMaskPlayerQuestLog4_5                 = FieldMask{Size: 1, Offset: 0xB1}
	FieldMaskPlayerQuestLog5_1                 = FieldMask{Size: 1, Offset: 0xB2}
	FieldMaskPlayerQuestLog5_2                 = FieldMask{Size: 1, Offset: 0xB3}
	FieldMaskPlayerQuestLog5_3                 = FieldMask{Size: 2, Offset: 0xB4}
	FieldMaskPlayerQuestLog5_5                 = FieldMask{Size: 1, Offset: 0xB6}
	FieldMaskPlayerQuestLog6_1                 = FieldMask{Size: 1, Offset: 0xB7}
	FieldMaskPlayerQuestLog6_2                 = FieldMask{Size: 1, Offset: 0xB8}
	FieldMaskPlayerQuestLog6_3                 = FieldMask{Size: 2, Offset: 0xB9}
	FieldMaskPlayerQuestLog6_5                 = FieldMask{Size: 1, Offset: 0xBB}
	FieldMaskPlayerQuestLog7_1                 = FieldMask{Size: 1, Offset: 0xBC}
	FieldMaskPlayerQuestLog7_2                 = FieldMask{Size: 1, Offset: 0xBD}
	FieldMaskPlayerQuestLog7_3                 = FieldMask{Size: 2, Offset: 0xBE}
	FieldMaskPlayerQuestLog7_5                 = FieldMask{Size: 1, Offset: 0xC0}
	FieldMaskPlayerQuestLog8_1                 = FieldMask{Size: 1, Offset: 0xC1}
	FieldMaskPlayerQuestLog8_2                 = FieldMask{Size: 1, Offset: 0xC2}
	FieldMaskPlayerQuestLog8_3                 = FieldMask{Size: 2, Offset: 0xC3}
	FieldMaskPlayerQuestLog8_5                 = FieldMask{Size: 1, Offset: 0xC5}
	FieldMaskPlayerQuestLog9_1                 = FieldMask{Size: 1, Offset: 0xC6}
	FieldMaskPlayerQuestLog9_2                 = FieldMask{Size: 1, Offset: 0xC7}
	FieldMaskPlayerQuestLog9_3                 = FieldMask{Size: 2, Offset: 0xC8}
	FieldMaskPlayerQuestLog9_5                 = FieldMask{Size: 1, Offset: 0xCA}
	FieldMaskPlayerQuestLog10_1                = FieldMask{Size: 1, Offset: 0xCB}
	FieldMaskPlayerQuestLog10_2                = FieldMask{Size: 1, Offset: 0xCC}
	FieldMaskPlayerQuestLog10_3                = FieldMask{Size: 2, Offset: 0xCD}
	FieldMaskPlayerQuestLog10_5                = FieldMask{Size: 1, Offset: 0xCF}
	FieldMaskPlayerQuestLog11_1                = FieldMask{Size: 1, Offset: 0xD0}
	FieldMaskPlayerQuestLog11_2                = FieldMask{Size: 1, Offset: 0xD1}
	FieldMaskPlayerQuestLog11_3                = FieldMask{Size: 2, Offset: 0xD2}
	FieldMaskPlayerQuestLog11_5                = FieldMask{Size: 1, Offset: 0xD4}
	FieldMaskPlayerQuestLog12_1                = FieldMask{Size: 1, Offset: 0xD5}
	FieldMaskPlayerQuestLog12_2                = FieldMask{Size: 1, Offset: 0xD6}
	FieldMaskPlayerQuestLog12_3                = FieldMask{Size: 2, Offset: 0xD7}
	FieldMaskPlayerQuestLog12_5                = FieldMask{Size: 1, Offset: 0xD9}
	FieldMaskPlayerQuestLog13_1                = FieldMask{Size: 1, Offset: 0xDA}
	FieldMaskPlayerQuestLog13_2                = FieldMask{Size: 1, Offset: 0xDB}
	FieldMaskPlayerQuestLog13_3                = FieldMask{Size: 2, Offset: 0xDC}
	FieldMaskPlayerQuestLog13_5                = FieldMask{Size: 1, Offset: 0xDE}
	FieldMaskPlayerQuestLog14_1                = FieldMask{Size: 1, Offset: 0xDF}
	FieldMaskPlayerQuestLog14_2                = FieldMask{Size: 1, Offset: 0xE0}
	FieldMaskPlayerQuestLog14_3                = FieldMask{Size: 2, Offset: 0xE1}
	FieldMaskPlayerQuestLog14_5                = FieldMask{Size: 1, Offset: 0xE3}
	FieldMaskPlayerQuestLog15_1                = FieldMask{Size: 1, Offset: 0xE4}
	FieldMaskPlayerQuestLog15_2                = FieldMask{Size: 1, Offset: 0xE5}
	FieldMaskPlayerQuestLog15_3                = FieldMask{Size: 2, Offset: 0xE6}
	FieldMaskPlayerQuestLog15_5                = FieldMask{Size: 1, Offset: 0xE8}
	FieldMaskPlayerQuestLog16_1                = FieldMask{Size: 1, Offset: 0xE9}
	FieldMaskPlayerQuestLog16_2                = FieldMask{Size: 1, Offset: 0xEA}
	FieldMaskPlayerQuestLog16_3                = FieldMask{Size: 2, Offset: 0xEB}
	FieldMaskPlayerQuestLog16_5                = FieldMask{Size: 1, Offset: 0xED}
	FieldMaskPlayerQuestLog17_1                = FieldMask{Size: 1, Offset: 0xEE}
	FieldMaskPlayerQuestLog17_2                = FieldMask{Size: 1, Offset: 0xEF}
	FieldMaskPlayerQuestLog17_3                = FieldMask{Size: 2, Offset: 0xF0}
	FieldMaskPlayerQuestLog17_5                = FieldMask{Size: 1, Offset: 0xF2}
	FieldMaskPlayerQuestLog18_1                = FieldMask{Size: 1, Offset: 0xF3}
	FieldMaskPlayerQuestLog18_2                = FieldMask{Size: 1, Offset: 0xF4}
	FieldMaskPlayerQuestLog18_3                = FieldMask{Size: 2, Offset: 0xF5}
	FieldMaskPlayerQuestLog18_5                = FieldMask{Size: 1, Offset: 0xF7}
	FieldMaskPlayerQuestLog19_1                = FieldMask{Size: 1, Offset: 0xF8}
	FieldMaskPlayerQuestLog19_2                = FieldMask{Size: 1, Offset: 0xF9}
	FieldMaskPlayerQuestLog19_3                = FieldMask{Size: 2, Offset: 0xFA}
	FieldMaskPlayerQuestLog19_5                = FieldMask{Size: 1, Offset: 0xFC}
	FieldMaskPlayerQuestLog20_1                = FieldMask{Size: 1, Offset: 0xFD}
	FieldMaskPlayerQuestLog20_2                = FieldMask{Size: 1, Offset: 0xFE}
	FieldMaskPlayerQuestLog20_3                = FieldMask{Size: 2, Offset: 0xFF}
	FieldMaskPlayerQuestLog20_5                = FieldMask{Size: 1, Offset: 0x101}
	FieldMaskPlayerQuestLog21_1                = FieldMask{Size: 1, Offset: 0x102}
	FieldMaskPlayerQuestLog21_2                = FieldMask{Size: 1, Offset: 0x103}
	FieldMaskPlayerQuestLog21_3                = FieldMask{Size: 2, Offset: 0x104}
	FieldMaskPlayerQuestLog21_5                = FieldMask{Size: 1, Offset: 0x106}
	FieldMaskPlayerQuestLog22_1                = FieldMask{Size: 1, Offset: 0x107}
	FieldMaskPlayerQuestLog22_2                = FieldMask{Size: 1, Offset: 0x108}
	FieldMaskPlayerQuestLog22_3                = FieldMask{Size: 2, Offset: 0x109}
	FieldMaskPlayerQuestLog22_5                = FieldMask{Size: 1, Offset: 0x10B}
	FieldMaskPlayerQuestLog23_1                = FieldMask{Size: 1, Offset: 0x10C}
	FieldMaskPlayerQuestLog23_2                = FieldMask{Size: 1, Offset: 0x10D}
	FieldMaskPlayerQuestLog23_3                = FieldMask{Size: 2, Offset: 0x10E}
	FieldMaskPlayerQuestLog23_5                = FieldMask{Size: 1, Offset: 0x110}
	FieldMaskPlayerQuestLog24_1                = FieldMask{Size: 1, Offset: 0x111}
	FieldMaskPlayerQuestLog24_2                = FieldMask{Size: 1, Offset: 0x112}
	FieldMaskPlayerQuestLog24_3                = FieldMask{Size: 2, Offset: 0x113}
	FieldMaskPlayerQuestLog24_5                = FieldMask{Size: 1, Offset: 0x115}
	FieldMaskPlayerQuestLog25_1                = FieldMask{Size: 1, Offset: 0x116}
	FieldMaskPlayerQuestLog25_2                = FieldMask{Size: 1, Offset: 0x117}
	FieldMaskPlayerQuestLog25_3                = FieldMask{Size: 2, Offset: 0x118}
	FieldMaskPlayerQuestLog25_5                = FieldMask{Size: 1, Offset: 0x11A}
	FieldMaskPlayerVisibleItem                 = FieldMask{Size: 3, Offset: 0x11B}
	FieldMaskPlayerChosenTitle                 = FieldMask{Size: 1, Offset: 0x141}
	FieldMaskPlayerFakeInebriation             = FieldMask{Size: 1, Offset: 0x142}
	FieldMaskPlayerFieldInv                    = FieldMask{Size: 3, Offset: 0x144}
	FieldMaskPlayerFarsight                    = FieldMask{Size: 2, Offset: 0x270}
	FieldMaskPlayerKnownTitles                 = FieldMask{Size: 2, Offset: 0x272}
	FieldMaskPlayerKnownTitles1                = FieldMask{Size: 2, Offset: 0x274}
	FieldMaskPlayerKnownTitles2                = FieldMask{Size: 2, Offset: 0x276}
	FieldMaskPlayerKnownCurrencies             = FieldMask{Size: 2, Offset: 0x278}
	FieldMaskPlayerXp                          = FieldMask{Size: 1, Offset: 0x27A}
	FieldMaskPlayerNextLevelXp                 = FieldMask{Size: 1, Offset: 0x27B}
	FieldMaskPlayerSkillInfo                   = FieldMask{Size: 3, Offset: 0x27C}
	FieldMaskPlayerCharacterPoints1            = FieldMask{Size: 1, Offset: 0x3FC}
	FieldMaskPlayerCharacterPoints2            = FieldMask{Size: 1, Offset: 0x3FD}
	FieldMaskPlayerTrackCreatures              = FieldMask{Size: 1, Offset: 0x3FE}
	FieldMaskPlayerTrackResources              = FieldMask{Size: 1, Offset: 0x3FF}
	FieldMaskPlayerBlockPercentage             = FieldMask{Size: 1, Offset: 0x400}
	FieldMaskPlayerDodgePercentage             = FieldMask{Size: 1, Offset: 0x401}
	FieldMaskPlayerParryPercentage             = FieldMask{Size: 1, Offset: 0x402}
	FieldMaskPlayerExpertise                   = FieldMask{Size: 1, Offset: 0x403}
	FieldMaskPlayerOffhandExpertise            = FieldMask{Size: 1, Offset: 0x404}
	FieldMaskPlayerCritPercentage              = FieldMask{Size: 1, Offset: 0x405}
	FieldMaskPlayerRangedCritPercentage        = FieldMask{Size: 1, Offset: 0x406}
	FieldMaskPlayerOffhandCritPercentage       = FieldMask{Size: 1, Offset: 0x407}
	FieldMaskPlayerSpellCritPercentage1        = FieldMask{Size: 7, Offset: 0x408}
	FieldMaskPlayerShieldBlock                 = FieldMask{Size: 1, Offset: 0x40F}
	FieldMaskPlayerShieldBlockCritPercentage   = FieldMask{Size: 1, Offset: 0x410}
	FieldMaskPlayerExploredZones1              = FieldMask{Size: 1, Offset: 0x411}
	FieldMaskPlayerRestStateExperience         = FieldMask{Size: 1, Offset: 0x491}
	FieldMaskPlayerCoinage                     = FieldMask{Size: 1, Offset: 0x492}
	FieldMaskPlayerModDamageDonePos            = FieldMask{Size: 7, Offset: 0x493}
	FieldMaskPlayerModDamageDoneNeg            = FieldMask{Size: 7, Offset: 0x49A}
	FieldMaskPlayerModDamageDonePct            = FieldMask{Size: 7, Offset: 0x4A1}
	FieldMaskPlayerModHealingDonePos           = FieldMask{Size: 1, Offset: 0x4A8}
	FieldMaskPlayerModHealingPct               = FieldMask{Size: 1, Offset: 0x4A9}
	FieldMaskPlayerModHealingDonePct           = FieldMask{Size: 1, Offset: 0x4AA}
	FieldMaskPlayerModTargetResistance         = FieldMask{Size: 1, Offset: 0x4AB}
	FieldMaskPlayerModTargetPhysicalResistance = FieldMask{Size: 1, Offset: 0x4AC}
	FieldMaskPlayerFeatures                    = FieldMask{Size: 1, Offset: 0x4AD}
	FieldMaskPlayerAmmoId                      = FieldMask{Size: 1, Offset: 0x4AE}
	FieldMaskPlayerSelfResSpell                = FieldMask{Size: 1, Offset: 0x4AF}
	FieldMaskPlayerPvpMedals                   = FieldMask{Size: 1, Offset: 0x4B0}
	FieldMaskPlayerBuybackPrice1               = FieldMask{Size: 1, Offset: 0x4B1}
	FieldMaskPlayerBuybackTimestamp1           = FieldMask{Size: 1, Offset: 0x4BD}
	FieldMaskPlayerKills                       = FieldMask{Size: 1, Offset: 0x4C9}
	FieldMaskPlayerTodayContribution           = FieldMask{Size: 1, Offset: 0x4CA}
	FieldMaskPlayerYesterdayContribution       = FieldMask{Size: 1, Offset: 0x4CB}
	FieldMaskPlayerLifetimeHonorableKills      = FieldMask{Size: 1, Offset: 0x4CC}
	FieldMaskPlayerBytes4                      = FieldMask{Size: 1, Offset: 0x4CD}
	FieldMaskPlayerWatchedFactionIndex         = FieldMask{Size: 1, Offset: 0x4CE}
	FieldMaskPlayerCombatRating1               = FieldMask{Size: 2, Offset: 0x4CF}
	FieldMaskPlayerArenaTeamInfo11             = FieldMask{Size: 2, Offset: 0x4E8}
	FieldMaskPlayerHonorCurrency               = FieldMask{Size: 1, Offset: 0x4FD}
	FieldMaskPlayerArenaCurrency               = FieldMask{Size: 1, Offset: 0x4FE}
	FieldMaskPlayerMaxLevel                    = FieldMask{Size: 1, Offset: 0x4FF}
	FieldMaskPlayerDailyQuests1                = FieldMask{Size: 2, Offset: 0x500}
	FieldMaskPlayerRuneRegen1                  = FieldMask{Size: 4, Offset: 0x519}
	FieldMaskPlayerNoReagentCost1              = FieldMask{Size: 3, Offset: 0x51D}
	FieldMaskPlayerGlyphSlots1                 = FieldMask{Size: 6, Offset: 0x520}
	FieldMaskPlayerGlyphs1                     = FieldMask{Size: 6, Offset: 0x526}
	FieldMaskPlayerGlyphsEnabled               = FieldMask{Size: 1, Offset: 0x52C}
	FieldMaskPlayerPetSpellPower               = FieldMask{Size: 1, Offset: 0x52D}

	FieldMaskGameObjectDisplayid      = FieldMask{Size: 1, Offset: 0x8}
	FieldMaskGameObjectFlags          = FieldMask{Size: 1, Offset: 0x9}
	FieldMaskGameObjectParentrotation = FieldMask{Size: 4, Offset: 0xA}
	FieldMaskGameObjectDynamic        = FieldMask{Size: 1, Offset: 0xE}
	FieldMaskGameObjectFaction        = FieldMask{Size: 1, Offset: 0xF}
	FieldMaskGameObjectLevel          = FieldMask{Size: 1, Offset: 0x10}
	FieldMaskGameObjectBytes1         = FieldMask{Size: 1, Offset: 0x11}

	FieldMaskDynamicObjectCaster   = FieldMask{Size: 2, Offset: 0x6}
	FieldMaskDynamicObjectBytes    = FieldMask{Size: 1, Offset: 0x8}
	FieldMaskDynamicObjectSpellid  = FieldMask{Size: 1, Offset: 0x9}
	FieldMaskDynamicObjectRadius   = FieldMask{Size: 1, Offset: 0xA}
	FieldMaskDynamicObjectCasttime = FieldMask{Size: 1, Offset: 0xB}

	FieldMaskCorpseOwner        = FieldMask{Size: 2, Offset: 0x6}
	FieldMaskCorpseParty        = FieldMask{Size: 2, Offset: 0x8}
	FieldMaskCorpseDisplayId    = FieldMask{Size: 1, Offset: 0xA}
	FieldMaskCorpseItem         = FieldMask{Size: 1, Offset: 0xB}
	FieldMaskCorpseBytes1       = FieldMask{Size: 1, Offset: 0x1E}
	FieldMaskCorpseBytes2       = FieldMask{Size: 1, Offset: 0x1F}
	FieldMaskCorpseGuild        = FieldMask{Size: 1, Offset: 0x20}
	FieldMaskCorpseFlags        = FieldMask{Size: 1, Offset: 0x21}
	FieldMaskCorpseDynamicFlags = FieldMask{Size: 1, Offset: 0x22}
)

// ValueMask stores the mask that communicates what fields are set in a ValueBuffer.
type ValueMask struct {
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
func (m *ValueMask) Bit(bit uint32) bool {
	if !m.anyBits || bit > m.largestBit {
		return false
	}

	index := bit / 32
	bitPos := bit % 32

	return m.mask[index]&(1<<bitPos) > 0
}

// Bytes returns a little-endian byte array of the value mask. The first byte is the size of the mask.
func (m *ValueMask) Bytes() []byte {
	buf := bytes.Buffer{}
	size := m.Len()
	buf.WriteByte(byte(size))
	binary.Write(&buf, binary.LittleEndian, m.mask[:size])

	return buf.Bytes()
}

// FieldMask returns whether all the bits for the provided mask have been set.
func (m *ValueMask) FieldMask(fieldMask FieldMask) bool {
	for i := uint32(0); i < fieldMask.Size; i++ {
		if !m.Bit(fieldMask.Offset + i) {
			return false
		}
	}

	return true
}

// Len returns the number of uint32s used to represent the mask.
func (m *ValueMask) Len() int {
	size := m.largestBit / 32

	if m.anyBits {
		size++
	}

	return int(size)
}

// SetFieldMask sets all the bits necessary for the provided field mask.
func (m *ValueMask) SetFieldMask(fieldMask FieldMask) {
	for i := uint32(0); i < fieldMask.Size; i++ {
		m.SetBit(fieldMask.Offset + i)
	}
}

// SetBit sets the nth bit in the update mask. The bit is zero-indexed with the first bit being zero.
func (m *ValueMask) SetBit(bit uint32) {
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
func (m *ValueMask) resize(n uint32) {
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
