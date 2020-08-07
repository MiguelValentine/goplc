package goplc

import (
	"github.com/MiguelValentine/goplc/enip/encapsulation"
	"log"
	"net"
	"testing"
	"time"
)

func TestEncapsulation(t *testing.T) {
	go tcpServer()

	a, b := New("127.0.0.1:10809", nil)
	err := a.Connect()
	log.Println(a, b, err)
	time.Sleep(time.Second * 5)
}

func tcpServer() {
	host := ":10809"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		log.Printf("resolve tcp addr failed: %v\n", err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("listen tcp port failed: %v\n", err)
		return
	}
	conn, err := listener.AcceptTCP()
	if err != nil {
		log.Printf("Accept failed:%v\n", err)
	}

	data := &encapsulation.Encapsulation{}
	data.Command = encapsulation.CommandRegisterSession
	_buf := data.Buffer()
	log.Println(_buf)

	_, _ = conn.Write(_buf[0:3])
	time.Sleep(time.Millisecond * 200)
	_, _ = conn.Write(_buf[3:])
	conn.Close()
}
