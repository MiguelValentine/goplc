package encapsulation

import (
	"bytes"
	"github.com/MiguelValentine/goplc/enip/etype"
	"github.com/MiguelValentine/goplc/enip/lib"
)

type RegisterSessionData struct {
	ProtocolVersion etype.XUINT
	OptionsFlags    etype.XUINT
}

func (r *Request) RegisterSession(context uint64) []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandRegisterSession
	pkg.Length = 4
	pkg.SenderContext = context

	data := &RegisterSessionData{}
	data.ProtocolVersion = 1
	data.OptionsFlags = 0

	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, data)

	pkg.Data = buffer.Bytes()

	return pkg.Buffer()
}
