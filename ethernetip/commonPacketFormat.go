package ethernetip

import (
	"bytes"
	_type "github.com/MiguelValentine/goplc/ethernetip/type"
	"github.com/MiguelValentine/goplc/lib"
)

type CommonPacketFormatType _type.UINT

const (
	TypeNull                      CommonPacketFormatType = 0x0000
	TypeListIdentityResponse      CommonPacketFormatType = 0x000c
	TypeConnectionBased           CommonPacketFormatType = 0x00a1
	TypeConnectionTransportPacket CommonPacketFormatType = 0x00b1
	TypeUnconnectedMessage        CommonPacketFormatType = 0x00b2
	TypeListServicesResponse      CommonPacketFormatType = 0x0100
	TypeSockInfoO2T               CommonPacketFormatType = 0x8000
	TypeSockInfoT2O               CommonPacketFormatType = 0x8001
	TypeSequencedAddress          CommonPacketFormatType = 0x8002
)

type CommonPacketFormat struct {
	ItemCount      _type.UINT
	AddressItem    CommonPacketFormatItem
	DataItem       CommonPacketFormatItem
	AdditionalItem []CommonPacketFormatItem
}

type CommonPacketFormatItem struct {
	TypeID CommonPacketFormatType
	Length _type.UINT
	Data   []byte
}

func (c *CommonPacketFormat) NullAddress() {}

func (c *CommonPacketFormat) ConnectedAddress(connectionID _type.UDINT) {
	c.AddressItem.TypeID = TypeConnectionBased
	c.AddressItem.Length = 4

	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, connectionID)

	c.AddressItem.Data = buffer.Bytes()
}

func (c *CommonPacketFormat) SequencedAddress(connectionID _type.UDINT, sequenceNumber _type.UDINT) {
	c.AddressItem.TypeID = TypeSequencedAddress
	c.AddressItem.Length = 8

	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, connectionID)
	lib.WriteByte(buffer, sequenceNumber)

	c.AddressItem.Data = buffer.Bytes()
}

func (c *CommonPacketFormat) UnconnectedData(data []byte) {
	c.DataItem.TypeID = TypeUnconnectedMessage
	c.DataItem.Length = _type.UINT(len(data))
	c.DataItem.Data = data
}

func (c *CommonPacketFormat) ConnectedData(data []byte) {
	c.DataItem.TypeID = TypeConnectionTransportPacket
	c.DataItem.Length = _type.UINT(len(data))
	c.DataItem.Data = data
}

func (c *CommonPacketFormat) AddItem(item CommonPacketFormatItem) {
	c.AdditionalItem = append(c.AdditionalItem, item)
}

func (c *CommonPacketFormat) Buffer() []byte {
	c.ItemCount = 2 + _type.UINT(len(c.AdditionalItem))

	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, c.ItemCount)

	lib.WriteByte(buffer, c.AddressItem.TypeID)
	lib.WriteByte(buffer, c.AddressItem.Length)
	lib.WriteByte(buffer, c.AddressItem.Data)

	lib.WriteByte(buffer, c.DataItem.TypeID)
	lib.WriteByte(buffer, c.DataItem.Length)
	lib.WriteByte(buffer, c.DataItem.Data)

	for i := 0; i < len(c.AdditionalItem); i++ {
		item := c.AdditionalItem[i]
		lib.WriteByte(buffer, item.TypeID)
		lib.WriteByte(buffer, item.Length)
		lib.WriteByte(buffer, item.Data)
	}

	return buffer.Bytes()
}

func CPFParser(buf []byte) *CommonPacketFormat {
	cpf := &CommonPacketFormat{}

	reader := bytes.NewReader(buf)
	lib.ReadByte(reader, &cpf.ItemCount)

	lib.ReadByte(reader, &cpf.AddressItem.TypeID)
	lib.ReadByte(reader, &cpf.AddressItem.Length)
	cpf.AddressItem.Data = make([]byte, cpf.AddressItem.Length)
	lib.ReadByte(reader, &cpf.AddressItem.Data)

	lib.ReadByte(reader, &cpf.DataItem.TypeID)
	lib.ReadByte(reader, &cpf.DataItem.Length)
	cpf.DataItem.Data = make([]byte, cpf.DataItem.Length)
	lib.ReadByte(reader, &cpf.DataItem.Data)

	return cpf
}
