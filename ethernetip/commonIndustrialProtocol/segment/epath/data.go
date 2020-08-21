package epath

import (
	"bytes"
	"github.com/MiguelValentine/goplc/ethernetip/commonIndustrialProtocol/segment"
	_type "github.com/MiguelValentine/goplc/ethernetip/type"
	"github.com/MiguelValentine/goplc/lib"
)

type DataTypes _type.USINT

const (
	DataTypeSimple DataTypes = 0x0
	DataTypeANSI   DataTypes = 0x11
)

func DataBuild(tp DataTypes, data []byte, padded bool) []byte {
	buffer := new(bytes.Buffer)

	firstByte := uint8(segment.SegmentTypeData) | uint8(tp)
	lib.WriteByte(buffer, firstByte)
	length := uint8(len(data))
	lib.WriteByte(buffer, length)
	lib.WriteByte(buffer, data)
	if padded && buffer.Len()%2 == 1 {
		lib.WriteByte(buffer, uint8(0))
	}

	return buffer.Bytes()
}
