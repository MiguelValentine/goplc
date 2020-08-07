package encapsulation

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/MiguelValentine/goplc/enip/etype"
	"github.com/MiguelValentine/goplc/enip/lib"
)

type Command etype.XUINT

const (
	CommandNOP               Command = 0x0000
	CommandListServices      Command = 0x0004
	CommandListIdentity      Command = 0x0063
	CommandListInterfaces    Command = 0x0064
	CommandRegisterSession   Command = 0x0065
	CommandUnRegisterSession Command = 0x0066
	CommandSendRRData        Command = 0x006F
	CommandSendUnitData      Command = 0x0070
	CommandIndicateStatus    Command = 0x0072
	CommandCancel            Command = 0x0073
)

var StatusMap map[etype.XUDINT]string
var CommandMap map[etype.XUINT]Command

func init() {
	StatusMap = make(map[etype.XUDINT]string)
	StatusMap[0x0000] = "Success"
	StatusMap[0x0001] = "The sender issued an invalid or unsupported encapsulation command."
	StatusMap[0x0002] = "Insufficient memory resources in the receiver to handle the command. This is not an application error. Instead, it only results if the encapsulation layer cannot obtain memory resources that it needs."
	StatusMap[0x0003] = "Poorly formed or incorrect data in the data portion of the encapsulation message."
	StatusMap[0x0064] = "An originator used an invalid session handle when sending an encapsulation message to the target."
	StatusMap[0x0065] = "The target received a message of invalid length"
	StatusMap[0x0069] = "Unsupported encapsulation protocol revision."

	CommandMap = make(map[etype.XUINT]Command)
	CommandMap[0x0000] = CommandNOP
	CommandMap[0x0004] = CommandListServices
	CommandMap[0x0063] = CommandListIdentity
	CommandMap[0x0064] = CommandListInterfaces
	CommandMap[0x0065] = CommandRegisterSession
	CommandMap[0x0066] = CommandUnRegisterSession
	CommandMap[0x006F] = CommandSendRRData
	CommandMap[0x0070] = CommandSendUnitData
	CommandMap[0x0072] = CommandIndicateStatus
	CommandMap[0x0073] = CommandCancel
}

type EncapsulationHeader struct {
	Command       Command
	Length        etype.XUINT
	SessionHandle etype.XUDINT
	Status        etype.XUDINT
	SenderContext uint64
	Options       etype.XUDINT
}

type Request struct{}
type Response struct{}

type Encapsulation struct {
	EncapsulationHeader
	Data []byte
}

func (e *Encapsulation) Buffer() []byte {
	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, e.EncapsulationHeader)
	lib.WriteByte(buffer, e.Data)
	return buffer.Bytes()
}

func Parse(data []byte) ([]*Encapsulation, uint64, error) {
	if len(data) < 24 {
		return nil, 0, nil
	}

	dataReader := bytes.NewReader(data)
	var _encapsulations []*Encapsulation
	count := uint64(dataReader.Len())

	for dataReader.Len() > 23 {
		_encapsulation := &Encapsulation{}

		lib.ReadByte(dataReader, &_encapsulation.EncapsulationHeader)

		if _encapsulation.Status != 0 {
			return nil, 0, errors.New(StatusMap[_encapsulation.Status])
		}

		_, ok := CommandMap[etype.XUINT(_encapsulation.Command)]
		if !ok {
			_err := errors.New(fmt.Sprintf("%s,%#x\n", "Encapsulation parser got unknow command", _encapsulation.Command))
			panic(_err)
			return nil, 0, _err
		}

		if int(_encapsulation.Length) > dataReader.Len() {
			count += 24
			break
		} else {
			if _encapsulation.Length > 0 {
				_data := make([]byte, _encapsulation.Length)
				_, err := dataReader.Read(data)
				if err != nil {
					panic(err)
					return nil, 0, err
				}

				_encapsulation.Data = _data
			}
			_encapsulations = append(_encapsulations, _encapsulation)
		}
	}

	count = count - uint64(dataReader.Len())
	return _encapsulations, count, nil
}

func Package(command Command, context uint64, options etype.XUDINT) {

}
