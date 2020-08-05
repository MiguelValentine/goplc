package enip

import "bytes"

type Command XUINT

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

var Status map[XUINT]string

func init() {
	Status = make(map[XUINT]string)
	Status[0x0000] = "Success"
	Status[0x0001] = "The sender issued an invalid or unsupported encapsulation command."
	Status[0x0002] = "Insufficient memory resources in the receiver to handle the command. This is not an application error. Instead, it only results if the encapsulation layer cannot obtain memory resources that it needs."
	Status[0x0003] = "Poorly formed or incorrect data in the data portion of the encapsulation message."
	Status[0x0064] = "An originator used an invalid session handle when sending an encapsulation message to the target."
	Status[0x0065] = "The target received a message of invalid length"
	Status[0x0069] = "Unsupported encapsulation protocol revision."
}

type encapsulationHeader struct {
	command       Command
	length        XUINT
	sessionHandle XUDINT
	status        XUDINT
	senderContext uint64
	options       XUDINT
}

type encapsulation struct {
	encapsulationHeader
	data *bytes.Buffer
}
