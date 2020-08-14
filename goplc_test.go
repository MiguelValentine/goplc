package goplc

import (
	"github.com/MiguelValentine/goplc/tag"
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	cfg := DefaultConfig()
	cfg.Log = log.New(os.Stdout, "", log.LstdFlags)

	a, _ := NewOriginator("10.211.55.7", 1, cfg)
	cfg.OnAttribute = func() {
		log.Printf("%+v\n", a.Controller)
		_tag := tag.NewTag([]byte("TESTA"))
		a.ReadTag(_tag)
	}

	_ = a.Connect()
	time.Sleep(time.Second * 600)

	//
	//b := tag.GenerateReadMessageRequest()
	//log.Printf("%#x\n", b)
}
