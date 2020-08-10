package encapsulation

func (r *Request) SendRRData(context uint64, data []byte) []byte {
	pkg := &Encapsulation{}
	pkg.Command = CommandSendRRData
	pkg.SenderContext = context
	return pkg.Buffer()
}
