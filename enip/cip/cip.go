package cip

import (
	"bytes"
	"github.com/MiguelValentine/goplc/enip/etype"
	"github.com/MiguelValentine/goplc/enip/lib"
)

type Type etype.XUINT

const (
	TypeNull                      Type = 0x0000
	TypeListIdentityResponse      Type = 0x000c
	TypeConnectionBased           Type = 0xa1
	TypeConnectionTransportPacket Type = 0x00b1
	TypeUnconnectedMessage        Type = 0x00b2
	TypeListServicesResponse      Type = 0x0100
	TypeSockInfoO2T               Type = 0x8000
	TypeSockInfoT2O               Type = 0x8001
	TypeSequencedAddress          Type = 0x8002
)

type Header struct {
	interfaceHandle etype.XUDINT
	timeout         etype.XUINT
}

type CPFItem struct {
	typeID Type
	length etype.XUINT
	data   []byte
}

func (i *CPFItem) Buffer() []byte {
	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, i.typeID)
	lib.WriteByte(buffer, i.length)
	lib.WriteByte(buffer, i.data)
	return buffer.Bytes()
}

func NewCPFItem(t Type, data []byte) *CPFItem {
	_cpfItem := &CPFItem{}
	_cpfItem.typeID = t
	_cpfItem.length = etype.XUINT(len(data))
	_cpfItem.data = data

	return _cpfItem
}

type CPF struct {
	itemCount etype.XUINT
	items     []*CPFItem
	optional  []byte
}

func (c *CPF) Buffer() []byte {
	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, c.itemCount)

	for _, item := range c.items {
		lib.WriteByte(buffer, item.Buffer())
	}

	lib.WriteByte(buffer, c.optional)
	return buffer.Bytes()
}

func NewCPF(items []*CPFItem, opt []byte) *CPF {
	_cpf := &CPF{}
	_cpf.itemCount = etype.XUINT(len(items))
	_cpf.items = items
	_cpf.optional = opt

	return _cpf
}

func Build(timeout etype.XUINT, cpf *CPF) []byte {
	_header := &Header{}
	_header.interfaceHandle = 0
	_header.timeout = timeout

	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, _header)
	lib.WriteByte(buffer, cpf.Buffer())

	return buffer.Bytes()
}
