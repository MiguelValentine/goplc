package logical

import (
	"bytes"
	"encoding/binary"
	"github.com/MiguelValentine/goplc/enip/cip/epath/segments/segmentsTypes"
)

type LogicalTpye uint8

const (
	ClassID     LogicalTpye = 0 << 2
	InstanceID  LogicalTpye = 1 << 2
	MemberID    LogicalTpye = 2 << 2
	ConnPoint   LogicalTpye = 3 << 2
	AttributeID LogicalTpye = 4 << 2
	Special     LogicalTpye = 5 << 2
	ServiceID   LogicalTpye = 6 << 2
)

// padded Default should be true
func Build(_type LogicalTpye, address uint32, padded bool) *bytes.Buffer {
	var result bytes.Buffer
	var body bytes.Buffer
	format := 0

	if address <= 255 {
		format = 0
		_ = binary.Write(&body, binary.BigEndian, uint8(address))
	} else if address > 255 && address <= 65535 {
		format = 1

		if padded {
			_ = binary.Write(&body, binary.BigEndian, uint8(0))
		}

		_ = binary.Write(&body, binary.LittleEndian, uint16(address))
	} else {
		format = 2

		if padded {
			_ = binary.Write(&body, binary.BigEndian, uint8(0))
		}

		_ = binary.Write(&body, binary.LittleEndian, address)
	}

	segmentByte := uint8(segmentsTypes.LOGICAL) | uint8(_type) | uint8(format)
	_ = binary.Write(&result, binary.BigEndian, segmentByte)
	_ = binary.Write(&result, binary.BigEndian, body.Bytes())

	return &result
}
