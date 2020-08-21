package ethernetip

import (
	"bytes"
	_type "github.com/MiguelValentine/goplc/ethernetip/type"
	"github.com/MiguelValentine/goplc/lib"
)

type SendUnitDataData struct {
	InterfaceHandle    _type.UDINT
	Timeout            _type.UINT
	CommonPacketFormat *CommonPacketFormat
}

func RequestSendUnitData(session _type.UDINT, context _type.ULINT, cpf *CommonPacketFormat) *Encapsulation {
	encapsulation := &Encapsulation{}
	encapsulation.Command = CommandSendUnitData
	encapsulation.SessionHandle = session
	encapsulation.SenderContext = context

	data := &SendUnitDataData{}
	data.InterfaceHandle = 0
	data.Timeout = 0
	data.CommonPacketFormat = cpf

	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, data.InterfaceHandle)
	lib.WriteByte(buffer, data.Timeout)
	lib.WriteByte(buffer, data.CommonPacketFormat.Buffer())

	encapsulation.Data = buffer.Bytes()
	encapsulation.Length = _type.UINT(len(encapsulation.Data))

	return encapsulation
}
