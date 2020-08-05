package enip

import (
	"bytes"
	"encoding/binary"
	"log"
)

func writeBinary(src *bytes.Buffer, data interface{}) {
	e := binary.Write(src, binary.LittleEndian, data)
	if e != nil {
		log.Printf("%s : %+v\n", "Binary write error", data)
		log.Fatalln(e)
	}
}

func toBinary(data interface{}) *bytes.Buffer {
	var _data bytes.Buffer
	writeBinary(&_data, data)
	return &_data
}
