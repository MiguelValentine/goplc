package epath

import (
	"bytes"
	"github.com/MiguelValentine/goplc/ethernetip/commonIndustrialProtocol/segment"
	"github.com/MiguelValentine/goplc/lib"
)

func PortBuild(link []byte, portID uint16, padded bool) []byte {
	extendedLinkTag := len(link) > 1
	extendedPortTag := !(portID < 15)

	buffer := new(bytes.Buffer)

	firstByte := uint8(segment.SegmentTypePort)
	if extendedLinkTag {
		firstByte = firstByte | 0x10
	}

	if !extendedPortTag {
		firstByte = firstByte | uint8(portID)
	} else {
		firstByte = firstByte | 0xf
	}

	lib.WriteByte(buffer, firstByte)

	if extendedLinkTag {
		lib.WriteByte(buffer, uint8(len(link)))
	}

	if extendedPortTag {
		lib.WriteByte(buffer, portID)
	}

	lib.WriteByte(buffer, link)

	if padded && buffer.Len()%2 == 1 {
		lib.WriteByte(buffer, uint8(0))
	}

	return buffer.Bytes()
}
