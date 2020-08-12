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
	TypeID Type
	Length etype.XUINT
	Data   []byte
}

func (i *CPFItem) Buffer() []byte {
	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, i.TypeID)
	lib.WriteByte(buffer, i.Length)
	lib.WriteByte(buffer, i.Data)
	return buffer.Bytes()
}

func NewCPFItem(t Type, data []byte) *CPFItem {
	_cpfItem := &CPFItem{}
	_cpfItem.TypeID = t
	_cpfItem.Length = etype.XUINT(len(data))
	_cpfItem.Data = data

	return _cpfItem
}

type CPF struct {
	ItemCount etype.XUINT
	Items     []*CPFItem
	Optional  []byte
}

func (c *CPF) Buffer() []byte {
	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, c.ItemCount)

	for _, item := range c.Items {
		lib.WriteByte(buffer, item.Buffer())
	}

	lib.WriteByte(buffer, c.Optional)
	return buffer.Bytes()
}

func NewCPF(items []*CPFItem, opt []byte) *CPF {
	_cpf := &CPF{}
	_cpf.ItemCount = etype.XUINT(len(items))
	_cpf.Items = items
	_cpf.Optional = opt

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

func Parser(buf []byte) (*Header, *CPF) {
	reader := bytes.NewReader(buf)

	header := &Header{}
	cpf := &CPF{}
	lib.ReadByte(reader, &header.interfaceHandle)
	lib.ReadByte(reader, &header.timeout)
	lib.ReadByte(reader, &cpf.ItemCount)

	cpf.Items = make([]*CPFItem, cpf.ItemCount)
	for i := 0; i < int(cpf.ItemCount); i++ {
		cpf.Items[i] = &CPFItem{}
		lib.ReadByte(reader, &cpf.Items[i].TypeID)
		lib.ReadByte(reader, &cpf.Items[i].Length)
		cpf.Items[i].Data = make([]byte, cpf.Items[i].Length)
		lib.ReadByte(reader, cpf.Items[i].Data)
	}

	if reader.Len() > 0 {
		cpf.Optional = make([]byte, reader.Len())
		lib.ReadByte(reader, cpf.Optional)
	}

	return header, cpf
}
