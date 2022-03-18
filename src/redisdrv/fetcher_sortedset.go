package redisdrv

import (
	"errors"
	"github.com/mediocregopher/radix.v2/redis"
	"reflect"
)

// 有序集合成员
type SortedSetMember struct {
	Member int64   // 成员,可扩展类型
	Score  float64 // 分数,固定类型
}

type SortedSetFetcher struct {
	Data   []SortedSetMember
}

func (self *SortedSetFetcher) Unmarshal(resp *redis.Resp) {
	sortedSetUnmarshal(resp, &self.Data)
}

func sortedSetUnmarshal(resp *redis.Resp, dataPtr interface{}) {

	var (
		err          error
		resps        []*redis.Resp
		vdata, slice reflect.Value
		sIndex       int
	)

	if resps, err = resp.Array(); err != nil {
		goto OnError
	}

	if len(resps)%2 != 0 {
		err = errors.New("list has odd number of elements")
		goto OnError
	}

	vdata = reflect.Indirect(reflect.ValueOf(dataPtr))

	slice = reflect.MakeSlice(vdata.Type(), len(resps)/2, len(resps)/2)

	for i := 0; i < len(resps); i += 2 {

		sliceValue := reflect.New(slice.Type().Elem())

		RespToAny(resps[i], sliceValue.Elem().FieldByName("Member").Addr().Interface())
		RespToAny(resps[i+1], sliceValue.Elem().FieldByName("Score").Addr().Interface())

		slice.Index(sIndex).Set(sliceValue.Elem())

		sIndex++
	}

	vdata.Set(slice)

	return

OnError:

	log.Errorln("[DB] sortedSetUnmarshal failed:", err)
	panic(ErrCode_DBFailed)

}
