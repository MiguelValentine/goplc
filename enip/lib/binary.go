package lib

import (
	"encoding/binary"
	"io"
)

func WriteByte(writer io.Writer, target interface{}) {
	err := binary.Write(writer, binary.LittleEndian, target)
	if err != nil {
		panic(err)
	}
}

func ReadByte(reader io.Reader, target interface{}) {
	err := binary.Read(reader, binary.LittleEndian, target)
	if err != nil {
		panic(err)
	}
}
