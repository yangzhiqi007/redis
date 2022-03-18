package redisdrv

import "testing"

func foo(t *testing.T, args ...interface{}) {

	var newValue = []interface{}{"b"}
	for i := 1; i < len(args); i++ {
		newValue = append(newValue, args[i])
	}

	t.Log(args...)
	t.Log(newValue)
}

func TestArray(t *testing.T) {

	foo(t, 1, 2, 3)
}
