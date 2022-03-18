package main

import (
	"log"
	"model"
	"redis"
	"redisdrv"
)

func main() {
	log.Println("hello world")

	redis.RedisConnect("100.64.25.10:16379", "", 10, 1)
	redisdrv.OperateDB(redis.RedisPool, func(client *redisdrv.RedisClient) {
		loader := redisdrv.NewObjectList(client)
		var (
			A model.A
		)
		loader.AddSaveableObject(&A, 2)
		redisdrv.LoadObject("A", loader)

		A.Name = "yangzhiqi2"
		A.Age = 26

		redisdrv.SaveObject("A", loader)
	})
}
