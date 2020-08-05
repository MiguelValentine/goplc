package enip

import (
	"bytes"
	"encoding/binary"
	"log"
)

func writeBinary(src *bytes.Buffer, data interface{}) {
	e := binary.Write(src, binary.LittleEndian, data)
	if e != nil {
		log.Printf("%s : %+v\n", "Binary Write Error", data)
		log.Fatalln(e)
	}
}
