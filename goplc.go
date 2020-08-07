package goplc

import (
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
}

func (p *plc) Connect() error {
	_conn, err := net.DialTCP("tcp", nil, p.tcpAddr)
	if err != nil {
		return err
	}

	err2 := _conn.SetKeepAlive(true)
	if err2 != nil {
		return err2
	}

	p.tcpConn = _conn
	p.Connected()
	return nil
}

func (p *plc) Connected() {
	if p.config.OnConnected != nil {
		p.config.OnConnected()
	}

	p.config.Println("PLC Connected!")
	p.config.EBF.Clean()

	go p.read()
}

func (p *plc) Disconnected(err error) {
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

func (p *plc) read() {
	buf := make([]byte, 1024*64)
	for {
		length, err := p.tcpConn.Read(buf)
		if err != nil {
			p.Disconnected(err)
			break
		}

		encapsulations, err := p.config.EBF.Read(buf[0:length])
		if err != nil {
			p.Disconnected(err)
			break
		}

		if encapsulations != nil {
			log.Printf("%+v\n", encapsulations[0])
		}
	}
}

func New(addr string, cfg *Config) (*plc, error) {
	_tcp, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	_plc := &plc{}
	_plc.tcpAddr = _tcp
	_plc.config = cfg

	if _plc.config == nil {
		_plc.config = defaultConfig
	}

	rand.Seed(time.Now().Unix())
	_plc.context = rand.Uint64()

	_plc.sender = make(chan []byte)
	return _plc, nil
}
