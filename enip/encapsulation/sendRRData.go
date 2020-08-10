package encapsulation

import "github.com/MiguelValentine/goplc/enip/etype"

func (r *Request) SendRRData(context uint64, data []byte) []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandSendRRData
	pkg.SenderContext = context
	return pkg.Buffer()
}

func (r *Response) SendRRData(context uint64, session etype.XUDINT, data []byte) []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandSendRRData
	pkg.SenderContext = context
	pkg.Data = data
	return pkg.Buffer()
}
