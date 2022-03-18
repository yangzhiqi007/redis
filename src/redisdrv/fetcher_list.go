package redisdrv

import (
	"github.com/mediocregopher/radix.v2/redis"
	"reflect"
)

type ListFetcher struct {
	Data []int64
}

func (self *ListFetcher) Unmarshal(resp *redis.Resp) {
	listUnmarshal(resp, &self.Data)
}

type ListFetcherStr struct {
	Data []string
}

func (self *ListFetcherStr) Unmarshal(resp *redis.Resp) {
	listUnmarshal(resp, &self.Data)
}

func listUnmarshal(resp *redis.Resp, dataPtr interface{}) {

	var (
		err          error
		resps        []*redis.Resp
		vdata, slice reflect.Value
	)

	if resps, err = resp.Array(); err != nil {
		goto OnError
	}

	vdata = reflect.Indirect(reflect.ValueOf(dataPtr))

	slice = reflect.MakeSlice(vdata.Type(), len(resps), len(resps))

	for i, resp := range resps {

		sliceValue := reflect.New(slice.Type().Elem())

		RespToAny(resp, sliceValue.Interface())

		slice.Index(i).Set(sliceValue.Elem())

	}

	vdata.Set(slice)

	return

OnError:

	log.Errorln("[DB] listUnmarshal failed:", err)
	panic(ErrCode_DBFailed)

}

type ListFetcherBytes struct {
	Data [][]byte
}

func (self *ListFetcherBytes) Unmarshal(resp *redis.Resp) {

	resps, err := resp.Array()
	if err != nil {
		log.Errorln("[DB] ListFetcherBytes failed:", err)
		panic(ErrCode_DBFailed)
	}

	self.Data = make([][]byte, len(resps))

	for i, resp := range resps {

		RespToAny(resp, &self.Data[i])
	}
}
