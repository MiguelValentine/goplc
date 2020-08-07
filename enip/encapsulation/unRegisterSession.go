package encapsulation

import "github.com/MiguelValentine/goplc/enip/etype"

func (r *Request) UnRegisterSession(context uint64, session etype.XUDINT) []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandUnRegisterSession
	pkg.SenderContext = context
	pkg.SessionHandle = session

	return pkg.Buffer()
}
