package ethernetip

import (
	"bytes"
	_type "github.com/MiguelValentine/goplc/ethernetip/type"
	"github.com/MiguelValentine/goplc/lib"
)

type SendRRDataData struct {
	InterfaceHandle    _type.UDINT
	Timeout            _type.UINT
	CommonPacketFormat *CommonPacketFormat
}

func RequestSendRRData(session _type.UDINT, context _type.ULINT, timeout _type.UINT, cpf *CommonPacketFormat) *Encapsulation {
	encapsulation := &Encapsulation{}
	encapsulation.Command = CommandSendRRData
	encapsulation.SessionHandle = session
	encapsulation.SenderContext = context

	data := &SendRRDataData{}
	data.InterfaceHandle = 0
	data.Timeout = timeout
	data.CommonPacketFormat = cpf

	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, data.InterfaceHandle)
	lib.WriteByte(buffer, data.Timeout)
	lib.WriteByte(buffer, data.CommonPacketFormat.Buffer())

	encapsulation.Data = buffer.Bytes()
	encapsulation.Length = _type.UINT(len(encapsulation.Data))

	return encapsulation
}

func SendRRDataParser(buf []byte) *CommonPacketFormat {
	data := &SendRRDataData{}

	reader := bytes.NewReader(buf)
	lib.ReadByte(reader, &data.InterfaceHandle)
	lib.ReadByte(reader, &data.Timeout)

	cpfBuf := make([]byte, reader.Len())
	lib.ReadByte(reader, &cpfBuf)

	return CPFParser(cpfBuf)
}
