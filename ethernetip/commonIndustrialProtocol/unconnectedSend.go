package commonIndustrialProtocol

import (
	"bytes"
	_type "github.com/MiguelValentine/goplc/ethernetip/type"
	"github.com/MiguelValentine/goplc/lib"
)

type UnconnectedSend struct {
	// Actual Time Out value = 2^time_tick x Time_out_tick
	TimeTick           _type.USINT
	TimeOutTicks       _type.USINT
	MessageRequestSize _type.UINT
	MessageRequest     *MessageRouterRequest
	Pad                _type.USINT
	RouterPathSize     _type.USINT
	Reserved           _type.USINT
	RouterPath         []byte
}

func (u *UnconnectedSend) Buffer() []byte {
	buffer := new(bytes.Buffer)

	lib.WriteByte(buffer, u.TimeTick)
	lib.WriteByte(buffer, u.TimeOutTicks)
	lib.WriteByte(buffer, _type.UINT(len(u.MessageRequest.Buffer())))
	lib.WriteByte(buffer, u.MessageRequest.Buffer())

	if len(u.MessageRequest.Buffer())%2 == 1 {
		lib.WriteByte(buffer, uint8(0))
	}

	lib.WriteByte(buffer, _type.USINT(len(u.RouterPath)/2))
	lib.WriteByte(buffer, uint8(0))
	lib.WriteByte(buffer, u.RouterPath)

	return buffer.Bytes()
}
