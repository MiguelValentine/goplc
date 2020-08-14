package goplc

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/MiguelValentine/goplc/enip/cip"
	"github.com/MiguelValentine/goplc/enip/cip/epath/segment"
	"github.com/MiguelValentine/goplc/enip/cip/messageRouter"
	"github.com/MiguelValentine/goplc/enip/cip/unconnectedSend"
	"github.com/MiguelValentine/goplc/enip/encapsulation"
	"github.com/MiguelValentine/goplc/enip/etype"
	"github.com/MiguelValentine/goplc/enip/lib"
	"github.com/MiguelValentine/goplc/tag"
	"io"
	"log"
	"math/rand"
	"net"
	"time"
)

type controller struct {
	VendorID     etype.XUINT
	DeviceType   etype.XUINT
	ProductCode  etype.XUINT
	Major        uint8
	Minor        uint8
	Status       uint16
	SerialNumber uint32
	Version      string
	Name         string
}

type plc struct {
	tcpAddr     *net.TCPAddr
	tcpConn     *net.TCPConn
	config      *Config
	sender      chan []byte
	context     uint64
	session     etype.XUDINT
	slot        uint8
	request     *encapsulation.Request
	path        []byte
	Controller  *controller
	writeHandle bool
	tagsContext map[uint64]*tag.Tag
}

func (p *plc) Connect() error {
	p.config.Println("Connecting...")
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
	if p.config.OnConnected != nil {
		p.config.OnConnected()
	}

	p.config.Println("PLC Connected!")
	p.config.EBF.Clean()

	if !p.writeHandle {
		go p.write()
	}

	go p.read()

	p.registerSession()
}

func (p *plc) registerSession() {
	p.config.Println("Register Session")
	p.sender <- p.request.RegisterSession(p.context)
}

func (p *plc) readControllerProps() {
	mr := messageRouter.Build(messageRouter.ServiceGetAttributeAll, [][]byte{
		segment.LogicalBuild(segment.LogicalTypeClassID, 0x01, true),
		segment.LogicalBuild(segment.LogicalTypeInstanceID, 0x01, true),
	}, nil)
	ucmm := unconnectedSend.Build(mr, p.path, 2000)
	p.sender <- p.request.SendRRData(p.context, p.session, 10, ucmm)
}

func (p *plc) ReadTag(tag *tag.Tag) {
	p.tagsContext[tag.Context] = tag
	ucmm := unconnectedSend.Build(tag.GenerateReadMessageRequest(), p.path, 2000)
	p.sender <- p.request.SendRRData(tag.Context, p.session, 10, ucmm)
}

func (p *plc) disconnected(err error) {
	if p.config.OnDisconnected != nil {
		p.config.OnDisconnected(err)
	}

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
		p.config.Println("Reconnecting...")
		time.Sleep(p.config.ReconnectionInterval)
		err := p.Connect()
		if err != nil {
			panic(err)
		}
	}
}

func (p *plc) write() {
	p.writeHandle = true
	for {
		select {
		case data := <-p.sender:
			_, _ = p.tcpConn.Write(data)
		}
	}
}

func (p *plc) read() {
	buf := make([]byte, 1024*64)
	var err error
	for {
		var length int
		length, err = p.tcpConn.Read(buf)
		if err != nil {
			break
		}

		err = p.config.EBF.Read(buf[0:length], p.encapsulationHandle)
		if err != nil {
			break
		}
	}

	go p.disconnected(err)
}

func (p *plc) encapsulationHandle(_encapsulation *encapsulation.Encapsulation) error {
	switch _encapsulation.Command {
	case encapsulation.CommandRegisterSession:
		p.session = _encapsulation.SessionHandle
		p.config.Printf("session=> %d\n", p.session)
		if p.config.OnRegistered != nil {
			p.config.OnRegistered()
		}
		p.readControllerProps()
	case encapsulation.CommandUnRegisterSession:
		return errors.New("UnRegisterSession")
	case encapsulation.CommandSendRRData:
		p.config.Printf("SendRRData=> %d\n", _encapsulation.Length)
		_, _cpf := cip.Parser(_encapsulation.Data)
		return p.sendRRDataHandle(_encapsulation, _cpf)
	case encapsulation.CommandSendUnitData:
	}

	return nil
}

func (p *plc) sendRRDataHandle(_encapsulation *encapsulation.Encapsulation, cpf *cip.CPF) error {
	mr := messageRouter.Parse(cpf.Items[1].Data)
	if mr.GeneralStatus != 0 {
		return errors.New(fmt.Sprintf("SendRRData Error %#x => %s\n", mr.GeneralStatus, string(mr.AdditionalStatus)))
	}
	switch mr.Service - 0x80 {
	case messageRouter.ServiceGetAttributeAll:
		p.getAttributeAllHandle(mr.ResponseData)
	case messageRouter.ServiceReadTag:
		_tag := p.tagsContext[_encapsulation.SenderContext]
		if _tag == nil {
			return errors.New("TAG not found")
		}
		dataReader := bytes.NewReader(mr.ResponseData)
		lib.ReadByte(dataReader, &_tag.XType)
		_tag.Value = make([]byte, dataReader.Len())
		lib.ReadByte(dataReader, _tag.Value)
		v := _tag.GetValue()
		log.Println(v)
	}
	return nil
}

func (p *plc) getAttributeAllHandle(data []byte) {
	dataReader := bytes.NewReader(data)
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

	if p.config.OnAttribute != nil {
		p.config.OnAttribute()
	}
}

func NewOriginator(addr string, slot uint8, cfg *Config) (*plc, error) {
	_plc := &plc{}
	_plc.slot = slot
	_plc.config = cfg
	if _plc.config == nil {
		_plc.config = defaultConfig
	}

	_tcp, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, _plc.config.ENIP_PORT))
	if err != nil {
		return nil, err
	}

	_plc.tcpAddr = _tcp
	_plc.request = &encapsulation.Request{}

	_plc.path = segment.PortBuild(1, []byte{slot})

	rand.Seed(time.Now().Unix())
	_plc.context = rand.Uint64()
	_plc.config.Printf("Random context: %d\n", _plc.context)
	_plc.Controller = &controller{}
	_plc.writeHandle = false
	_plc.tagsContext = make(map[uint64]*tag.Tag)

	_plc.sender = make(chan []byte)
	return _plc, nil
}
