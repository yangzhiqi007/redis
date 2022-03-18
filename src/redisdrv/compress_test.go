package redisdrv

import (
	"reflect"
	"testing"
)

type dummy int

func (dummy) EnableObjectCompress() {

}

func TestCompress(t *testing.T) {

	testData := []byte("hello")

	data, err := encodeData(new(dummy), testData)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	ret, err := decodeData(data)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if !reflect.DeepEqual(ret, testData) {
		t.FailNow()
	}
}
