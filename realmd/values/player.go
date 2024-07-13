package values

type Player struct {
	*ObjectData
	*UnitData
	*PlayerData
}

func NewPlayer() *Player {
	return &Player{
		ObjectData: NewObjectData(),
		UnitData:   NewUnitData(),
		PlayerData: NewPlayerData(),
	}
}

func (p *Player) Marshal(onlyDirty bool) ([]byte, *blockMask) {
	var mask blockMask
	var data, d []byte
	var s []structSection

	d, s = p.ObjectData.Marshal(onlyDirty)
	data = append(data, d...)
	mask.Update(s, ObjectDataOffset)

	d, s = p.UnitData.Marshal(onlyDirty)
	data = append(data, d...)
	mask.Update(s, UnitDataOffset)

	d, s = p.PlayerData.Marshal(onlyDirty)
	data = append(data, d...)
	mask.Update(s, PlayerDataOffset)

	return data, &mask
}
