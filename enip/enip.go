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

func (e *EBF) Read(data []byte, handle func(*encapsulation.Encapsulation) error) error {
	e.bufferData = append(e.bufferData, data...)
	read, err, es := encapsulation.Parse(e.bufferData)

	e.bufferData = e.bufferData[read:]
	for _, _encapsulation := range es {
		err1 := handle(_encapsulation)
		if err1 != nil {
			return err1
		}
	}
	return err
}
