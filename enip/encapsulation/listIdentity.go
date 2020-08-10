package encapsulation

func (r *Request) ListIdentity() []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandListIdentity
	return pkg.Buffer()
}

func (r *Response) ListIdentity(data []byte) []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandListIdentity
	pkg.Data = data
	return pkg.Buffer()
}
