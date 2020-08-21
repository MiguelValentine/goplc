package ethernetip

import (
	"bytes"
	"encoding/binary"
	_type "github.com/MiguelValentine/goplc/ethernetip/type"
)

type SocketAddr struct {
	SinFamily _type.UINT
	SinPort   _type.UINT
	SinAddr   _type.UDINT
	SinZero   _type.ULINT
}

func (s *SocketAddr) Buffer() []byte {
	buffer := new(bytes.Buffer)
	_ = binary.Write(buffer, binary.BigEndian, s.SinFamily)
	_ = binary.Write(buffer, binary.BigEndian, s.SinFamily)
	_ = binary.Write(buffer, binary.BigEndian, s.SinFamily)
	_ = binary.Write(buffer, binary.BigEndian, s.SinFamily)

	return buffer.Bytes()
}
