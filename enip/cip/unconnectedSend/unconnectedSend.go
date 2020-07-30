package unconnectedSend

import (
	"bytes"
	"encoding/binary"
	"github.com/MiguelValentine/goplc/enip/cip/epath/segments/logical"
	"github.com/MiguelValentine/goplc/enip/cip/messageRouter"
	"math"
)

const UNCONNECTED_SEND_SERVICE = messageRouter.READ_TAG_FRAGMENTED

var UNCONNECTED_SEND_PATH bytes.Buffer

func init() {
	_ = binary.Write(&UNCONNECTED_SEND_PATH, binary.BigEndian, logical.Build(logical.ClassID, uint32(0x06), true).Bytes())
	_ = binary.Write(&UNCONNECTED_SEND_PATH, binary.BigEndian, logical.Build(logical.InstanceID, uint32(1), true).Bytes())
}

func GenerateEncodedTimeout(timeout uint64) (uint8, uint8) {
	diff := math.MaxFloat64
	j := float64(0)
	timeTick := float64(0)
	ticks := float64(0)
	for i := float64(0); i < 16; i++ {
		for j = i; j < 256; j++ {
			newDiff := math.Abs(float64(timeout) - math.Pow(2, i)*j)
			if newDiff <= diff {
				diff = newDiff
				timeTick = i
				ticks = j
			}
		}
	}
	return uint8(timeTick), uint8(ticks)
}

func Build(messageRequest *bytes.Buffer, path *bytes.Buffer, timeout uint64) *bytes.Buffer {
	if timeout < 100 {
		timeout = 1000
	}

	var buf bytes.Buffer
	timeTick, ticks := GenerateEncodedTimeout(timeout)
	_ = binary.Write(&buf, binary.BigEndian, timeTick)
	_ = binary.Write(&buf, binary.BigEndian, ticks)
	_ = binary.Write(&buf, binary.LittleEndian, uint16(messageRequest.Len()))
	_ = binary.Write(&buf, binary.BigEndian, messageRequest.Bytes())

	if messageRequest.Len()%2 == 1 {
		_ = binary.Write(&buf, binary.BigEndian, uint8(0))
	}

	pathLen := math.Ceil(float64(path.Len()) / 2)
	_ = binary.Write(&buf, binary.BigEndian, uint8(pathLen))
	_ = binary.Write(&buf, binary.BigEndian, uint8(0))
	_ = binary.Write(&buf, binary.BigEndian, path.Bytes())

	return messageRouter.Build(UNCONNECTED_SEND_SERVICE, &UNCONNECTED_SEND_PATH, &buf)
}
