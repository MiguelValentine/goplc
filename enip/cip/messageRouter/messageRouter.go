package messageRouter

import "github.com/MiguelValentine/goplc/enip/etype"

type Service etype.XUSINT

const (
	ServiceGetAttributeAll       Service = 0x01
	ServiceSetAttributeAll       Service = 0x02
	ServiceGetAttributeList      Service = 0x03
	ServiceSetAttributeList      Service = 0x04
	ServiceReset                 Service = 0x05
	ServiceStart                 Service = 0x06
	ServiceStop                  Service = 0x07
	ServiceCreate                Service = 0x08
	ServiceDelete                Service = 0x09
	ServiceMultipleServicePacket Service = 0x0a
	ServiceApplyAttributes       Service = 0x0d
	ServiceGetAttributeSingle    Service = 0x0e
	ServiceSetAttributeSingle    Service = 0x10
	ServiceFindNext              Service = 0x11
	ServiceRestore               Service = 0x15
	ServiceSave                  Service = 0x16
	ServiceGetMember             Service = 0x18
	ServiceSetMember             Service = 0x19
	ServiceInsertMember          Service = 0x1A
	ServiceRemoveMember          Service = 0x1B
	ServiceGroupSync             Service = 0x1C
	ServiceReadTag               Service = 0x4c
	ServiceWriteTag              Service = 0x4d
	ServiceReadTagFragmented     Service = 0x52
	ServiceWriteTagFragmented    Service = 0x53
	FragmentedReadModifyWriteTag Service = 0x4e
)

type Request struct {
	service         Service
	requestPathSize etype.XUSINT
	requestPath     []byte
	requestData     []byte
}
