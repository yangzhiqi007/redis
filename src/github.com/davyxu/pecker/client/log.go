package client

import (
	"github.com/davyxu/golog"
)

var log = golog.New("client")

func init() {
	log.SetParts()
}
