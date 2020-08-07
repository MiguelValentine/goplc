package cip

import "github.com/MiguelValentine/goplc/enip/etype"

type ItemHeader struct {
	TypeID etype.XUINT
	Length etype.XUINT
}

type Item struct {
	ItemHeader
	Data []byte
}

type CPF struct {
	ItemCount   etype.XUINT
	AddressItem Item
	DataItem    Item
	Optional    []byte
}
