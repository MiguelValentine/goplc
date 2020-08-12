package encapsulation

import (
	"github.com/MiguelValentine/goplc/enip/cip"
	"github.com/MiguelValentine/goplc/enip/etype"
)

func (r *Request) SendRRData(context uint64, session etype.XUDINT, timeout etype.XUINT, data []byte) []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandSendRRData
	pkg.SenderContext = context
	pkg.SessionHandle = session

	cpf := cip.NewCPF([]*cip.CPFItem{
		cip.NewCPFItem(cip.TypeNull, nil),
		cip.NewCPFItem(cip.TypeUnconnectedMessage, data),
	}, nil)

	pkg.Data = cip.Build(timeout, cpf)
	pkg.Length = etype.XUINT(len(pkg.Data))
	return pkg.Buffer()
}
