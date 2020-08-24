package tagGroup

import (
	"bytes"
	"github.com/MiguelValentine/goplc/ethernetip/commonIndustrialProtocol"
	"github.com/MiguelValentine/goplc/ethernetip/commonIndustrialProtocol/segment"
	"github.com/MiguelValentine/goplc/ethernetip/commonIndustrialProtocol/segment/epath"
	"github.com/MiguelValentine/goplc/lib"
	"github.com/MiguelValentine/goplc/tag"
)

const ServiceMultipleServicePacket = commonIndustrialProtocol.Service(0x0a)

type TagGroup struct {
	tags []*tag.Tag
}

func (tg *TagGroup) Add(tag *tag.Tag) {
	tg.tags = append(tg.tags, tag)
}

func (tg *TagGroup) GenerateReadMessageRequest() *commonIndustrialProtocol.MessageRouterRequest {
	mr := &commonIndustrialProtocol.MessageRouterRequest{}
	mr.Service = ServiceMultipleServicePacket

	mr.RequestPath = segment.Paths(
		epath.LogicalBuild(epath.LogicalTypeClassID, 2, true),
		epath.LogicalBuild(epath.LogicalTypeInstanceID, 1, true),
	)

	data := new(bytes.Buffer)
	msg := new(bytes.Buffer)

	lib.WriteByte(data, uint16(len(tg.tags)))
	offset := uint16(2 + 2*len(tg.tags))

	for i := 0; i < len(tg.tags); i++ {
		lib.WriteByte(data, offset)
		_tag := tg.tags[i]
		req := _tag.GenerateReadMessageRequest().Buffer()
		offset = offset + uint16(len(req))
		lib.WriteByte(msg, req)
	}

	lib.WriteByte(data, msg.Bytes())
	mr.RequestData = data.Bytes()

	return mr
}

func (tg *TagGroup) ReadTagParser(mr *commonIndustrialProtocol.MessageRouterResponse) {
	length := uint16(0)
	dataReader := bytes.NewReader(mr.ResponseData)
	lib.ReadByte(dataReader, &length)
	offsets := make([]uint16, length)

	for i := uint16(0); i < length; i++ {
		lib.ReadByte(dataReader, &offsets[i])
	}

	for i := uint16(0); i < length; i++ {
		_length := uint16(0)
		if i == length-1 {
			_length = uint16(dataReader.Len())
		} else {
			_length = offsets[i+1] - offsets[i]
		}
		msg := make([]byte, _length)
		lib.ReadByte(dataReader, &msg)
		tagMr := commonIndustrialProtocol.MRParser(msg)
		tg.tags[i].ReadTagParser(tagMr)
	}
}
