package goplc

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestEncapsulation(t *testing.T) {
	cfg := DefaultConfig()

	a, _ := NewOriginator("10.211.55.7", 1, cfg)
	cfg.Log = log.New(os.Stdout, "", log.LstdFlags)
	cfg.OnAttribute = func() {
		log.Printf("%+v\n", a.Controller)
	}

	_ = a.Connect()
	time.Sleep(time.Second * 50)
}
