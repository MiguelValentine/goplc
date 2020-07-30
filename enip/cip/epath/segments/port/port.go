package port

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/MiguelValentine/goplc/enip/cip/epath/segments/segmentsTypes"
	"reflect"
)

func Build(port uint16, link interface{}) (error, *bytes.Buffer) {
	_type := reflect.TypeOf(link).String()
	if _type != "string" && _type != "uint8" {
		return errors.New("link number must be a uint8 or string"), nil
	}

	var result bytes.Buffer
	var linkBuf bytes.Buffer
	portIdentifierByte := uint8(segmentsTypes.PORT)

	switch _type {
	case "string":
		_ = binary.Write(&linkBuf, binary.BigEndian, []byte(link.(string)))
	case "uint8":
		_ = binary.Write(&linkBuf, binary.BigEndian, link.(uint8))
	}

	if port < 15 {
		portIdentifierByte |= uint8(port)
	} else {
		portIdentifierByte |= 0x0f
	}

	if linkBuf.Len() > 1 {
		portIdentifierByte |= 0x10
	}

	_ = binary.Write(&result, binary.BigEndian, portIdentifierByte)
	if port < 15 {
		if linkBuf.Len() > 1 {
			_ = binary.Write(&result, binary.BigEndian, uint8(linkBuf.Len()))
		}
	} else {
		if linkBuf.Len() > 1 {
			_ = binary.Write(&result, binary.BigEndian, uint8(linkBuf.Len()))
			_ = binary.Write(&result, binary.LittleEndian, port)
		} else {
			_ = binary.Write(&result, binary.LittleEndian, port)
		}
	}

	_ = binary.Write(&result, binary.BigEndian, linkBuf.Bytes())
	if result.Len()%2 == 1 {
		_ = binary.Write(&result, binary.BigEndian, uint8(0))
	}

	return nil, &result
}
