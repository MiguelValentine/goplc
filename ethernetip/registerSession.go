package ethernetip

import (
	"bytes"
	_type "github.com/MiguelValentine/goplc/ethernetip/type"
	"github.com/MiguelValentine/goplc/lib"
)

type RegisterSessionData struct {
	ProtocolVersion _type.UINT
	OptionsFlags    _type.UINT
}

func RequestRegisterSession(context _type.ULINT) *Encapsulation {
	encapsulation := &Encapsulation{}
	encapsulation.Command = CommandRegisterSession
	encapsulation.Length = 4
	encapsulation.SenderContext = context

	data := &RegisterSessionData{}
	data.ProtocolVersion = 1
	data.OptionsFlags = 0

	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, data)

	encapsulation.Data = buffer.Bytes()

	return encapsulation
}
