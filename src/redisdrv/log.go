package redisdrv

import (
	"github.com/davyxu/golog"
)

var log = golog.New("redisdrv")

func init() {
	log.SetLevelByString("info")
}
