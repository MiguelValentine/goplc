package encapsulation

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Command uint8

const (
	NOP               Command = 0x00
	ListServices      Command = 0x04
	ListIdentity      Command = 0x63
	ListInterfaces    Command = 0x64
	RegisterSession   Command = 0x65
	UnregisterSession Command = 0x66
	SendRRData        Command = 0x6f
	SendUnitData      Command = 0x70
	IndicateStatus    Command = 0x72
	Cancel            Command = 0x73
)

func PaeseStatus(status uint8) string {
	switch status {
	case 0x00:
		return "SUCCESS"
	case 0x01:
		return "FAIL: Sender issued an invalid ecapsulation command."
	case 0x02:
		return "FAIL: Insufficient memory resources to handle command."
	case 0x03:
		return "FAIL: Poorly formed or incorrect data in encapsulation packet."
	case 0x64:
		return "FAIL: Originator used an invalid session handle."
	case 0x65:
		return "FAIL: Target received a message of invalid length."
	case 0x69:
		return "FAIL: Unsupported encapsulation protocol revision."
	default:
		return fmt.Sprintf("FAIL: General failure <%x> occured.", status)
	}
}

type cpf struct {
	ItemIDs map[string]uint16
}

type cpfDataItem struct {
	TypeID uint16
	data   *bytes.Buffer
}

func (c *cpf) isCmd(cmd string) bool {
	_, ok := c.ItemIDs[cmd]
	return ok
}

func (c *cpf) build(dataItems []*cpfDataItem) *bytes.Buffer {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, uint16(len(dataItems)))

	for _, item := range dataItems {
		var buf1 bytes.Buffer
		var buf2 bytes.Buffer

		_ = binary.Write(&buf2, binary.BigEndian, item.data.Bytes())
		_ = binary.Write(&buf1, binary.LittleEndian, item.TypeID)
		_ = binary.Write(&buf1, binary.LittleEndian, uint16(buf2.Len()))

		_ = binary.Write(&buf, binary.BigEndian, buf1.Bytes())

		if buf2.Len() > 0 {
			_ = binary.Write(&buf, binary.BigEndian, buf2.Bytes())
		}
	}

	return &buf
}

func (c *cpf) parse(buf *bytes.Buffer) []*cpfDataItem {
	itemCountB := make([]byte, 2)
	_ = binary.Read(buf, binary.BigEndian, itemCountB)

	itemCount := binary.LittleEndian.Uint16(itemCountB)

	var result []*cpfDataItem

	for i := uint16(0); i < itemCount; i++ {
		TypeIDB := make([]byte, 2)
		_ = binary.Read(buf, binary.BigEndian, TypeIDB)
		TypeID := binary.LittleEndian.Uint16(TypeIDB)

		lengthB := make([]byte, 2)
		_ = binary.Read(buf, binary.BigEndian, lengthB)
		length := binary.LittleEndian.Uint16(lengthB)

		data := make([]byte, length)
		_ = binary.Read(buf, binary.BigEndian, data)

		item := &cpfDataItem{}
		item.TypeID = TypeID
		item.data = bytes.NewBuffer(data)
		result = append(result, item)
	}

	return result
}

var CPF *cpf

type header struct {
}

func init() {
	CPF := &cpf{}
	CPF.ItemIDs = make(map[string]uint16)
	CPF.ItemIDs["Null"] = 0x00
	CPF.ItemIDs["ListIdentity"] = 0x0c
	CPF.ItemIDs["ConnectionBased"] = 0xa1
	CPF.ItemIDs["ConnectedTransportPacket"] = 0xb1
	CPF.ItemIDs["UCMM"] = 0xb2
	CPF.ItemIDs["ListServices"] = 0x100
	CPF.ItemIDs["SockaddrO2T"] = 0x8000
	CPF.ItemIDs["SockaddrT2O"] = 0x8001
	CPF.ItemIDs["SequencedAddrItem"] = 0x8002
}
