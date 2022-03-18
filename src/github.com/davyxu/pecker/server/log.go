package server

import (
	"github.com/davyxu/golog"
)

var log = golog.New("server")

func init() {
	log.SetParts(golog.LogPart_Time)
}
