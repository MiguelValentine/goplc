package cip

import (
	"bytes"
	"github.com/MiguelValentine/goplc/enip/etype"
	"github.com/MiguelValentine/goplc/enip/lib"
)

type Identity struct {
	ItemTypeCode                 etype.XUINT
	ItemLength                   etype.XUINT
	EncapsulationProtocolVersion etype.XUINT
	SocketAddr                   *SocketAddress
	VendorID                     etype.XUINT
	DeviceType                   etype.XUINT
	ProductCode                  etype.XUINT
	Revision                     etype.XUSINT
	Status                       etype.XUINT
	SerialNumber                 etype.XUDINT
	ProductName                  []byte
	State                        etype.XUSINT
}

func (i *Identity) Buffer() []byte {
	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, i.ItemTypeCode)
	lib.WriteByte(buffer, i.ItemLength)
	lib.WriteByte(buffer, i.EncapsulationProtocolVersion)
	lib.WriteByte(buffer, i.SocketAddr.Buffer())
	lib.WriteByte(buffer, i.VendorID)
	lib.WriteByte(buffer, i.DeviceType)
	lib.WriteByte(buffer, i.ProductCode)
	lib.WriteByte(buffer, i.Revision)
	lib.WriteByte(buffer, i.Status)
	lib.WriteByte(buffer, i.SerialNumber)
	lib.WriteByte(buffer, i.ProductName)
	lib.WriteByte(buffer, i.State)
	return buffer.Bytes()
}
