package enip

import (
	"bytes"
)

type commonPacketFormat struct {
	itemCount   XUINT
	addressItem []*itemFormat
	dataItem    []*itemFormat
	optional    *bytes.Buffer
}

func (c *commonPacketFormat) Buffer() *bytes.Buffer {
	_result := toBinary(c.itemCount)

	for _, addressItem := range c.addressItem {
		writeBinary(_result, addressItem.Buffer().Bytes())
	}

	for _, dataItem := range c.dataItem {
		writeBinary(_result, dataItem.Buffer().Bytes())
	}

	if c.optional != nil {
		writeBinary(_result, c.optional.Bytes())
	}

	return _result
}

func unconnectedMessages() *commonPacketFormat {
	_result := &commonPacketFormat{}
	_result.itemCount = XUINT(2)
	_addressItem := nullAddressItem()
	_result.addressItem = append(_result.addressItem, _addressItem)

	return _result
}
