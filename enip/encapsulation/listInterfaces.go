package encapsulation

func (r *Request) ListInterfaces() []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandListInterfaces
	return pkg.Buffer()
}

func (r *Response) ListInterfaces(data []byte) []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandListInterfaces
	pkg.Data = data
	return pkg.Buffer()
}
