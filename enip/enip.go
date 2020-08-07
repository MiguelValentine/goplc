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

func (e *EBF) Read(data []byte) ([]*encapsulation.Encapsulation, error) {
	e.bufferData = append(e.bufferData, data...)
	encapsulations, read, err := encapsulation.Parse(e.bufferData)
	e.bufferData = e.bufferData[read:]
	return encapsulations, err
}
