package redisdrv

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
	"time"
)

const (
	TryLockTimeOutMS int64 = 2 * 1000               // 尝试获取锁超时时间(毫秒)
	TryLockInterval        = 300 * time.Millisecond // 尝试获取锁间隔时间
	LockKeyExpireSec       = 3                      // db中锁对应的key的过期时间(秒)

)

// 自1970年起的UTC绝对毫秒
func TimeEpochMS() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}

func IsLockOK(resp *redis.Resp) bool {
	if !resp.IsType(redis.Nil) {

		if v, err := resp.Str(); err != nil {
			log.Errorf("GetSerialRespond type failed, %s", err)
			return false
		} else {
			return v == "OK"
		}

	}

	return false
}

// 锁
type Lock struct {
	Key string // 锁在redis中的key， 规则: "lock:xxx"
}

func GetKey_Lock(name string) string {
	return fmt.Sprintf("lock:%s", name)
}

func TryLock(c *RedisClient, key string) (*Lock, int) {

	// 锁在redis中key的名字
	lock := &Lock{GetKey_Lock(key)}

	log.Debugln("[DB] TryLock , key:", lock.Key)

	beginMS := TimeEpochMS()

	for {

		resp := WrapCommand(c, "SET", lock.Key, 1, "EX", LockKeyExpireSec, "NX")

		if resp.Err != nil {
			log.Errorln("[DB] TryLock failed, key:", lock.Key, resp.Err)
			return nil, ErrCode_DBFailed
		}

		var lockSuccess bool

		// 返回Nil表示Key已存在
		if !resp.IsType(redis.Nil) {

			if v, err := resp.Str(); err != nil {
				log.Errorln("[DB] TryLock failed, key:", lock.Key, err)
				return nil, ErrCode_DBFailed
			} else {
				lockSuccess = v == "OK"
			}
		}

		if lockSuccess {
			return lock, 0
		} else {

			if (TimeEpochMS() - beginMS) > int64(TryLockTimeOutMS) {

				log.Errorln("[DB] TryLock timeout, key:", lock.Key)
				return nil, ErrCode_DBFailed
			} else {
				time.Sleep(TryLockInterval) // TODO　改用After
			}
		}
	}

}

func Unlock(c *RedisClient, lock *Lock) int {

	if lock == nil {
		return 0
	}

	resp := WrapCommand(c, "DEL", lock.Key)

	if resp.Err != nil {
		log.Errorln("[DB] Unlock failed, key:", lock.Key, resp.Err)
		return ErrCode_DBFailed
	}

	return 0
}
