package segment

import (
	"bytes"
	"github.com/MiguelValentine/goplc/enip/lib"
)

type Port struct {
	id   uint8
	size uint8
	port uint16
	link []byte
}

func PortBuild(port uint16, link []byte) []byte {
	buffer := new(bytes.Buffer)
	_port := &Port{}
	_port.id = uint8(TypePort)
	_port.size = uint8(len(link))
	_port.port = port
	_port.link = link

	extendedLinkTag := _port.size > 1
	extendedPortTag := !(port < 15)

	if extendedLinkTag {
		_port.id = _port.id | 0x10
	}

	if !extendedPortTag {
		_port.id = uint8(TypePort) | uint8(port)
	} else {
		_port.id = uint8(TypePort) | 0xf
	}

	lib.WriteByte(buffer, _port.id)

	if extendedLinkTag {
		lib.WriteByte(buffer, _port.size)
	}

	if extendedPortTag {
		lib.WriteByte(buffer, _port.port)
	}

	lib.WriteByte(buffer, _port.link)

	if extendedLinkTag && _port.size%2 == 1 {
		lib.WriteByte(buffer, uint8(0))
	}

	return buffer.Bytes()
}
