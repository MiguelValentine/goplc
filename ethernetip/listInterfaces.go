package ethernetip

func RequestListInterfaces() *Encapsulation {
	encapsulation := &Encapsulation{}
	encapsulation.Command = CommandListInterfaces

	return encapsulation
}

func HandleListInterfaces(encapsulation *Encapsulation) {

}
