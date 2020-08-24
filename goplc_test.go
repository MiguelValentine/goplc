package goplc

import (
	"github.com/MiguelValentine/goplc/tag"
	"github.com/MiguelValentine/goplc/tagGroup"
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	testTagGroupInterval()
	time.Sleep(time.Second * 600)
}
func testTagGroupInterval() {
	cfg := DefaultConfig()
	//cfg.Logger = log.New(os.Stdout, "", log.LstdFlags)

	tagA := tag.NewTag("TESTA")
	tagA.OnChange = func(value interface{}) {
		log.Println(value)
	}
	tagC := tag.NewTag("TESTC")
	tagC.OnChange = func(value interface{}) {
		log.Printf("%s => %d\n", tagC.Name(), value)
	}
	tg := &tagGroup.TagGroup{}
	tg.Add(tagA)
	tg.Add(tagC)

	plc, _ := NewOriginator("10.211.55.7", 1, cfg)
	cfg.OnConnected = func() {
		plc.ReadTagGroupInterval(tg, time.Second)
	}

	_ = plc.Connect()
}

func testTagGroup() {
	cfg := DefaultConfig()
	//cfg.Logger = log.New(os.Stdout, "", log.LstdFlags)

	tagA := tag.NewTag("TESTA")
	tagA.OnChange = func(value interface{}) {
		log.Println(value)
	}
	tagC := tag.NewTag("TESTC")
	tagC.OnChange = func(value interface{}) {
		log.Println(value)
	}
	tg := &tagGroup.TagGroup{}
	tg.Add(tagA)
	tg.Add(tagC)

	plc, _ := NewOriginator("10.211.55.7", 1, cfg)
	cfg.OnConnected = func() {
		plc.ReadTagGroup(tg)
	}

	_ = plc.Connect()
}

func testTag() {
	cfg := DefaultConfig()
	//cfg.Logger = log.New(os.Stdout, "", log.LstdFlags)

	tagA := tag.NewTag("TESTA")
	tagA.OnChange = func(value interface{}) {
		log.Println(value)
	}

	plc, _ := NewOriginator("10.211.55.7", 1, cfg)
	cfg.OnConnected = func() {
		log.Printf("%+v\n", plc.Controller)
		plc.ReadTag(tagA).Then(func() {
			tagA.SetValue(int32(123))
			plc.WriteTag(tagA).Then(func() {
				plc.ReadTag(tagA)
			})
		})
	}

	_ = plc.Connect()
}
