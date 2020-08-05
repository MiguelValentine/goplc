package enip

import (
	"bytes"
	"log"
)

type itemFormat struct {
	typeID XUINT
	length XUINT
	data   *bytes.Buffer
}

func (i *itemFormat) Buffer() *bytes.Buffer {
	_data := toBinary(i.typeID)
	writeBinary(_data, i.length)
	if i.data != nil {
		writeBinary(_data, i.data.Bytes())
	}
	return _data
}

func nullAddressItem() *itemFormat {
	return &itemFormat{}
}

func connectedAddressItem(data XUDINT) *itemFormat {
	_data := toBinary(data)

	return &itemFormat{
		typeID: 0xA1,
		length: 4,
		data:   _data,
	}
}

func sequencedAddressItem(connectionID XUDINT, sequenceNumber XUDINT) *itemFormat {
	_data := toBinary(connectionID)
	writeBinary(_data, sequenceNumber)

	return &itemFormat{
		typeID: 0x8002,
		length: 8,
		data:   _data,
	}
}

func unconnectedDataItem(data *bytes.Buffer) *itemFormat {
	if data == nil {
		log.Panicln("Data can't be nil")
	}

	return &itemFormat{
		typeID: 0xB2,
		length: XUINT(data.Len()),
		data:   data,
	}
}

func connectedDataItem(data *bytes.Buffer) *itemFormat {
	if data == nil {
		log.Panicln("Data can't be nil")
	}

	return &itemFormat{
		typeID: 0xB1,
		length: XUINT(data.Len()),
		data:   data,
	}
}

func sockaddrInfoItem(o2t bool, sin *bytes.Buffer) *itemFormat {
	var _typeID XUINT
	if o2t {
		_typeID = 0x8000
	} else {
		_typeID = 0x8001
	}

	return &itemFormat{
		typeID: _typeID,
		length: 16,
		data:   sin,
	}
}
