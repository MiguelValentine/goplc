package segment

import (
	"bytes"
	"github.com/MiguelValentine/goplc/enip/lib"
)

type LogicalType uint8

const (
	LogicalTypeClassID     LogicalType = 0 << 2
	LogicalTypeInstanceID  LogicalType = 1 << 2
	LogicalTypeMemberID    LogicalType = 2 << 2
	LogicalTypeConnPoint   LogicalType = 3 << 2
	LogicalTypeAttributeID LogicalType = 4 << 2
	LogicalTypeSpecial     LogicalType = 5 << 2
	LogicalTypeServiceID   LogicalType = 6 << 2
)

func LogicalBuild(_type LogicalType, address uint32, padded bool) []byte {
	format := uint8(0)

	if address <= 255 {
		format = 0
	} else if address > 255 && address <= 65535 {
		format = 1
	} else {
		format = 2
	}

	id := uint8(TypeLogical) | uint8(_type) | format

	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, id)

	if address > 255 && padded {
		lib.WriteByte(buffer, uint8(0))
	}

	if address <= 255 {
		lib.WriteByte(buffer, uint8(address))
	} else if address > 255 && address <= 65535 {
		lib.WriteByte(buffer, uint16(address))
	} else {
		lib.WriteByte(buffer, address)
	}

	return buffer.Bytes()
}
