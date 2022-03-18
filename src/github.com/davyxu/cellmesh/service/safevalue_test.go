package service

import (
	"github.com/davyxu/cellmesh/discovery"
	memsd "github.com/davyxu/cellmesh/discovery/memsd/api"
	"reflect"
	"testing"
)

func TestSafeGetValue(t *testing.T) {

	var origin []byte
	for i := 0; i < 12; i++ {
		//origin = append(origin, byte(rand.Int31n(127)))
		origin = append(origin, byte(i))
	}

	sdConfig := memsd.DefaultConfig()
	sdConfig.Address = GetDiscoveryAddr()
	discovery.Default = memsd.NewDiscovery(sdConfig)

	err := discovery.SafeSetValue(discovery.Default, "config/test", origin, true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var outData []byte
	err = discovery.SafeGetValue(discovery.Default, "config/test", &outData, true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(origin, outData) {
		t.FailNow()
	}
}
