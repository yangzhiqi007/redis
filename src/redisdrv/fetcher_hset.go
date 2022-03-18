package redisdrv

import (
	"errors"
	"github.com/mediocregopher/radix.v2/redis"
	"reflect"
)

// Hash成员
type HScanMember struct {
	Field string // key,固定类型
	Value int64  // value,固定类型
}

// Hash字段和值遍历器
type HScanFetcher struct {
	Cursor int32         // 迭代一次后返回的游标
	Data   []HScanMember // 迭代一次后返回的数据
}

func (self *HScanFetcher) Unmarshal(resp *redis.Resp) {
	hScanUnmarshal(resp, &self.Cursor, &self.Data)
}

func hScanUnmarshal(resp *redis.Resp, cursorPtr interface{}, dataPtr interface{}) {

	var (
		err          error
		resps        []*redis.Resp
		subResps     []*redis.Resp
		vdata, slice reflect.Value
		sIndex       int
	)

	if resps, err = resp.Array(); err != nil {
		goto OnError
	}

	if len(resps) != 2 {
		err = errors.New("hash resps is not 2")
		goto OnError
	}

	RespToAny(resps[0], cursorPtr)

	if subResps, err = resps[1].Array(); err != nil {
		goto OnError
	}

	if len(subResps)%2 != 0 {
		err = errors.New("hash sub resps has odd number of elements")
		goto OnError
	}

	vdata = reflect.Indirect(reflect.ValueOf(dataPtr))

	slice = reflect.MakeSlice(vdata.Type(), len(subResps)/2, len(subResps)/2)

	for i := 0; i < len(subResps); i += 2 {

		sliceValue := reflect.New(slice.Type().Elem())

		RespToAny(subResps[i], sliceValue.Elem().FieldByName("Field").Addr().Interface())
		RespToAny(subResps[i+1], sliceValue.Elem().FieldByName("Value").Addr().Interface())

		slice.Index(sIndex).Set(sliceValue.Elem())

		sIndex++
	}

	vdata.Set(slice)

	return

OnError:

	log.Errorln("[DB] hScanUnmarshal failed:", err)
	panic(ErrCode_DBFailed)

}
