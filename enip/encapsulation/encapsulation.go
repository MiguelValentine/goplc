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

type CPFItemID uint16

const (
	CPFItemID_Null                     CPFItemID = 0x00
	CPFItemID_ListIdentity             CPFItemID = 0x0c
	CPFItemID_ConnectionBased          CPFItemID = 0xa1
	CPFItemID_ConnectedTransportPacket CPFItemID = 0xb1
	CPFItemID_UCMM                     CPFItemID = 0xb2
	CPFItemID_ListServices             CPFItemID = 0x100
	CPFItemID_SockaddrO2T              CPFItemID = 0x8000
	CPFItemID_SockaddrT2O              CPFItemID = 0x8001
	CPFItemID_SequencedAddrItem        CPFItemID = 0x8002
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
	ItemIDs map[string]CPFItemID
}

type cpfDataItem struct {
	TypeID CPFItemID
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
		var TypeID CPFItemID
		_ = binary.Read(buf, binary.BigEndian, &TypeID)

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

type header struct {
	Command Command
	Length  uint16
	Session uint32
	Status  uint32
	Context uint64
	Options uint32
	Data    *bytes.Buffer
}

func (h *header) build() *bytes.Buffer {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, uint16(h.Command))
	_ = binary.Write(&buf, binary.LittleEndian, uint16(h.Data.Len()))
	_ = binary.Write(&buf, binary.LittleEndian, h.Session)
	_ = binary.Write(&buf, binary.LittleEndian, h.Status)
	_ = binary.Write(&buf, binary.LittleEndian, h.Context)
	_ = binary.Write(&buf, binary.LittleEndian, h.Options)
	_ = binary.Write(&buf, binary.LittleEndian, h.Data.Bytes())

	return &buf
}

func (h *header) parse(buf *bytes.Buffer) {
	_ = binary.Read(buf, binary.LittleEndian, &h.Command)
	_ = binary.Read(buf, binary.LittleEndian, &h.Length)
	_ = binary.Read(buf, binary.LittleEndian, &h.Session)
	_ = binary.Read(buf, binary.LittleEndian, &h.Status)
	_ = binary.Read(buf, binary.LittleEndian, &h.Context)
	_ = binary.Read(buf, binary.LittleEndian, &h.Options)
	_ = binary.Read(buf, binary.LittleEndian, h.Data)
}

func (h *header) registerSession() {
	var buf bytes.Buffer

	// Protocol Version (Required to be 1)
	_ = binary.Write(&buf, binary.LittleEndian, uint16(0x01))

	// Opton Flags (Reserved for Future List)
	_ = binary.Write(&buf, binary.LittleEndian, uint16(0))

	h.Command = RegisterSession
	h.Data = &buf
}

func (h *header) unregisterSession(session uint32) {
	h.Command = UnregisterSession
	h.Session = session
}

func (h *header) sendRRData(timeout uint16, data *bytes.Buffer, session uint32) {
	h.Command = SendRRData
	h.Session = session

	var _timeout bytes.Buffer

	// Interface Handle ID (Shall be 0 for CIP)
	_ = binary.Write(&_timeout, binary.LittleEndian, uint32(0))

	// Timeout (sec)
	_ = binary.Write(&_timeout, binary.LittleEndian, timeout)

	var buf bytes.Buffer

	cpfBuf := NewCPF().build([]*cpfDataItem{
		&cpfDataItem{
			TypeID: CPFItemID_Null,
			data:   &bytes.Buffer{},
		},
		&cpfDataItem{
			TypeID: CPFItemID_UCMM,
			data:   data,
		},
	})

	_ = binary.Write(&buf, binary.LittleEndian, _timeout)
	_ = binary.Write(&buf, binary.LittleEndian, cpfBuf)

	h.Data = &buf
}

func (h *header) sendUnitData(session uint32, connectID uint32, seqNumber uint32, data *bytes.Buffer) {
	h.Session = session
	h.Command = SendUnitData

	var _timeout bytes.Buffer

	// Interface Handle ID (Shall be 0 for CIP)
	_ = binary.Write(&_timeout, binary.LittleEndian, uint32(0))

	// Timeout (sec)
	_ = binary.Write(&_timeout, binary.LittleEndian, uint16(0))

	var _seqaddrBuf bytes.Buffer
	_ = binary.Write(&_seqaddrBuf, binary.LittleEndian, connectID)
	_ = binary.Write(&_seqaddrBuf, binary.LittleEndian, seqNumber)

	cpfBuf := NewCPF().build([]*cpfDataItem{
		&cpfDataItem{
			TypeID: CPFItemID_SequencedAddrItem,
			data:   &_seqaddrBuf,
		},
		&cpfDataItem{
			TypeID: CPFItemID_ConnectedTransportPacket,
			data:   data,
		},
	})

	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, _timeout)
	_ = binary.Write(&buf, binary.LittleEndian, cpfBuf)

	h.Data = &buf
}

func NewCPF() *cpf {
	_cpf := &cpf{}
	_cpf.ItemIDs = make(map[string]CPFItemID)
	_cpf.ItemIDs["Null"] = CPFItemID_Null
	_cpf.ItemIDs["ListIdentity"] = CPFItemID_ListIdentity
	_cpf.ItemIDs["ConnectionBased"] = CPFItemID_ConnectionBased
	_cpf.ItemIDs["ConnectedTransportPacket"] = CPFItemID_ConnectedTransportPacket
	_cpf.ItemIDs["UCMM"] = CPFItemID_UCMM
	_cpf.ItemIDs["ListServices"] = CPFItemID_ListServices
	_cpf.ItemIDs["SockaddrO2T"] = CPFItemID_SockaddrO2T
	_cpf.ItemIDs["SockaddrT2O"] = CPFItemID_SockaddrT2O
	_cpf.ItemIDs["SequencedAddrItem"] = CPFItemID_SequencedAddrItem
	return _cpf
}

func NewHeader() *header {
	_header := &header{}
	return _header
}
