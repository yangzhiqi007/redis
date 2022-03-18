package redisdrv

import (
	"github.com/mediocregopher/radix.v2/redis"
	"reflect"
)

// 地理集合成员
type GeoSetMemberInt64 struct {
	Member   int64
	Distance float64
}

type GeoSetFetcher struct {
	Data []GeoSetMemberInt64
}

func (self *GeoSetFetcher) Unmarshal(resp *redis.Resp) {
	geoSetUnmarshal(resp, &self.Data)
}

func geoSetUnmarshal(resp *redis.Resp, dataPtr interface{}) {

	var (
		err               error
		respsOut, respsIn []*redis.Resp
		vdata             reflect.Value
		slice             reflect.Value
		sIndex            int
	)

	if respsOut, err = resp.Array(); err != nil {
		goto OnError
	}

	vdata = reflect.Indirect(reflect.ValueOf(dataPtr))

	slice = reflect.MakeSlice(vdata.Type(), len(respsOut), len(respsOut))

	for _, v := range respsOut {

		if respsIn, err = v.Array(); err != nil {
			goto OnError
		}

		sliceValue := reflect.New(slice.Type().Elem())

		RespToAny(respsIn[0], sliceValue.Elem().FieldByName("Member").Addr().Interface())
		RespToAny(respsIn[1], sliceValue.Elem().FieldByName("Distance").Addr().Interface())

		slice.Index(sIndex).Set(sliceValue.Elem())

		sIndex++
	}

	vdata.Set(slice)

	return

OnError:

	log.Errorln("[DB] geoSetUnmarshal failed:", err)
	panic(ErrCode_DBFailed)

}
