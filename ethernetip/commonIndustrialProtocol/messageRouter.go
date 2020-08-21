package commonIndustrialProtocol

import (
	"bytes"
	_type "github.com/MiguelValentine/goplc/ethernetip/type"
	"github.com/MiguelValentine/goplc/lib"
)

type Service _type.USINT

type MessageRouterRequest struct {
	// service is different for each epath
	Service         Service
	RequestPathSize _type.USINT
	RequestPath     []byte
	RequestData     []byte
}

type MessageRouterResponse struct {
	ReplyService           Service
	Reserved               _type.USINT
	GeneralStatus          _type.USINT
	SizeOfAdditionalStatus _type.USINT
	AdditionalStatus       []byte
	ResponseData           []byte
}

func (req *MessageRouterRequest) Buffer() []byte {
	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, req.Service)

	req.RequestPathSize = _type.USINT(len(req.RequestPath) / 2)
	lib.WriteByte(buffer, req.RequestPathSize)
	lib.WriteByte(buffer, req.RequestPath)
	lib.WriteByte(buffer, req.RequestData)

	return buffer.Bytes()
}

func MRParser(buf []byte) *MessageRouterResponse {
	mr := &MessageRouterResponse{}

	reader := bytes.NewReader(buf)
	lib.ReadByte(reader, &mr.ReplyService)
	lib.ReadByte(reader, &mr.Reserved)
	lib.ReadByte(reader, &mr.GeneralStatus)
	lib.ReadByte(reader, &mr.SizeOfAdditionalStatus)
	mr.AdditionalStatus = make([]byte, mr.SizeOfAdditionalStatus)
	lib.ReadByte(reader, &mr.AdditionalStatus)
	mr.ResponseData = make([]byte, reader.Len())
	lib.ReadByte(reader, &mr.ResponseData)

	return mr
}
