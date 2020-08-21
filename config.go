package goplc

import (
	"log"
	"time"
)

type Config struct {
	Port                 uint16
	ReconnectionInterval time.Duration
	Logger               *log.Logger
	OnConnected          func()
}

var defaultConfig *Config

func DefaultConfig() *Config {
	_defaultConfig := &Config{}
	_defaultConfig.Port = 0xAF12
	_defaultConfig.ReconnectionInterval = 0
	_defaultConfig.Logger = nil

	return _defaultConfig
}

func (c *Config) Println(v ...interface{}) {
	if c.Logger != nil {
		c.Logger.Println(v...)
	}
}

func (c *Config) Printf(format string, v ...interface{}) {
	if c.Logger != nil {
		c.Logger.Printf(format, v...)
	}
}

func init() {
	defaultConfig = DefaultConfig()
}
