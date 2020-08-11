package enip

import (
	"github.com/MiguelValentine/goplc/enip/encapsulation"
)

type EBF struct {
	bufferData []byte
}

func (e *EBF) Clean() {
	e.bufferData = []byte{}
}

func (e *EBF) Read(data []byte, handle func(*encapsulation.Encapsulation)) error {
	e.bufferData = append(e.bufferData, data...)
	read, err := encapsulation.Parse(e.bufferData, handle)
	e.bufferData = e.bufferData[read:]
	return err
}
