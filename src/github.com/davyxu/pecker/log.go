package main

import (
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func init() {
	log.SetParts()
}
