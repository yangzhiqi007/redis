package redisdrv

type keyPair struct {
	Fetcher  ObjectKeyFetcher
	Env      interface{}
	UserData interface{}
	Save     bool
}

func (self *keyPair) HGet(ser *ObjectList) {

	mainKey := self.Fetcher.GetMainKey(self.Env)
	hashKey := self.Fetcher.GetHashKey(self.Env)

	if log.IsDebugEnabled() {
		log.Debugf("[DB] %d HGET mainKey:%v hashKey:%v", ser.id, mainKey, hashKey)
	}

	WrapPipeAppend(ser.c, "HGET", mainKey, hashKey)
}

func (self *keyPair) HSet(ser *ObjectList, data []byte) {

	mainKey := self.Fetcher.GetMainKey(self.Env)
	hashKey := self.Fetcher.GetHashKey(self.Env)

	if log.IsDebugEnabled() {
		log.Debugf("[DB] %d HSET mainKey:%v hashKey:%v", ser.id, mainKey, hashKey)
	}

	WrapPipeAppend(ser.c, "HSET", mainKey, hashKey, data)
}
