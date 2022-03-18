package redisdrv

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	ErrCode_DBFailed = -1
)

func argsToString(args []interface{}) string {

	var sb strings.Builder
	for _, arg := range args {

		switch v := arg.(type) {
		case []byte:
			sb.WriteString(fmt.Sprintf("bytes(len:%d)", len(v)))
		case string:
			sb.WriteString("'")
			sb.WriteString(v)
			sb.WriteString("'")
		default:
			sb.WriteString(fmt.Sprintf("%v", arg))
		}

		sb.WriteString(" ")

	}

	return sb.String()
}

func getStackFileString(level int) string {

	_, file, line, ok := runtime.Caller(level)

	if !ok {
		return "??"
	}

	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}
