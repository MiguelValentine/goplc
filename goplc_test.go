package goplc

import (
	"github.com/MiguelValentine/goplc/tag"
	"github.com/MiguelValentine/goplc/tagGroup"
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	//testListAllTags()
	//testListTemplate()
	//testTagGroupInterval()
	//testTagGroup()
	testTag()
	time.Sleep(time.Second * 600)
}

func testListTemplate() {
	cfg := DefaultConfig()
	plc, _ := NewOriginator("10.211.55.7", 1, cfg)
	plc.OnConnected = func() {
		plc.ListTemplate(0xafce)
	}

	_ = plc.Connect()
}

func testListAllTags() {
	cfg := DefaultConfig()
	plc, _ := NewOriginator("10.211.55.7", 1, cfg)
	plc.OnConnected = func() {
		plc.ListAllTags(0)
	}

	_ = plc.Connect()
}

func testTagGroupInterval() {
	cfg := DefaultConfig()
	//cfg.Logger = log.New(os.Stdout, "", log.LstdFlags)

	tagA := tag.NewTag("TESTE")
	tagA.OnChange = func(value interface{}) {
		log.Println(value)
	}

	tagC := tag.NewTag("TESTD")
	tagC.OnChange = func(value interface{}) {
		log.Printf("%s => %f\n", tagC.Name(), value)
	}
	tg := &tagGroup.TagGroup{}
	tg.Add(tagA)
	tg.Add(tagC)

	plc, _ := NewOriginator("10.211.55.7", 1, cfg)
	plc.OnConnected = func() {
		plc.ReadTagGroupInterval(tg, time.Millisecond*50)
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
	tagC := tag.NewTag("TESTD")
	tagC.OnChange = func(value interface{}) {
		log.Println(value)
	}
	tg := &tagGroup.TagGroup{}
	tg.Add(tagA)
	tg.Add(tagC)

	plc, _ := NewOriginator("10.211.55.7", 1, cfg)
	plc.OnConnected = func() {
		plc.ReadTagGroup(tg)
	}

	_ = plc.Connect()
}

func testTag() {
	//cfg := DefaultConfig()
	//cfg.Logger = log.New(os.Stdout, "", log.LstdFlags)

	tagA := tag.NewTag("TESTA")
	tagA.OnData = func(value interface{}) {
		log.Println(value)
	}

	plc, _ := NewOriginator("10.211.55.7", 1, nil)
	plc.OnConnected = func() {
		log.Printf("%+v\n", plc.Controller)
		plc.ReadTag(tagA).Then(func() {
			tagA.SetValue(int32(1238))
			plc.WriteTag(tagA).Then(func() {
				plc.ReadTag(tagA)
			})
		})
	}

	_ = plc.Connect()
}
