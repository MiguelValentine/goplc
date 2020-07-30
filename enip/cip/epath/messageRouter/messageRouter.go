package messageRouter

import (
	"bytes"
	"encoding/binary"
	"math"
)

type Service uint8

var ServicesMap map[uint8]Service

const (
	GET_ATTRIBUTE_ALL       Service = 0x01
	GET_ATTRIBUTE_SINGLE    Service = 0x0e
	RESET                   Service = 0x05
	START                   Service = 0x06
	STOP                    Service = 0x07
	CREATE                  Service = 0x08
	DELETE                  Service = 0x09
	MULTIPLE_SERVICE_PACKET Service = 0x0a
	APPLY_ATTRIBUTES        Service = 0x0d
	SET_ATTRIBUTE_SINGLE    Service = 0x10
	FIND_NEXT               Service = 0x11
	READ_TAG                Service = 0x4c
	WRITE_TAG               Service = 0x4d
	READ_MODIFY_WRITE_TAG   Service = 0x4e
	READ_TAG_FRAGMENTED     Service = 0x52
	WRITE_TAG_FRAGMENTED    Service = 0x53
)

func init() {
	ServicesMap = make(map[uint8]Service)
	ServicesMap[0x01] = GET_ATTRIBUTE_ALL
	ServicesMap[0x0e] = GET_ATTRIBUTE_SINGLE
	ServicesMap[0x05] = RESET
	ServicesMap[0x06] = START
	ServicesMap[0x07] = STOP
	ServicesMap[0x08] = CREATE
	ServicesMap[0x09] = DELETE
	ServicesMap[0x0a] = MULTIPLE_SERVICE_PACKET
	ServicesMap[0x0d] = APPLY_ATTRIBUTES
	ServicesMap[0x10] = SET_ATTRIBUTE_SINGLE
	ServicesMap[0x11] = FIND_NEXT
	ServicesMap[0x4c] = READ_TAG
	ServicesMap[0x4d] = WRITE_TAG
	ServicesMap[0x4e] = READ_MODIFY_WRITE_TAG
	ServicesMap[0x52] = READ_TAG_FRAGMENTED
	ServicesMap[0x53] = WRITE_TAG_FRAGMENTED
}

func Build(service Service, path *bytes.Buffer, data *bytes.Buffer) *bytes.Buffer {
	var result bytes.Buffer

	pathLen := uint8(math.Ceil(float64(path.Len()) / 2))
	_ = binary.Write(&result, binary.BigEndian, service)
	_ = binary.Write(&result, binary.BigEndian, pathLen)
	_ = binary.Write(&result, binary.BigEndian, path.Bytes())
	if path.Len()%2 == 1 {
		_ = binary.Write(&result, binary.BigEndian, uint8(0))
	}
	_ = binary.Write(&result, binary.BigEndian, data.Bytes())

	return &result
}

type MessageRouter struct {
	Service              Service
	GeneralStatusCode    uint8
	ExtendedStatusLength uint8
	ExtendedStatus       []uint16
	data                 *bytes.Buffer
}

//func Parse(buf *bytes.Buffer) *MessageRouter{
//	result := MessageRouter{}
//	serviceCode := make([]byte,2)
//	_,_ = buf.Read(serviceCode)
//	generalStatusCode := make([]byte,1)
//	_,_ = buf.Read(generalStatusCode)
//	extendedStatusLength := make([]byte,1)
//	_,_ = buf.Read(extendedStatusLength)
//
//	result.Service = ServicesMap[serviceCode[0]]
//	result.GeneralStatusCode = generalStatusCode[0]
//	result.ExtendedStatusLength = extendedStatusLength[0]
//
//	for i := uint8(0); i < result.ExtendedStatusLength; i++ {
//		binary.
//		result.ExtendedStatus = append(result.ExtendedStatus, )
//	}
//
//	return &result
//}
