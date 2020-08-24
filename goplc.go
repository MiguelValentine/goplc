package goplc

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/MiguelValentine/goplc/ethernetip"
	"github.com/MiguelValentine/goplc/ethernetip/commonIndustrialProtocol"
	"github.com/MiguelValentine/goplc/ethernetip/commonIndustrialProtocol/segment"
	"github.com/MiguelValentine/goplc/ethernetip/commonIndustrialProtocol/segment/epath"
	_type "github.com/MiguelValentine/goplc/ethernetip/type"
	"github.com/MiguelValentine/goplc/lib"
	"github.com/MiguelValentine/goplc/tag"
	"github.com/MiguelValentine/goplc/tagGroup"
	"io"
	"math"
	"math/rand"
	"net"
	"time"
)

type controller struct {
	VendorID     _type.UINT
	DeviceType   _type.UINT
	ProductCode  _type.UINT
	Major        _type.USINT
	Minor        _type.USINT
	Status       _type.UINT
	SerialNumber _type.UDINT
	Version      string
	Name         string
}

type plc struct {
	tcpAddr    *net.TCPAddr
	tcpConn    *net.TCPConn
	config     *Config
	session    _type.UDINT
	Controller *controller

	writeRoute bool
	sender     chan []byte
	bufferData []byte
	TargetPath []byte
	HandleMap  map[ethernetip.Command]func(*ethernetip.Encapsulation)

	ContextPool map[_type.ULINT]func(*commonIndustrialProtocol.MessageRouterResponse)
}

func (p *plc) Connect() error {
	p.config.Println("PLC Connecting...")
	_conn, err := net.DialTCP("tcp", nil, p.tcpAddr)
	if err != nil {
		return err
	}

	err2 := _conn.SetKeepAlive(true)
	if err2 != nil {
		return err2
	}

	p.tcpConn = _conn
	p.connected()
	return nil
}

func (p *plc) connected() {
	p.config.Println("PLC Connected!")
	p.bufferData = []byte{}

	if !p.writeRoute {
		go p.write()
	}

	go p.read()

	p.config.Println("Register Session...")
	encapsulation := ethernetip.RequestRegisterSession(0)

	p.sender <- encapsulation.Buffer()
}

func (p *plc) disconnected(err error) {
	if err == io.EOF {
		p.config.Println("PLC Disconnected!")
		p.config.Println("EOF")
	} else {
		p.config.Println("PLC Disconnected!")
		p.config.Println(err)
	}

	_ = p.tcpConn.Close()
	p.tcpConn = nil

	if p.config.ReconnectionInterval != 0 {
		time.Sleep(p.config.ReconnectionInterval)
		p.config.Println("PLC Reconnecting!")
		err := p.Connect()
		if err != nil {
			panic(err)
		}
	}
}

func (p *plc) read() {
	defer func() {
		if err := recover(); err != nil {
			go p.disconnected(err.(error))
		}
	}()

	buf := make([]byte, 1024*64)
	var err error
	for {
		var length int
		length, err = p.tcpConn.Read(buf)
		if err != nil {
			break
		}
		p.config.Printf("Receive => %d bytes\n", length)

		p.bufferData = append(p.bufferData, buf[0:length]...)
		if len(p.bufferData) > 24 {
			read, encapsulations := ethernetip.Slice(p.bufferData)
			p.bufferData = p.bufferData[read:]

			for _, encapsulation := range encapsulations {
				if encapsulation.Status == ethernetip.StatusSuccess {
					if exec, ok := p.HandleMap[encapsulation.Command]; ok {
						exec(encapsulation)
					} else {
						p.config.Printf("Received encapsulation Command: %#x ,but no registered handler!\n", encapsulation.Command)
					}
				}
			}
		}
	}
}

func (p *plc) handleRegisterSession(encapsulation *ethernetip.Encapsulation) {
	p.session = encapsulation.SessionHandle
	p.config.Printf("Session => %#x\n", p.session)

	// get_attribute_all
	mr1 := &commonIndustrialProtocol.MessageRouterRequest{}
	mr1.Service = 0x01
	mr1.RequestPath = segment.Paths(
		epath.LogicalBuild(epath.LogicalTypeClassID, 01, true),
		epath.LogicalBuild(epath.LogicalTypeInstanceID, 01, true),
	)

	p.UcmmSend(3, 250, math.MaxUint64, mr1)
	p.ContextPool[math.MaxUint64] = p.getAttributeAll
}

func (p *plc) getAttributeAll(mr *commonIndustrialProtocol.MessageRouterResponse) {
	p.config.Printf("%+v\n", mr)

	dataReader := bytes.NewReader(mr.ResponseData)
	lib.ReadByte(dataReader, &p.Controller.VendorID)
	lib.ReadByte(dataReader, &p.Controller.DeviceType)
	lib.ReadByte(dataReader, &p.Controller.ProductCode)
	lib.ReadByte(dataReader, &p.Controller.Major)
	lib.ReadByte(dataReader, &p.Controller.Minor)
	lib.ReadByte(dataReader, &p.Controller.Status)
	lib.ReadByte(dataReader, &p.Controller.SerialNumber)
	nameLen := uint8(0)
	lib.ReadByte(dataReader, &nameLen)
	nameBuf := make([]byte, nameLen)
	lib.ReadByte(dataReader, nameBuf)

	p.Controller.Name = string(nameBuf)
	p.Controller.Version = fmt.Sprintf("%d.%d", p.Controller.Major, p.Controller.Minor)

	if p.config.OnConnected != nil {
		p.config.OnConnected()
	}
}

