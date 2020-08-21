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
	cfg.Logger = log.New(os.Stdout, "", log.LstdFlags)

	tagA := tag.NewTag("TESTA")
	tagA.Onchange = func(value interface{}) {
		log.Println(value)
	}

	plc, _ := NewOriginator("10.211.55.7", 1, cfg)
	cfg.OnConnected = func() {
		log.Printf("%+v\n", plc.Controller)
		plc.ReadTag(tagA)
	}

	_ = plc.Connect()
	time.Sleep(time.Second * 600)
}
