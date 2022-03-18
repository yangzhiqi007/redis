package redisdrv

import "github.com/mediocregopher/radix.v2/redis"

func CmdQuery(c *RedisClient, cmd string, dataPtr interface{}, args ...interface{}) {

	if log.IsDebugEnabled() {
		log.Debugf("[DB] CmdQuery, %s %s    --->%s", cmd, argsToString(args), getStackFileString(2))
	}

	resp := WrapCommand(c, cmd, args...)

	RespToAny(resp, dataPtr)
}

func WrapCommand(c *redis.Client, cmd string, args ...interface{}) *redis.Resp {
	var resp *redis.Resp
	if NeedWrap() {

		if len(args) > 0 {
			if key, ok := args[0].(string); ok {

				var newValue = []interface{}{WrapKey(key)}
				for i := 1; i < len(args); i++ {
					newValue = append(newValue, args[i])
				}

				if ABaseLog {
					zjLog.Debugf("[DB] zj Cmd %s %s    --->%s", cmd, argsToString(newValue), getStackFileString(2))
				}

				resp = c.Cmd(cmd, newValue...)
			}
		}
	}

	if resp == nil {

		if ABaseLog {
			zjLog.Debugf("[DB] Cmd %s %s    --->%s", cmd, argsToString(args), getStackFileString(2))
		}
		resp = c.Cmd(cmd, args...)
	}

	return resp
}

func WrapPipeAppend(c *redis.Client, cmd string, args ...interface{}) {

	if NeedWrap() {

		if len(args) > 0 {
			if key, ok := args[0].(string); ok {

				var newValue = []interface{}{WrapKey(key)}
				for i := 1; i < len(args); i++ {
					newValue = append(newValue, args[i])
				}

				if ABaseLog {
					zjLog.Debugf("[DB] zj PipeAppend %s %s    --->%s", cmd, argsToString(newValue), getStackFileString(2))
				}
				c.PipeAppend(cmd, newValue...)
				return
			}
		}
	}

	if ABaseLog {
		zjLog.Debugf("[DB] PipeAppend %s %s    --->%s", cmd, argsToString(args), getStackFileString(2))
	}
	c.PipeAppend(cmd, args...)
}

func CmdExec(c *RedisClient, cmd string, args ...interface{}) {

	if log.IsDebugEnabled() {
		log.Debugf("[DB] CmdExec, %s %s    --->%s", cmd, argsToString(args), getStackFileString(2))
	}

	resp := WrapCommand(c, cmd, args...)

	if resp.Err != nil {
		log.Errorln("[DB] CmdExec failed:", resp.Err)
		panic(ErrCode_DBFailed)
	}

}

func PipeAppend(c *RedisClient, cmd string, args ...interface{}) {

	if log.IsDebugEnabled() {
		log.Debugf("[DB] PipeAppend, %s %s    --->%s", cmd, argsToString(args), getStackFileString(2))
	}

	WrapPipeAppend(c, cmd, args...)
}

// 批量执行
func PipeExec(c *RedisClient, callback func()) {

	c.PipeClear()

	callback()

	count := c.PendingCount()

	if log.IsDebugEnabled() {
		log.Debugf("[DB] PipeExec, %d    --->%s", count, getStackFileString(2))
	}

	for i := 0; i < count; i++ {

		if resp := c.PipeResp(); resp.Err != nil {
			log.Errorln("[DB] PipeExec failed:", resp.Err)
			panic(ErrCode_DBFailed)
		}
	}

	return

}
