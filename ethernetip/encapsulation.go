package ethernetip

import (
	"bytes"
	_type "github.com/MiguelValentine/goplc/ethernetip/type"
	"github.com/MiguelValentine/goplc/lib"
)

type Command _type.UINT

const (
	// TCP
	CommandNOP Command = 0x0000

	// TCP && UDP
	CommandListServices   Command = 0x0004
	CommandListIdentity   Command = 0x0063
	CommandListInterfaces Command = 0x0064

	// TCP
	CommandRegisterSession   Command = 0x0065
	CommandUnRegisterSession Command = 0x0066
	CommandSendRRData        Command = 0x006F
	CommandSendUnitData      Command = 0x0070
	CommandIndicateStatus    Command = 0x0072
	CommandCancel            Command = 0x0073
)

type Status _type.UDINT

var StatusMap map[Status]string

const (
	StatusSuccess Status = 0x0000

	StatusUnsupportedCommand Status = 0x0001
	StatusOutOfMemory        Status = 0x0002
	StatusIncorrectData      Status = 0x0003
	StatusInvalidSession     Status = 0x0064
	StatusInvalidLength      Status = 0x0065
	StatusUnsupportedVersion Status = 0x0069
)

func init() {
	StatusMap = make(map[Status]string)
	StatusMap[StatusSuccess] = "Success."
	StatusMap[StatusUnsupportedCommand] = "The sender issued an invalid or unsupported encapsulation command."
	StatusMap[StatusOutOfMemory] = "Insufficient memory resources in the receiver to handle the command. This is not an application error. Instead, it only results if the encapsulation layer cannot obtain memory resources that it needs."
	StatusMap[StatusIncorrectData] = "Poorly formed or incorrect data in the data portion of the encapsulation message."
	StatusMap[StatusInvalidSession] = "An originator used an invalid session handle when sending an encapsulation message to the target."
	StatusMap[StatusInvalidLength] = "The target received a message of invalid length."
	StatusMap[StatusUnsupportedVersion] = "Unsupported encapsulation protocol revision."
}

type EncapsulationHeader struct {
	Command       Command
	Length        _type.UINT
	SessionHandle _type.UDINT
	Status        Status
	SenderContext _type.ULINT
	Options       _type.UDINT
}

type Encapsulation struct {
	EncapsulationHeader
	Data []byte
}

func (e *Encapsulation) Buffer() []byte {
	e.Length = _type.UINT(len(e.Data))

	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, e.EncapsulationHeader)
	lib.WriteByte(buffer, e.Data)

	return buffer.Bytes()
}

func Slice(data []byte) (uint64, []*Encapsulation) {
	var result []*Encapsulation

	dataReader := bytes.NewReader(data)
	count := dataReader.Len()

	for dataReader.Len() > 23 {
		_encapsulation := &Encapsulation{}
		lib.ReadByte(dataReader, &_encapsulation.EncapsulationHeader)

		if int(_encapsulation.Length) > dataReader.Len() {
			count += 24
			break
		} else {
			if _encapsulation.Length > 0 {
				_encapsulation.Data = make([]byte, _encapsulation.Length)
				_, e := dataReader.Read(_encapsulation.Data)
				if e != nil {
					panic(e)
				}

				result = append(result, _encapsulation)
			}
		}
	}

	count = count - dataReader.Len()
	return uint64(count), result
}
