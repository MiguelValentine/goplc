package tag

import (
	"bytes"
	"github.com/MiguelValentine/goplc/ethernetip/commonIndustrialProtocol"
	"github.com/MiguelValentine/goplc/ethernetip/commonIndustrialProtocol/segment/epath"
	"github.com/MiguelValentine/goplc/lib"
	"log"
)

const ServiceReadTag = commonIndustrialProtocol.Service(0x4c)

type Tag struct {
	name      []byte
	readCount uint16
	xtype     DataType
	value     []byte
	Onchange  func(interface{})
}

func (t *Tag) GenerateReadMessageRequest() *commonIndustrialProtocol.MessageRouterRequest {
	mr := &commonIndustrialProtocol.MessageRouterRequest{}
	mr.Service = ServiceReadTag
	mr.RequestPath = epath.DataBuild(epath.DataTypeANSI, t.name, true)

	data := new(bytes.Buffer)
	lib.WriteByte(data, t.readCount)
	mr.RequestData = data.Bytes()

	return mr
}

func (t *Tag) Parser(mr *commonIndustrialProtocol.MessageRouterResponse) {
	log.Printf("%+v\n", mr)
	dataReader := bytes.NewReader(mr.ResponseData)
	lib.ReadByte(dataReader, &t.xtype)
	t.value = make([]byte, dataReader.Len())
	lib.ReadByte(dataReader, t.value)
	if t.Onchange != nil {
		t.Onchange(t.GetValue())
	}
}

func NewTag(name string) *Tag {
	_tag := &Tag{}
	_tag.name = []byte(name)
	_tag.readCount = 1
	_tag.xtype = NULL
	return _tag
}

func NewTagWithType(name string, tp DataType) *Tag {
	_tag := &Tag{}
	_tag.name = []byte(name)
	_tag.readCount = 1
	_tag.xtype = tp
	return _tag
}

func (t *Tag) GetValue() interface{} {
	reader := bytes.NewReader(t.value)

	switch t.xtype {
	case NULL:
		return nil
	//case BOOL:
	//return 0xFF == t.Value[0]
	case SINT:
		result := int8(0)
		lib.ReadByte(reader, &result)
		return result
	case INT:
		result := int16(0)
		lib.ReadByte(reader, &result)
		return result
	case DINT:
		result := int32(0)
		lib.ReadByte(reader, &result)
		return result
	case LINT:
		result := int64(0)
		lib.ReadByte(reader, &result)
		return result
	default:
		return nil
	}
}
