package messageRouter

import (
	"bytes"
	"github.com/MiguelValentine/goplc/enip/etype"
	"github.com/MiguelValentine/goplc/enip/lib"
)

type Service etype.XUSINT

const (
	ServiceGetAttributeAll              Service = 0x01
	ServiceSetAttributeAll              Service = 0x02
	ServiceGetAttributeList             Service = 0x03
	ServiceSetAttributeList             Service = 0x04
	ServiceReset                        Service = 0x05
	ServiceStart                        Service = 0x06
	ServiceStop                         Service = 0x07
	ServiceCreate                       Service = 0x08
	ServiceDelete                       Service = 0x09
	ServiceMultipleServicePacket        Service = 0x0a
	ServiceApplyAttributes              Service = 0x0d
	ServiceGetAttributeSingle           Service = 0x0e
	ServiceSetAttributeSingle           Service = 0x10
	ServiceFindNext                     Service = 0x11
	ServiceRestore                      Service = 0x15
	ServiceSave                         Service = 0x16
	ServiceGetMember                    Service = 0x18
	ServiceSetMember                    Service = 0x19
	ServiceInsertMember                 Service = 0x1A
	ServiceRemoveMember                 Service = 0x1B
	ServiceGroupSync                    Service = 0x1C
	ServiceReadTag                      Service = 0x4c
	ServiceWriteTag                     Service = 0x4d
	ServiceFragmentedReadModifyWriteTag Service = 0x4e
	ServiceReadTagFragmented            Service = 0x52
	ServiceUnconnectedSend              Service = 0x52
	ServiceWriteTagFragmented           Service = 0x53
	ServiceGetInstanceAttributeList     Service = 0x55
)

type request struct {
	Service         Service
	RequestPathSize etype.XUSINT
	RequestPath     []byte
	RequestData     []byte
}

type response struct {
	Service                Service
	Reserved               uint8
	GeneralStatus          etype.XUSINT
	SizeOfAdditionalStatus etype.XUSINT
	AdditionalStatus       []byte
	ResponseData           []byte
}

func (r *request) Buffer() []byte {
	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, r.Service)
	lib.WriteByte(buffer, r.RequestPathSize)
	lib.WriteByte(buffer, r.RequestPath)
	lib.WriteByte(buffer, r.RequestData)

	return buffer.Bytes()
}

func Build(service Service, paths [][]byte, data []byte) []byte {
	buffer := new(bytes.Buffer)
	for _, path := range paths {
		lib.WriteByte(buffer, path)
	}

	_path := buffer.Bytes()

	_request := &request{}
	_request.Service = service
	_request.RequestPathSize = etype.XUSINT(len(_path) / 2)
	_request.RequestPath = _path
	_request.RequestData = data

	return _request.Buffer()
}

func Parse(buf []byte) *response {
	r := &response{}
	reader := bytes.NewReader(buf)
	lib.ReadByte(reader, &r.Service)
	lib.ReadByte(reader, &r.Reserved)
	lib.ReadByte(reader, &r.GeneralStatus)
	lib.ReadByte(reader, &r.SizeOfAdditionalStatus)
	r.AdditionalStatus = make([]byte, r.SizeOfAdditionalStatus*2)
	lib.ReadByte(reader, r.AdditionalStatus)
	r.ResponseData = make([]byte, reader.Len())
	lib.ReadByte(reader, r.ResponseData)

	return r
}
