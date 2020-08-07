package encapsulation

func (r *Request) Nop() []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandNOP
	return pkg.Buffer()
}
