package segment

import "github.com/MiguelValentine/goplc/enip/etype"

const (
	TypePort      etype.XUSINT = 0 << 5
	TypeLogical   etype.XUSINT = 1 << 5
	TypeNetwork   etype.XUSINT = 2 << 5
	TypeSymbolic  etype.XUSINT = 3 << 5
	TypeData      etype.XUSINT = 4 << 5
	TypeDataType1 etype.XUSINT = 5 << 5
	TypeDataType2 etype.XUSINT = 6 << 5
)
