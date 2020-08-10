package goplc

import (
	"fmt"
	"github.com/MiguelValentine/goplc/enip/encapsulation"
	"github.com/MiguelValentine/goplc/enip/etype"
	"io"
	"log"
	"math/rand"
	"net"
	"time"
)

type plc struct {
	tcpAddr *net.TCPAddr
	tcpConn *net.TCPConn
	config  *Config
	sender  chan []byte
	context uint64
	session etype.XUDINT
}

func (p *plc) Connect() error {
	log.Println("Connecting...")
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

	go p.read()
	go p.write()

	p.registerSession()
}

func (p *plc) registerSession() {
	r := encapsulation.Request{}
	p.sender <- r.RegisterSession(p.context)
}

func (p *plc) disconnected(err error) {
	if p.config.OnDisconnected != nil {
		p.config.OnDisconnected(err)
	}

	if err != io.EOF {
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
	for {
		select {
		case data := <-p.sender:
			log.Println(data)
			_, _ = p.tcpConn.Write(data)
		}
	}
}

func (p *plc) read() {
	buf := make([]byte, 1024*64)
	for {
		length, err := p.tcpConn.Read(buf)
		if err != nil {
			p.disconnected(err)
			break
		}

		encapsulations, err := p.config.EBF.Read(buf[0:length])
		if err != nil {
			p.disconnected(err)
			break
		}

		for _, _encapsulation := range encapsulations {
			switch _encapsulation.Command {
			case encapsulation.CommandRegisterSession:
				p.session = _encapsulation.SessionHandle
				p.config.OnRegistered()
			}
		}
	}
}

func NewOriginator(addr string, cfg *Config) (*plc, error) {
	_plc := &plc{}
	_plc.config = cfg
	if _plc.config == nil {
		_plc.config = defaultConfig
	}

	_tcp, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, _plc.config.ENIP_PORT))
	if err != nil {
		return nil, err
	}

	_plc.tcpAddr = _tcp

	rand.Seed(time.Now().Unix())
	_plc.context = rand.Uint64()
	log.Printf("Random context: %d\n", _plc.context)

	_plc.sender = make(chan []byte)
	return _plc, nil
}
