package http

import (
	"github.com/davyxu/golog"
)

var log = golog.New("httppeer")

func SetLevelByString(level string) {

	log.SetLevelByString(level)
}
