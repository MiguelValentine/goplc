package enip

import (
	"bytes"
)

type Request struct{}

func (r *Request) NOP() *encapsulation {
	_encapsulation := &encapsulation{}
	_encapsulation.command = CommandNOP

	return _encapsulation
}

func (r *Request) ListIdentity() *encapsulation {
	_encapsulation := &encapsulation{}
	_encapsulation.command = CommandListIdentity

	return _encapsulation
}

type RegisterSessionData struct {
	ProtocolVersion XUINT
	OptionsFlags    XUINT
}

func (x *RegisterSessionData) Buffer() *bytes.Buffer {
	var _data bytes.Buffer
	writeBinary(&_data, x.ProtocolVersion)
	writeBinary(&_data, x.OptionsFlags)
	return &_data
}

func (r *Request) RegisterSession(context uint64) *encapsulation {
	_encapsulation := &encapsulation{}
	_encapsulation.command = CommandRegisterSession
	_encapsulation.senderContext = context

	_data := RegisterSessionData{
		ProtocolVersion: 1,
		OptionsFlags:    0,
	}

	_dataBuffer := _data.Buffer()

	_encapsulation.length = XUINT(_dataBuffer.Len())
	_encapsulation.data = _dataBuffer
	return _encapsulation
}

func (r *Request) UnRegisterSession(session XUDINT, context uint64) *encapsulation {
	_encapsulation := &encapsulation{}
	_encapsulation.command = CommandUnRegisterSession

	_encapsulation.sessionHandle = session
	_encapsulation.senderContext = context

	return _encapsulation
}

func (r *Request) ListServices(context uint64) *encapsulation {
	_encapsulation := &encapsulation{}
	_encapsulation.command = CommandListServices

	_encapsulation.senderContext = context

	return _encapsulation
}

type SendRRDataData struct {
	InterfaceHandle    XUDINT
	Timeout            XUINT
	EncapsulatedPacket *bytes.Buffer
}

func (x *SendRRDataData) Buffer() *bytes.Buffer {
	var _data bytes.Buffer
	writeBinary(&_data, x.InterfaceHandle)
	writeBinary(&_data, x.Timeout)

	if x.EncapsulatedPacket != nil {
		writeBinary(&_data, x.EncapsulatedPacket.Bytes())
	}

	return &_data
}

func (r *Request) SendRRData(session XUDINT, context uint64, timeout XUINT, cpf *bytes.Buffer) *encapsulation {
	_encapsulation := &encapsulation{}
	_encapsulation.command = CommandSendRRData

	_encapsulation.sessionHandle = session
	_encapsulation.senderContext = context

	_data := SendRRDataData{
		InterfaceHandle:    0,
		Timeout:            timeout,
		EncapsulatedPacket: cpf,
	}

	_dataBuffer := _data.Buffer()
	_encapsulation.length = XUINT(_dataBuffer.Len())
	_encapsulation.data = _dataBuffer
	return _encapsulation
}

type SendUnitData struct {
	InterfaceHandle    XUDINT
	Timeout            XUINT
	EncapsulatedPacket *bytes.Buffer
}

func (x *SendUnitData) Buffer() *bytes.Buffer {
	var _data bytes.Buffer
	writeBinary(&_data, x.InterfaceHandle)
	writeBinary(&_data, x.Timeout)

	if x.EncapsulatedPacket != nil {
		writeBinary(&_data, x.EncapsulatedPacket.Bytes())
	}

	return &_data
}

func (r *Request) SendUnitData(session XUDINT, context uint64, timeout XUINT, cpf *bytes.Buffer) *encapsulation {
	_encapsulation := &encapsulation{}
	_encapsulation.command = CommandSendUnitData

	_encapsulation.sessionHandle = session
	_encapsulation.senderContext = context

	_data := SendRRDataData{
		InterfaceHandle:    0,
		Timeout:            timeout,
		EncapsulatedPacket: cpf,
	}

	_dataBuffer := _data.Buffer()
	_encapsulation.length = XUINT(_dataBuffer.Len())
	_encapsulation.data = _dataBuffer
	return _encapsulation
}
