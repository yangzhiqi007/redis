package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"time"
)

// 服务互联消息处理
type SvcEventHooker struct {
}

func sendPing(duration time.Duration, ses cellnet.Session) {
	if duration == 0 {
		return
	}

	time.AfterFunc(duration, func() {
		if ses != nil {
			ses.Send(&ServicePingREQ{
				Time:  time.Now().Unix(),
				SvcID: GetLocalSvcID(),
			})
		}
	})
}

const PingDuration = 0 //time.Second * 60 * 3

func (SvcEventHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch msg := inputEvent.Message().(type) {
	case *ServiceIdentifyACK:

		if pre := GetRemoteService(msg.SvcID); pre == nil {

			// 添加连接上来的对方服务
			AddRemoteService(inputEvent.Session(), msg.SvcID, msg.SvcName)
		}
	case *ServicePingREQ: // Acceptor处理
		if svc := GetRemoteService(msg.SvcID); svc != nil {
			inputEvent.Session().Send(&ServicePingACK{
				Time: msg.Time,
			})
		}
	case *ServicePingACK:

		sendPing(PingDuration, inputEvent.Session())

	case *cellnet.SessionConnected:

		ctx := inputEvent.Session().Peer().(cellnet.ContextSet)

		var sd *discovery.ServiceDesc
		if ctx.FetchContext("sd", &sd) {

			// 用Connector的名称（一般是ProcName）让远程知道自己是什么服务，用于网关等需要反向发送消息的标识
			inputEvent.Session().Send(&ServiceIdentifyACK{
				SvcName: GetProcName(),
				SvcID:   GetLocalSvcID(),
			})

			AddRemoteService(inputEvent.Session(), sd.ID, sd.Name)

			sendPing(PingDuration, inputEvent.Session())

		} else {

			log.Errorf("Make sure call multi.AddPeer before peer.Start, peer: %s", inputEvent.Session().Peer().TypeName())
		}

	case *cellnet.SessionClosed:

		RemoveRemoteService(inputEvent.Session())
	}

	return inputEvent

}

func (SvcEventHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func init() {

	// 服务器间通讯协议
	proc.RegisterProcessor("tcp.svc", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(SvcEventHooker), new(tcp.MsgHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})

	// 与客户端通信的处理器
	proc.RegisterProcessor("tcp.client", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(tcp.MsgHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
