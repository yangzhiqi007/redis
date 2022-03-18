package redis

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"time"
)

var (
	// RedisClient redis 客户端对象. 在 InitRedis() 中赋值
	RedisClient *redis.Client

	// RedisPool redis 连接池对象. 在 InitRedis() 中赋值
	RedisPool *pool.Pool
)

func RedisConnect(addr, password string, maxConnCount int, dbIndex int) {

	var faildWaitSec = 2

	for {
		pool, err := pool.NewCustom("tcp", addr, maxConnCount, func(network, addr string) (*redis.Client, error) {

			client, err := redis.DialTimeout(network, addr, time.Second*5)
			if err != nil {
				fmt.Errorf("redis.Dial %s", err.Error())
				return nil, err
			}

			if len(password) > 0 {
				if err = client.Cmd("AUTH", password).Err; err != nil {
					fmt.Errorf("redis.Auth %s %s", password, err.Error())
					client.Close()
					return nil, err
				}
			}

			if err = client.Cmd("SELECT", dbIndex).Err; err != nil {
				fmt.Errorf("redis.SELECT %d %s", dbIndex, err.Error())
				client.Close()
				return nil, err
			}

			fmt.Printf("Create redis pool connection: %s index: %d", addr, dbIndex)

			return client, nil
		})

		if err != nil {
			fmt.Errorf("Redis connect failed: %s, wait %d secods retry...", err, faildWaitSec)

			time.Sleep(time.Duration(faildWaitSec) * time.Second)

			if faildWaitSec < 10 {
				faildWaitSec += 2
			}

			continue
		}

		RedisPool = pool

		break
	}

}
