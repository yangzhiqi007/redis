package memsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	"github.com/davyxu/cellmesh/discovery/memsd/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"time"
)

func (self *memDiscovery) clearCache() {

	self.svcCacheGuard.Lock()
	self.svcCache = map[string][]*discovery.ServiceDesc{}
	self.svcCacheGuard.Unlock()

	self.kvCacheGuard.Lock()
	self.kvCache = map[string][]byte{}
	self.kvCacheGuard.Unlock()
}

const PingTimesKey = "pingTimes"

//func (self *memDiscovery) recvPing(ses cellnet.Session) {
//
//	atomic.AddInt64(&self.pongTimes, 1)
//}

func (self *memDiscovery) sendPing(ses cellnet.Session) {

	if self.config.PingInterval == 0 {
		return
	}

	time.AfterFunc(self.config.PingInterval, func() {

		self.pingGuard.Lock()

		defer self.pingGuard.Unlock()

		var pingTimes int64
		if ses.(cellnet.ContextSet).FetchContext(PingTimesKey, &pingTimes) {
			if self.pongTimes != pingTimes {
				log.Warnf("mem sd ping timeout!, ping:%d pong: %d timeout: %v", pingTimes, self.pongTimes, self.config.PingCheckTimeout)
				return
			}
		}

		ses.(cellnet.ContextSet).SetContext(PingTimesKey, pingTimes+1)
		ses.Send(&proto.PingACK{})
	})
}

func (self *memDiscovery) Close() {
	if self.ses != nil {
		self.ses.Close()
	}
}

func (self *memDiscovery) TestPing(timeOut time.Duration) (ret error) {

	log.Debugf("memsd test ping start...")
	return remoteCallEx(self.Session(), &proto.PingACK{}, func(ack *proto.PingACK) {
		log.Debugf("memsd test ping done!")
	}, timeOut)
}

func (self *memDiscovery) connect(addr string) {
	p := peer.NewGenericPeer("tcp.Connector", "memsd", addr, model.Queue)

	proc.BindProcessorHandler(p, "memsd.cli", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *cellnet.SessionConnected:

			self.sesGuard.Lock()
			self.ses = ev.Session()
			self.sesGuard.Unlock()
			self.clearCache()
			ev.Session().Send(&proto.AuthREQ{
				Token:         self.token,
				ClientVersion: model.Version,
			})
		case *cellnet.SessionClosed:

			self.token = ""
			log.Errorf("memsd discovery lost! reason:%s", msg.Reason.String())

		case *proto.AuthACK:

			if msg.ServerVersion == "" {
				log.Errorf("memsd auth failed, client version: %s, unknown server version", model.Version)
				break
			}

			if msg.Code != 0 {
				log.Errorf("memsd auth failed: %s", msg.Code.String())
			} else {

				self.token = msg.Token

				if self.initWg != nil {
					// Pull的消息还要在queue里处理，这里确认处理完成后才算初始化完成
					self.initWg.Done()
				}

				log.Infof("memsd discovery ready!")

				self.triggerNotify("ready", 0, "")

				self.sendPing(ev.Session())
			}

		case *proto.PingACK:

			//self.recvPing(ev.Session())

			log.Debugf("recv pingack from memsd")

			self.sendPing(ev.Session())

		case *proto.ValueChangeNotifyACK:

			if model.IsServiceKey(msg.Key) {
				self.updateSvcCache(msg.SvcName, msg.Value)
			} else {
				self.updateKVCache(msg.Key, msg.Value)
			}

		case *proto.ValueDeleteNotifyACK:

			if model.IsServiceKey(msg.Key) {
				svcid := model.GetSvcIDByServiceKey(msg.Key)
				self.deleteSvcCache(svcid, msg.SvcName)
			} else {
				self.deleteKVCache(msg.Key)
			}
		}
	})

	// noDelay
	p.(cellnet.TCPSocketOption).SetSocketBuffer(1024*1024, 1024*1024, true)

	// 断线后自动重连
	p.(cellnet.TCPConnector).SetReconnectDuration(time.Second * 5)

	p.Start()

	for {

		if p.(cellnet.PeerReadyChecker).IsReady() {
			break
		}

		time.Sleep(time.Millisecond * 500)
	}

}
