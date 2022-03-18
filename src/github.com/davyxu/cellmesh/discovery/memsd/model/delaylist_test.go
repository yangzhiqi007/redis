package model

import (
	"container/list"
	"testing"
)

func TestDelay(t *testing.T) {

	var l = list.New()
	for i := 0; i < 10; i++ {
		l.PushBack(i)
	}

	consume(t, "a", l)
	consume(t, "b", l)
	consume(t, "c", l)
	consume(t, "d", l)
	consume(t, "e", l)

}

func consume(t *testing.T, s string, l *list.List) {
	const n = 3
	for i := 0; i <= n; i++ {

		e := l.Front()
		if e == nil {
			break
		}
		t.Logf("%s %d", s, e.Value)
		l.Remove(e)
	}

}
