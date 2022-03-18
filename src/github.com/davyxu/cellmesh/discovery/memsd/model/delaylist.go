package model

import (
	"container/list"
	"github.com/davyxu/cellnet/timer"
	"github.com/davyxu/golog"
	"time"
)

var (
	delayList = list.New()
)

type Payload struct {
	Value   int
	Handler func()
}

func DelayProc(p Payload) {
	delayList.PushBack(p)
}

var log = golog.New("memsdm")

func ProcTask() {

	timer.NewLoop(Queue, time.Second, func(loop *timer.Loop) {

		var beginCount = -1

		var totalPay int
		for totalPay < 100 {

			e := delayList.Front()
			if e == nil {
				break
			}

			if beginCount == -1 {

				beginCount = delayList.Len()
			}

			p := e.Value.(Payload)
			p.Handler()
			totalPay += p.Value

			delayList.Remove(e)

		}

		if beginCount != -1 {
			log.Debugf("Delay proc : %d -> %d  sesCount: %d  valueCount: %d", beginCount, delayList.Len(), SessionCount(), ValueCount())
		}

	}, nil).Start()

}
