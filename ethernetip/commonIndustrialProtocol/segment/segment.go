package segment

import (
	"bytes"
	_type "github.com/MiguelValentine/goplc/ethernetip/type"
	"github.com/MiguelValentine/goplc/lib"
)

type SegmentType _type.USINT

const (
	SegmentTypePort      SegmentType = 0 << 5
	SegmentTypeLogical   SegmentType = 1 << 5
	SegmentTypeNetwork   SegmentType = 2 << 5
	SegmentTypeSymbolic  SegmentType = 3 << 5
	SegmentTypeData      SegmentType = 4 << 5
	SegmentTypeDataType1 SegmentType = 5 << 5
	SegmentTypeDataType2 SegmentType = 6 << 5
)

func Paths(arg ...[]byte) []byte {
	buffer := new(bytes.Buffer)
	for i := 0; i < len(arg); i++ {
		lib.WriteByte(buffer, arg[i])
	}
	return buffer.Bytes()
}
