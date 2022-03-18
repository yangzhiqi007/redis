package memsd

import "time"

type Config struct {
	Address          string
	RequestTimeout   time.Duration
	PingCheckTimeout time.Duration
	PingInterval     time.Duration
}

func DefaultConfig() Config {

	return Config{
		Address:        ":8900",
		RequestTimeout: time.Second * 10,
	}
}

func (self *memDiscovery) GetConfig() Config {
	return self.config
}

func (self *memDiscovery) SetConfig(c Config) {
	self.config = c
}
