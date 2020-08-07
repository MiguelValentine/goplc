package encapsulation

import "github.com/MiguelValentine/goplc/enip/cip"

func (r *Request) ListIdentity() []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandListIdentity
	return pkg.Buffer()
}

func (r *Response) ListIdentity(i *cip.Identity) []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandListIdentity
	pkg.Data = i.Buffer()
	return pkg.Buffer()
}
