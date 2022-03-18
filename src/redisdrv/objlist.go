package redisdrv

import (
	"github.com/mediocregopher/radix.v2/redis"
	"sync/atomic"
)

type ObjectKeyFetcher interface {

	// 获取Redis中主Key（最外面的一层）
	GetMainKey(env interface{}) interface{}

	// 获取HASH类型值的key, hash容器的key，一般为玩家身上的
	GetHashKey(env interface{}) interface{}
}

// 对hash数据结构且value为msgpack类型的封装，
type ObjectList struct {
	c *redis.Client

	pairs []*keyPair

	id int64

	mapper map[interface{}]*keyPair
}

type msgPackObject interface {
	MarshalMsg(b []byte) (o []byte, err error)
	UnmarshalMsg(bts []byte) (o []byte, err error)
}

// 查询结构中的Env，只会保留一个，不会重复查询
func (self *ObjectList) MappingObjectByEnv() {

	self.mapper = make(map[interface{}]*keyPair)
}

// fetcher中提供的对象在Load之后会发生Save行为
func (self *ObjectList) AddSaveableObject(fetcher ObjectKeyFetcher, env interface{}) {

	self.addObjectEx(fetcher, env, nil, true)
}

// fetcher对象只加载不发生Save
func (self *ObjectList) AddObject(fetcher ObjectKeyFetcher, env interface{}) {
	self.addObjectEx(fetcher, env, nil, false)
}

func (self *ObjectList) AddObjectEx(fetcher ObjectKeyFetcher, env, userData interface{}) {
	self.addObjectEx(fetcher, env, userData, false)
}

func (self *ObjectList) AddSaveableObjectEx(fetcher ObjectKeyFetcher, env, userData interface{}) {
	self.addObjectEx(fetcher, env, userData, true)
}

func (self *ObjectList) addObjectEx(fetcher ObjectKeyFetcher, env, userData interface{}, save bool) {

	pair := &keyPair{
		Fetcher:  fetcher,
		Env:      env,
		UserData: userData,
		Save:     save,
	}

	if self.mapper != nil && env != nil {

		// key重复，无需缓冲
		if _, ok := self.mapper[env]; ok {
			return
		}

		self.mapper[env] = pair

	}

	self.pairs = append(self.pairs, pair)

	return
}

func (self *ObjectList) FindObjectByEnv(env interface{}) interface{} {

	if self.mapper == nil {
		return nil
	}

	if v, ok := self.mapper[env]; ok {
		return v.Fetcher
	}

	return nil
}

func (self *ObjectList) MarkSaveableByEnv(env interface{}) {

	for _, pair := range self.pairs {
		if pair.Env == env {
			pair.Save = true
			break
		}
	}
}

func (self *ObjectList) GetObject(index int) interface{} {

	return self.pairs[index].Fetcher
}

func (self *ObjectList) GetObjectUserData(index int) interface{} {

	return self.pairs[index].UserData
}

func (self *ObjectList) GetEnv(index int) interface{} {

	return self.pairs[index].Env
}

func (self *ObjectList) ObjectCount() int {
	return len(self.pairs)
}

func (self *ObjectList) serialize() {
	for _, obj := range self.pairs {

		if !obj.Save {
			continue
		}

		msgpObj := obj.Fetcher.(msgPackObject)

		data, err := msgpObj.MarshalMsg(nil)
		if err != nil {
			log.Errorln(err)
			panic(ErrCode_DBFailed)
		}

		data, err = encodeData(obj.Fetcher, data)
		if err != nil {
			log.Errorln(err)
			panic(ErrCode_DBFailed)
		}

		obj.HSet(self, data)
	}

}

func (self *ObjectList) deserialize() error {
	for _, obj := range self.pairs {

		resp := self.c.PipeResp()

		if resp.IsType(redis.Nil) {
			continue
		}

		if resp.Err != nil {
			return resp.Err
		}

		data, err := resp.Bytes()

		if err != nil {
			return err
		}

		msgpObj := obj.Fetcher.(msgPackObject)

		data, err = decodeData(data)
		if err != nil {
			return err
		}

		_, err = msgpObj.UnmarshalMsg(data)

		if err != nil {
			return err
		}
	}

	return nil

}

func (self *ObjectList) loadCmd() {
	for _, p := range self.pairs {
		p.HGet(self)
	}
}

func NewObjectList(c *RedisClient) *ObjectList {
	return &ObjectList{
		c:  c,
		id: getSeq(),
	}
}

var (
	globalSeq int64
)

func getSeq() int64 {

	return atomic.AddInt64(&globalSeq, 1)
}
