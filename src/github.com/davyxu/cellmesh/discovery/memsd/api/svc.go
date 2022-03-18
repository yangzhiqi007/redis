package memsd

import (
	"encoding/json"
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	"github.com/davyxu/cellmesh/discovery/memsd/proto"
	"github.com/davyxu/cellnet/util"
	"time"
)

func (self *memDiscovery) Register(svc *discovery.ServiceDesc) (retErr error) {

	if svc.Name == "" {
		return errors.New("expect svc name")
	}

	if svc.ID == "" {
		return errors.New("expect svc id")
	}

	data, err := json.Marshal(svc)
	if err != nil {
		return err
	}

	callErr := self.remoteCall(&proto.SetValueREQ{
		Key:     model.ServiceKeyPrefix + svc.ID,
		Value:   data,
		SvcName: svc.Name,
	}, func(ack *proto.SetValueACK) {
		retErr = codeToError(ack.Code)
	})

	if retErr != nil {
		return
	}

	return callErr
}

func (self *memDiscovery) Deregister(svcid string) error {

	return self.DeleteValue(model.ServiceKeyPrefix + svcid)
}

func (self *memDiscovery) Query(name string) (ret []*discovery.ServiceDesc) {

	self.svcCacheGuard.RLock()
	defer self.svcCacheGuard.RUnlock()

	return self.svcCache[name]
}

func (self *memDiscovery) QueryAll() (ret []*discovery.ServiceDesc) {

	self.svcCacheGuard.RLock()
	defer self.svcCacheGuard.RUnlock()

	for _, list := range self.svcCache {
		ret = append(ret, list...)
	}

	return
}

func (self *memDiscovery) ClearService() {
	self.remoteCall(&proto.ClearSvcREQ{}, func(ack *proto.ClearSvcACK) {})
}

func (self *memDiscovery) triggerNotify(mode string, timeout time.Duration, desc string) {

	self.notifyMap.Range(func(key, value interface{}) bool {

		if value == nil {
			return true
		}

		ctx := value.(*notifyContext)

		if ctx.mode != mode {
			return true
		}

		c := key.(chan struct{})

		if timeout == 0 {

			select {
			case c <- struct{}{}:
			default:
			}

		} else {
			select {
			case c <- struct{}{}:
			case <-time.After(timeout):
				// 接收通知阻塞太久，或者没有释放侦听的channel
				log.Errorf("notify(%s) timeout, desc: '%s' regstack: %s ", ctx.mode, desc, ctx.stack)
			}
		}

		return true
	})

}

func (self *memDiscovery) RegisterNotify(mode string) (ret chan struct{}) {
	ret = make(chan struct{}, 10)

	switch mode {
	case "add", "ready":
		self.notifyMap.Store(ret, &notifyContext{
			mode:  mode,
			stack: util.StackToString(5),
		})
	default:
		panic("unknown notify mode: " + mode)
	}

	return
}

func (self *memDiscovery) DeregisterNotify(mode string, c chan struct{}) {

	switch mode {
	case "add", "ready":
		self.notifyMap.Store(c, nil)
	default:
		panic("unknown notify mode: " + mode)
	}
}

func (self *memDiscovery) updateSvcCache(svcName string, value []byte) {
	self.svcCacheGuard.Lock()

	list := self.svcCache[svcName]

	var desc discovery.ServiceDesc
	err := json.Unmarshal(value, &desc)
	if err != nil {
		log.Errorf("ServiceDesc unmarshal failed, %s", err)
		self.svcCacheGuard.Unlock()
		return
	}

	var notfound = true
	for index, svc := range list {
		if svc.ID == desc.ID {
			list[index] = &desc
			notfound = false
			break
		}
	}

	if notfound {
		list = append(list, &desc)
	}

	self.svcCache[svcName] = list
	self.svcCacheGuard.Unlock()

	self.triggerNotify("add", time.Second*10, desc.String())
}

func (self *memDiscovery) deleteSvcCache(svcid, svcName string) {

	self.svcCacheGuard.Lock()
	defer self.svcCacheGuard.Unlock()

	list := self.svcCache[svcName]

	for index, svc := range list {
		if svc.ID == svcid {
			list = append(list[:index], list[index+1:]...)
			break
		}
	}

	self.svcCache[svcName] = list
}
