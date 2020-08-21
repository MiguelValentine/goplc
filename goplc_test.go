package goplc

import (
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	//cfg := DefaultConfig()
	//cfg.Log = log.New(os.Stdout, "", log.LstdFlags)
	//
	//a, _ := NewOriginator("10.211.55.7", 1, cfg)
	//cfg.OnAttribute = func() {
	//	log.Printf("%+v\n", a.Controller)
	//	_tag := tag.NewTag([]byte("TESTA"))
	//	a.ReadTag(_tag)
	//}
	//
	//_ = a.Connect()
	//time.Sleep(time.Second * 600)

	//
	//b := tag.GenerateReadMessageRequest()
	//log.Printf("%#x\n", b)

	//cpf := ethernetip.CommonPacketFormat{}
	//cpf.SequencedAddress(123, 333)
	//cpf.ConnectedData([]byte("Hello"))
	//
	//log.Printf("%+v\n", cpf)
	//log.Printf("% #x\n", cpf.Buffer())

	cfg := DefaultConfig()
	//cfg.Logger = log.New(os.Stdout, "", log.LstdFlags)

	plc, _ := NewOriginator("10.211.55.7", 1, cfg)

	cfg.OnConnected = func() {
		log.Printf("%+v\n", plc.Controller)
	}

	_ = plc.Connect()
	time.Sleep(time.Second * 600)
}
