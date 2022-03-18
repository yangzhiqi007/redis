package redisdrv

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

type RedisClient = redis.Client

func OperateDB(pool *pool.Pool, callback func(*RedisClient)) (code int) {

	if pool == nil {
		return ErrCode_DBFailed
	}

	c, err := pool.Get()
	if err != nil {
		log.Errorln("[DB] get client failed:", err)
		return ErrCode_DBFailed
	}

	defer func() {

		pool.Put(c)

		//drv.RemoveRef()

		//log.Warnf("%p remove ref %d %s", drv, ref, getStackFileString(3))

		switch err := recover().(type) {
		case int:
			code = err
		case nil:
		default:
			panic(err)
		}

	}()

	callback(c)

	return 0
}

func LoadObject(debugEnv string, objLists ...*ObjectList) {

	if len(objLists) == 0 {
		panic(ErrCode_DBFailed)
	}

	if log.IsDebugEnabled() {
		log.Debugf("[DB] LoadObject | %s", debugEnv)
	}

	one := objLists[0]

	one.c.PipeClear()

	for _, objList := range objLists {

		objList.loadCmd()
	}

	for _, objList := range objLists {

		err := objList.deserialize()

		if err != nil {
			log.Errorln(err)
			panic(ErrCode_DBFailed)
		}

	}

}

func SaveObject(debugEnv string, objLists ...*ObjectList) {

	if len(objLists) == 0 {
		panic(ErrCode_DBFailed)
	}

	one := objLists[0]

	if log.IsDebugEnabled() {
		log.Debugf("[DB] SaveObject | %s", debugEnv)
	}

	one.c.PipeClear()

	for _, objList := range objLists {

		objList.serialize()
	}

	count := one.c.PendingCount()

	for i := 0; i < count; i++ {
		resp := one.c.PipeResp()

		if resp.Err != nil {
			log.Errorln(resp.Err)
			panic(ErrCode_DBFailed)
		}
	}

}
