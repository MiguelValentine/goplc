package encapsulation

func (r *Request) ListServices(context uint64) []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandListServices
	pkg.SenderContext = context

	return pkg.Buffer()
}
