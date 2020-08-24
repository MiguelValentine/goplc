<div align="center">
    <img src="https://img.shields.io/github/go-mod/go-version/MiguelValentine/goplc?style=flat-square" alt="golang version">
    <img src="https://img.shields.io/github/license/MiguelValentine/goplc?style=flat-square" alt="MIT license"/>
    <img src="https://img.shields.io/github/stars/MiguelValentine/goplc.svg?&amp;style=social&amp;logo=github&amp;label=Stars" alt="GitHub stars">
</div>

# goplc

A golang based API(Originator) for interfacing with Rockwell Control/CompactLogix PLCs.

## Thanks

<a href="https://github.com/cmseaton42/node-ethernet-ip">node-ethernet-ip</a> project gave me a lot of inspiration。

## The API

### The Basics

#### Hello world

```go
import (
    "github.com/MiguelValentine/goplc"
)

func test(){
    //IP Slot Config
    plc, _ := NewOriginator("1.1.1.1", 1, nil)
    plc.OnConnected = func() {
        log.Printf("%+v\n",plc.Controller)
        // &{VendorID:1 DeviceType:14 ProductCode:53 Major:32 Minor:11 Status:12400 SerialNumber:1881423856 Version:32.11 Name:Emulator R32.11}
    }
    
    _ = plc.Connect()
}
```

#### Config

```go
import (
    "github.com/MiguelValentine/goplc"
)

func test(){
    cfg := DefaultConfig()

    //Reconnection Interval  default will not automatically reconnect
    cfg.ReconnectionInterval = time.Second

    //Logger  default no log
    cfg.Logger = log.New(os.Stdout, "", log.LstdFlags)

    //default 0xAF12
    cfg.Port = 1111

    plc, _ := NewOriginator("1.1.1.1", 1, cfg)
    plc.OnConnected = func() {
        log.Printf("%+v\n",plc.Controller)
        // &{VendorID:1 DeviceType:14 ProductCode:53 Major:32 Minor:11 Status:12400 SerialNumber:1881423856 Version:32.11 Name:Emulator R32.11}
    }
    
    _ = plc.Connect()
}
```

##### ListAllTags

```go
import (
    "github.com/MiguelValentine/goplc"
)

func testListAllTags() {
    cfg := DefaultConfig()
    plc, _ := NewOriginator("1.1.1.1", 1, cfg)
    plc.OnConnected = func() {
        plc.ListAllTags(0)
        //  2020/08/24 17:20:10 Map:Local : (0x1069)
        //  2020/08/24 17:20:10 Task:MainTask : (0x1070)
        //  2020/08/24 17:20:10 Program:MainProgram : (0x1068)
        //  2020/08/24 17:20:10 TESTA : DINT(0xc4)
        //  2020/08/24 17:20:10 TESTB : (0xafce)
        //  2020/08/24 17:20:10 TESTC : DINT(0xc4)
        //  2020/08/24 17:20:10 TESTD : REAL(0xca)
        //  2020/08/24 17:20:10 TESTE : (0x8fce)
    }

    _ = plc.Connect()
}
```

#### ReadAndWriteTag

```go
import (
    "github.com/MiguelValentine/goplc"
    "github.com/MiguelValentine/goplc/tag"
)

func testTag() {
    tagA := tag.NewTag("TESTA")

    //OnData will be called every time you sync
    //Try to use OnChange if not necessary
    tagA.OnData = func(value interface{}) {
        log.Println(value)
        //1237
        //1238
    }

    plc, _ := NewOriginator("1.1.1.1", 1, nil)
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
```

#### TagGroup

```go
import (
    "github.com/MiguelValentine/goplc"
    "github.com/MiguelValentine/goplc/tag"
    "github.com/MiguelValentine/goplc/tagGroup"
)

func testTagGroup() {
    cfg := DefaultConfig()

    tagA := tag.NewTag("TESTA")
    tagA.OnChange = func(value interface{}) {
        log.Println(value)
    }
    tagB := tag.NewTag("TESTB")
    tagB.OnChange = func(value interface{}) {
        log.Println(value)
    }
    tg := &tagGroup.TagGroup{}
    tg.Add(tagA)
    tg.Add(tagB)

    plc, _ := NewOriginator("1.1.1.1", 1, cfg)
    plc.OnConnected = func() {
        plc.ReadTagGroup(tg)
    }

    _ = plc.Connect()
}
```

#### Sync

```go
import (
    "github.com/MiguelValentine/goplc"
    "github.com/MiguelValentine/goplc/tag"
    "github.com/MiguelValentine/goplc/tagGroup"
)

func testTagGroup() {
    cfg := DefaultConfig()

    tagA := tag.NewTag("TESTA")
    tagA.OnChange = func(value interface{}) {
        log.Printf("%s => %f\n", tagA.Name(), value)
    }

    tagB := tag.NewTag("TESTB")
    tagB.OnChange = func(value interface{}) {
        log.Printf("%s => %f\n", tagB.Name(), value)
    }
    tg := &tagGroup.TagGroup{}
    tg.Add(tagA)
    tg.Add(tagB)

    plc, _ := NewOriginator("1.1.1.1", 1, cfg)
    plc.OnConnected = func() {
        //setTimeInterval
        plc.ReadTagGroupInterval(tg, time.Millisecond*50)
    }

    _ = plc.Connect()
}
```


## License

This project is licensed under the MIT License - see the [LICENCE](https://github.com/MiguelValentine/goplc/blob/master/LICENSE) file for details