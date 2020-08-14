package segment

type DataTypes uint8
type DataElementType uint8

const (
	DataTypeSimple DataTypes = 0x80
	DataTypeANSI   DataTypes = 0x91
)

const (
	Uint8ElementID  DataElementType = 0x28
	Uint16ElementID DataElementType = 0x29
	Uint32ElementID DataElementType = 0x2A
)
