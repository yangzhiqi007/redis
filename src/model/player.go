package model

import (
	"fmt"
	"log"
)

type GameModel struct {
}

func (self *GameModel) GetMainKey(env interface{}) interface{} {
	return fmt.Sprintf("m:%d", env) // teamid
}

//go:generate msgp

type A struct {
	GameModel `msg:"-"`
	Name      string `msg:"N"`
	Age       int8   `msg:"A"`
}

func (self *A) GetHashKey(env interface{}) interface{} {
	return "TestA"
}

func (self *A) Hsh() {
	log.Println("111111111111111111")
}

func (self *A) SSS() {
	log.Println("22222222222222222222")
}
