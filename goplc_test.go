package goplc

import (
	"log"
	"net"
	"testing"
	"time"
)

func TestTest(t *testing.T) {
	go x()
	time.Sleep(time.Second * 2)
	a, b := New("127.0.0.1:10809", nil)
	err := a.Connect()
	log.Println(a, b, err)
	time.Sleep(time.Second * 5)
}

func x() {
	host := ":10809"

	// 获取tcp地址
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

	conn.Write([]byte("123"))
	conn.Close()

	conn2, _ := listener.AcceptTCP()
	conn2.Write([]byte("333"))
	conn2.Close()
}
