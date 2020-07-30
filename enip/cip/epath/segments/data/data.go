package data

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
)

type dataType uint8

const (
	Simple    dataType = 0x80
	ANSI_EXTD dataType = 0x91
)

type elementType uint8

const (
	UINT8  elementType = 0x28
	UINT16 elementType = 0x29
	UINT32 elementType = 0x2a
)

func elementBuild(data uint32) *bytes.Buffer {
	var _type elementType
	var dataBuf bytes.Buffer

	if data < 256 {
		_type = UINT8
		_ = binary.Write(&dataBuf, binary.LittleEndian, uint8(data))
	} else if data < 65536 {
		_type = UINT16
		_ = binary.Write(&dataBuf, binary.LittleEndian, uint8(0))
		_ = binary.Write(&dataBuf, binary.LittleEndian, uint16(data))
	} else {
		_type = UINT32
		_ = binary.Write(&dataBuf, binary.LittleEndian, uint8(0))
		_ = binary.Write(&dataBuf, binary.LittleEndian, data)
	}

	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, _type)
	_ = binary.Write(&buf, binary.LittleEndian, dataBuf.Bytes())
	return &buf
}

func symbolicBuild(data *bytes.Buffer, ANSI bool) *bytes.Buffer {
	var buf bytes.Buffer

	if ANSI {
		_ = binary.Write(&buf, binary.LittleEndian, ANSI_EXTD)
		_ = binary.Write(&buf, binary.LittleEndian, uint8(data.Len()))
	} else {
		_ = binary.Write(&buf, binary.LittleEndian, Simple)
		_ = binary.Write(&buf, binary.LittleEndian, uint8(math.Ceil(float64(data.Len())/2)))
	}

	_ = binary.Write(&buf, binary.BigEndian, data.Bytes())

	if buf.Len()%2 == 1 {
		_ = binary.Write(&buf, binary.LittleEndian, uint8(0))
	}

	return &buf
}

func Build(data interface{}, ANSI bool) *bytes.Buffer {
	if reflect.TypeOf(data).String() == "uint32" {
		return elementBuild(data.(uint32))
	} else if reflect.TypeOf(data).String() == "*bytes.Buffer" {
		return symbolicBuild(data.(*bytes.Buffer), ANSI)
	} else {
		return nil
	}
}
