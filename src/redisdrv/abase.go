package redisdrv

import (
	"fmt"
	"github.com/davyxu/golog"
)

var (
	ABaseTable string // [xxx]
	ABaseLog   bool
	zjLog      = golog.New("zjabase")
)

func WrapKey(key string) string {
	if ABaseTable == "" {
		return key
	}

	return ABaseTable + key
}

func NeedWrap() bool {
	return ABaseTable != ""
}

func InitWrapKey(tableName string, enableLog bool) {
	if tableName == "" {
		return
	}

	ABaseTable = fmt.Sprintf("[%s]", tableName)
	ABaseLog = enableLog
	log.Debugf("zj abase table: %s, log: %v", tableName, enableLog)
}
