package model

import (
	"github.com/davyxu/golog"
)

var log = golog.New("model")

func init() {
	log.SetParts()
}
