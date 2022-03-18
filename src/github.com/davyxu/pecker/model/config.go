package model

import (
	"encoding/json"
	"io/ioutil"
)

type AddressAliasDefine struct {
	Addr  string `json:"addr"`
	Alias string `json:"alias"`
}

type Config struct {
	AddrAlias []AddressAliasDefine `json:"addr_alias"`
}

func (self *Config) GetAddress(addrStr string) string {

	for _, def := range self.AddrAlias {
		if def.Alias == addrStr {
			return def.Addr
		}
	}
	return addrStr
}

func (self *Config) Load(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, self)
	if err != nil {
		return err
	}

	return nil
}
