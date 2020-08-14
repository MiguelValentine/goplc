package unconnectedSend

import (
	"bytes"
	"github.com/MiguelValentine/goplc/enip/cip/epath/segment"
	"github.com/MiguelValentine/goplc/enip/cip/messageRouter"
	"github.com/MiguelValentine/goplc/enip/etype"
	"github.com/MiguelValentine/goplc/enip/lib"
	"math"
)

type messageRequest struct {
	service         etype.XUSINT
	requestPathSize etype.XUSINT
	requestPath     []byte
	requestData     []byte
}

type UCMMSend struct {
	timeTick           uint8
	timeoutTick        etype.XUSINT
	messageRequestSize etype.XUINT
	messageRequest
	pad           etype.XUSINT
	routePathSize etype.XUSINT
	reserved      etype.XUSINT
	routePath     []byte
}

var UCMMSendPath [][]byte

func init() {
	UCMMSendPath = [][]byte{
		segment.LogicalBuild(segment.LogicalTypeClassID, 0x06, true),
		segment.LogicalBuild(segment.LogicalTypeInstanceID, 0x01, true),
	}
}

func generateEncodedTimeout(timeout uint64) (uint8, uint8) {
	diff := math.MaxFloat64
	timeTick := uint16(0)
	ticks := uint16(0)

	for i := uint16(0); i < 16; i++ {
		for j := uint16(1); j < 256; j++ {
			newDiff := math.Abs(float64(timeout) - math.Pow(float64(2), float64(i))*float64(j))
			if newDiff <= diff {
				diff = newDiff
				timeTick = i
				ticks = j
			}
		}
	}

	return uint8(timeTick), uint8(ticks)
}

func Build(data []byte, path []byte, timeout uint64) []byte {
	timeTick, ticks := generateEncodedTimeout(timeout)

	buffer := new(bytes.Buffer)
	lib.WriteByte(buffer, timeTick)
	lib.WriteByte(buffer, ticks)

	msgLength := uint16(len(data))
	lib.WriteByte(buffer, msgLength)
	lib.WriteByte(buffer, data)
	if msgLength%2 == 1 {
		lib.WriteByte(buffer, uint8(0))
	}

	pathLength := uint8(len(path) / 2)
	lib.WriteByte(buffer, pathLength)
	lib.WriteByte(buffer, uint8(0))
	lib.WriteByte(buffer, path)

	return messageRouter.Build(messageRouter.ServiceUnconnectedSend, UCMMSendPath, buffer.Bytes())
}
