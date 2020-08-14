package tag

import (
	"bytes"
	"github.com/MiguelValentine/goplc/enip/cip/dataTypes"
	"github.com/MiguelValentine/goplc/enip/cip/epath/segment"
	"github.com/MiguelValentine/goplc/enip/cip/messageRouter"
	"github.com/MiguelValentine/goplc/enip/lib"
)

type Tag struct {
	name      []byte
	readCount uint16
}

func (t *Tag) GenerateReadMessageRequest() []byte {
	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, segment.DataTypeANSI)
	lib.WriteByte(buffer, uint8(len(t.name)))
	lib.WriteByte(buffer, t.name)

	data := new(bytes.Buffer)
	lib.WriteByte(data, t.readCount)

	return messageRouter.Build(messageRouter.ServiceReadTag, [][]byte{buffer.Bytes()}, data.Bytes())
}

func NewTag(name []byte, dataType dataTypes.DataType) *Tag {
	_tag := &Tag{}
	_tag.name = name
	_tag.readCount = 1

	return _tag
}
