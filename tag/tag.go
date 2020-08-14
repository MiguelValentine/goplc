package tag

import (
	"bytes"
	"github.com/MiguelValentine/goplc/enip/cip/dataTypes"
	"github.com/MiguelValentine/goplc/enip/cip/epath/segment"
	"github.com/MiguelValentine/goplc/enip/cip/messageRouter"
	"github.com/MiguelValentine/goplc/enip/lib"
	"math/rand"
	"reflect"
	"time"
)

type Tag struct {
	Name      []byte
	readCount uint16
	XType     dataTypes.DataType
	Context   uint64
	Value     []byte
}

func (t *Tag) GenerateReadMessageRequest() []byte {
	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, segment.DataTypeANSI)
	lib.WriteByte(buffer, uint8(len(t.Name)))
	lib.WriteByte(buffer, t.Name)

	if buffer.Len()%2 != 0 {
		lib.WriteByte(buffer, uint8(0))
	}

	data := new(bytes.Buffer)
	lib.WriteByte(data, t.readCount)

	return messageRouter.Build(messageRouter.ServiceReadTag, [][]byte{buffer.Bytes()}, data.Bytes())
}

func NewTag(name []byte) *Tag {
	_tag := &Tag{}
	_tag.Name = name
	_tag.readCount = 1
	_tag.XType = dataTypes.AUTO
	rand.Seed(time.Now().Unix())
	_tag.Context = rand.Uint64()

	return _tag
}

func NewTagWithType(name []byte, dataType dataTypes.DataType) *Tag {
	_tag := &Tag{}
	_tag.Name = name
	_tag.readCount = 1
	_tag.XType = dataType
	rand.Seed(time.Now().Unix())
	_tag.Context = rand.Uint64()

	return _tag
}

func (t *Tag) GetValue() interface{} {
	kind := dataTypes.CoverMap[t.XType]
	reader := bytes.NewReader(t.Value)

	switch kind {
	case reflect.Invalid:
		return nil
	case reflect.Bool:
		return 0xFF == t.Value[0]
	case reflect.Int8:
		result := int8(0)
		lib.ReadByte(reader, &result)
		return result
	case reflect.Int16:
		result := int16(0)
		lib.ReadByte(reader, &result)
		return result
	case reflect.Int32:
		result := int32(0)
		lib.ReadByte(reader, &result)
		return result
	case reflect.Int64:
		result := int64(0)
		lib.ReadByte(reader, &result)
		return result
	case reflect.Uint8:
		result := uint8(0)
		lib.ReadByte(reader, &result)
		return result
	case reflect.Uint16:
		result := uint16(0)
		lib.ReadByte(reader, &result)
		return result
	case reflect.Uint32:
		result := uint32(0)
		lib.ReadByte(reader, &result)
		return result
	case reflect.Float32:
		result := float32(0)
		lib.ReadByte(reader, &result)
		return result
	case reflect.Float64:
		result := float64(0)
		lib.ReadByte(reader, &result)
		return result
	case reflect.String:
		return string(t.Value)
	default:
		return nil
	}
}
