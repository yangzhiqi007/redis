package model

import (
	"encoding/json"
	"io/ioutil"
)

type Database struct {
	ServerManager
	TaskManager
	SessionManager
}

func (self *Database) Load() error {
	data, err := ioutil.ReadFile(DBFileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, self)
	if err != nil {
		return err
	}

	return nil
}

func (self *Database) Save() error {
	data, err := json.MarshalIndent(self, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(DBFileName, data, 0666)
}

var (
	DB         Database
	DBFileName string
)
