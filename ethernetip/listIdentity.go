package ethernetip

import _type "github.com/MiguelValentine/goplc/ethernetip/type"

type CIPIDData struct {
	EncapsulationProtocolVersion _type.UINT
	SocketAddress                SocketAddr
	VendorID                     _type.UINT
	DeviceType                   _type.UINT
	ProductCode                  _type.UINT
	Revision                     _type.USINT
	Status                       _type.UINT
	SerialNumber                 _type.UDINT
	ProductName                  []byte
	State                        _type.USINT
}

func RequestListIdentity() *Encapsulation {
	encapsulation := &Encapsulation{}
	encapsulation.Command = CommandListIdentity

	return encapsulation
}

func HandleListIdentity(encapsulation *Encapsulation) {

}
