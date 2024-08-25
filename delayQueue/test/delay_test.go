package test

import (
	"bash_algorithm/delayQueue"
	"bash_algorithm/delayQueue/redis"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestConsuming(t *testing.T) {
	var data delayQueue.RedisJobData
	redisConn := redis.GetRedisDb()
	conn := context.Background()
	data.SetJobPool(10, conn)
	pool, err := delayQueue.NewPool(20)
	if err != nil {
		panic(err)
	}
	wg := new(sync.WaitGroup)

	c := time.Tick(1 * time.Second)

	for next := range c {
		fmt.Println("我在执行了")
		err := delayQueue.TimerDelayBucket(redisConn, conn, pool, wg)
		if err != nil {
			fmt.Println("定时timer发生错误：", next, err)
		}
	}
	wg.Wait()

	pool.Close()
}
