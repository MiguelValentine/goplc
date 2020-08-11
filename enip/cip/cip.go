package cip

import "github.com/MiguelValentine/goplc/enip/etype"

type Header struct {
	interfaceHandle etype.XUDINT
	timeout         etype.XUINT
}

type CPFItem struct {
	typeID etype.XUINT
	length etype.XUINT
	data   []byte
}

type CPF struct {
	itemCount etype.XUINT
	items     []CPFItem
	optional  []byte
}

type Body struct {
	Header
	CPF
}