func (p *plc) UcmmSend(timeTicks _type.USINT, timeoutTicks _type.USINT, context _type.ULINT, mr1 *commonIndustrialProtocol.MessageRouterRequest) {
	ucmm := &commonIndustrialProtocol.UnconnectedSend{}
	ucmm.TimeTick = timeTicks
	ucmm.TimeOutTicks = timeoutTicks
	ucmm.MessageRequest = mr1
	ucmm.RouterPath = p.TargetPath

	mr2 := &commonIndustrialProtocol.MessageRouterRequest{}
	mr2.Service = 0x52
	mr2.RequestPath = segment.Paths(
		epath.LogicalBuild(epath.LogicalTypeClassID, 06, true),
		epath.LogicalBuild(epath.LogicalTypeInstanceID, 01, true),
	)
	mr2.RequestData = ucmm.Buffer()

	cpf := &ethernetip.CommonPacketFormat{}
	cpf.UnconnectedData(mr2.Buffer())
	pkg := ethernetip.RequestSendRRData(p.session, context, 10, cpf)

	p.sender <- pkg.Buffer()
}

func (p *plc) handleSendData(encapsulation *ethernetip.Encapsulation) {
	cpf := ethernetip.SendRRDataParser(encapsulation.Data)
	mr := commonIndustrialProtocol.MRParser(cpf.DataItem.Data)
	if mr.GeneralStatus != 0 {
		panic(errors.New(fmt.Sprintf("target error => Service Code: %#x | Status: %#x | Addtional: %s", mr.ReplyService, mr.GeneralStatus, mr.AdditionalStatus)))
	} else {
		p.ContextPool[encapsulation.SenderContext](mr)
		delete(p.ContextPool, encapsulation.SenderContext)
	}
}

func (p *plc) write() {
	p.writeRoute = true
	for {
		select {
		case data := <-p.sender:
			p.config.Printf("Send => %d bytes\n", len(data))
			_, _ = p.tcpConn.Write(data)
		}
	}
}

func (p *plc) ReadTag(tag *tag.Tag) *tag.Tag {
	rand.Seed(time.Now().UnixNano())
	context := _type.ULINT(rand.Uint64())
	p.ContextPool[context] = tag.ReadTagParser
	p.UcmmSend(3, 250, context, tag.GenerateReadMessageRequest())
	return tag
}

func (p *plc) WriteTag(tag *tag.Tag) *tag.Tag {
	rand.Seed(time.Now().UnixNano())
	context := _type.ULINT(rand.Uint64())
	p.ContextPool[context] = tag.WriteTagParser
	p.UcmmSend(3, 250, context, tag.GenerateWriteMessageRequest())
	return tag
}

func (p *plc) ReadTagGroup(tg *tagGroup.TagGroup) {
	rand.Seed(time.Now().UnixNano())
	context := _type.ULINT(rand.Uint64())
	p.ContextPool[context] = tg.ReadTagParser
	p.UcmmSend(3, 250, context, tg.GenerateReadMessageRequest())
}

func (p *plc) ReadTagGroupInterval(tg *tagGroup.TagGroup, d time.Duration) {
	lib.Cron(d, func() {
		p.ReadTagGroup(tg)
	})
}

func NewOriginator(addr string, slot uint8, cfg *Config) (*plc, error) {
	_plc := &plc{}
	_plc.config = cfg
	_plc.config = cfg
	if _plc.config == nil {
		_plc.config = defaultConfig
	}

	_tcp, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, _plc.config.Port))
	if err != nil {
		return nil, err
	}

	_plc.tcpAddr = _tcp
	_plc.Controller = &controller{}
	_plc.sender = make(chan []byte)
	_plc.HandleMap = make(map[ethernetip.Command]func(*ethernetip.Encapsulation))
	_plc.TargetPath = epath.PortBuild([]byte{slot}, 1, true)
	_plc.ContextPool = make(map[_type.ULINT]func(*commonIndustrialProtocol.MessageRouterResponse))

	_plc.HandleMap[ethernetip.CommandNOP] = ethernetip.HandleNop
	_plc.HandleMap[ethernetip.CommandListIdentity] = ethernetip.HandleListIdentity
	_plc.HandleMap[ethernetip.CommandListInterfaces] = ethernetip.HandleListInterfaces
	_plc.HandleMap[ethernetip.CommandRegisterSession] = _plc.handleRegisterSession
	_plc.HandleMap[ethernetip.CommandSendRRData] = _plc.handleSendData
	_plc.HandleMap[ethernetip.CommandSendUnitData] = _plc.handleSendData

	return _plc, nil
}
