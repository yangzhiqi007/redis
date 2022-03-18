package redisdrv

import (
	"errors"
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
	"reflect"
)

type RespCoder interface {
	Unmarshal(resp *redis.Resp)
}

type valueOject struct {
	Value    interface{}
	UserData interface{}
}

// 除ObjectList模型外的批量加载器
type ValueLoader struct {
	c *redis.Client

	values []valueOject
}

func (self *ValueLoader) AddObject(cmd string, dataPtr interface{}, args ...interface{}) {

	self.AddObjectEx(cmd, dataPtr, nil, args)
}

func (self *ValueLoader) AddObjectEx(cmd string, dataPtr, userData interface{}, args ...interface{}) {

	if log.IsDebugEnabled() {
		log.Debugf("[DB] PipeAppend, %s %s    --->%s", cmd, argsToString(args), getStackFileString(2))
	}

	self.values = append(self.values, valueOject{
		Value:    dataPtr,
		UserData: userData,
	})

	WrapPipeAppend(self.c, cmd, args...)
}

func (self *ValueLoader) Load(debugEnv string) {

	if log.IsDebugEnabled() {
		log.Debugf("[DB] LoadValue | %s", debugEnv)
	}

	for _, v := range self.values {
		resp := self.c.PipeResp()
		RespToAny(resp, v.Value)
	}

}

func (self *ValueLoader) GetObject(index int) interface{} {

	return self.values[index].Value
}

func (self *ValueLoader) GetObjectUserData(index int) interface{} {

	return self.values[index].UserData
}

func (self *ValueLoader) ObjectCount() int {
	return len(self.values)
}

func NewValueLoader(c *RedisClient) *ValueLoader {

	c.PipeClear()

	return &ValueLoader{
		c: c,
	}
}

func RespToAny(resp *redis.Resp, dataPtr interface{}) {

	var logErr error

	if resp.Err != nil {
		logErr = errors.New(fmt.Sprintf("redis resp err:%s", resp.Err.Error()))
		goto OnError

	}

	if resp.IsType(redis.Nil) {
		return
	}

	//log.Debugf("%+T", dataPtr)

	switch vp := dataPtr.(type) {

	case RespCoder:
		vp.Unmarshal(resp)
	case *int:
		if v, err := resp.Int64(); err == nil {
			*vp = int(v)
		} else {
			logErr = err
			goto OnError
		}
	case *int32:
		if v, err := resp.Int64(); err == nil {
			*vp = int32(v)
		} else {
			logErr = err
			goto OnError
		}
	case *int64:
		if v, err := resp.Int64(); err == nil {
			*vp = v
		} else {
			logErr = err
			goto OnError
		}
	case *float64:
		if v, err := resp.Float64(); err == nil {
			*vp = v
		} else {
			logErr = err
			goto OnError
		}
	case *[]byte:
		if v, err := resp.Bytes(); err == nil {
			*vp = v
		} else {
			logErr = err
			goto OnError
		}
	case *string:
		if v, err := resp.Str(); err == nil {
			*vp = v
		} else {
			logErr = err
			goto OnError
		}
	case *[]string:
		if v, err := resp.List(); err == nil {
			*vp = v
		} else {
			logErr = err
			goto OnError
		}
	case interface {
		UnmarshalMsg(bts []byte) (o []byte, err error)
	}:

		data, err := resp.Bytes()

		if err != nil {
			logErr = err
			goto OnError
		}

		_, err = vp.UnmarshalMsg(data)

		if err != nil {
			logErr = err
			goto OnError
		}

	default:
		panic("[DB] unsupport redis value to get: " + reflect.TypeOf(vp).Elem().Name())
	}

	return

OnError:
	log.Errorf("[DB] RespToAny failed:	%s --->%s", logErr, getStackFileString(2))
	panic(ErrCode_DBFailed)

}
