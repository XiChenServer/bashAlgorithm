package delayQueue

import (
	"bash_algorithm/delayQueue/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	redis2 "github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"time"
)

var (
	JobPoolKey         = "job_pool_key_"
	BaseDelayBucketKey = "base_delay_bucket"
	BaseReadyQueueKey  = "base_ready_queue"
	delayAddTime       = 1
)

type RedisJobData struct {
	Topic string
	ID    string
	Delay int
	TTR   int
	Body  *BodyContent
}
type BodyContent struct {
	OrderID   int
	OrderName string
}

// SetJobPool 组装数据
func (d *RedisJobData) SetJobPool(number int, ctx context.Context) bool {
	redisCoon := redis.GetRedisDb()
	for i := 0; i < number; i++ {
		d.Topic = "order_queue"
		d.ID = uuid.NewString()
		d.Delay = 1
		d.TTR = 3
		d.Body = &BodyContent{
			OrderID:   i,
			OrderName: "order_name_" + strconv.Itoa(i),
		}
		key := JobPoolKey + strconv.Itoa(i)
		delayKey := strconv.Itoa(i)
		data, _ := json.Marshal(d)
		//写入job pool
		_, err := redisCoon.Set(ctx, key, data, 0*time.Second).Result()
		if err != nil {
			fmt.Println("添加失败: ", err)
			return false
		}
		fmt.Println("添加成功: ", err)
		nowTime := time.Now().Unix()

		//delayTime := int(nowTime) + r.Intn(delayAddTime)
		delayTime := int(nowTime) + delayAddTime
		//写入delay queue
		redisCoon.ZAdd(ctx, BaseDelayBucketKey, redis2.Z{
			Score:  float64(delayTime),
			Member: delayKey,
		})
		//为了可以更好的演示，这每个过期时间增加几秒，防止一次性消费了
		if i%10 == 0 && i > 0 {
			delayAddTime += 10
		}
	}
	return true
}

// 定时timer.tick查询bucket中是否有过期的数据，如果有放入消费队列中
func TimerDelayBucket(redisCoon *redis2.Client, ctx context.Context, p *Pool, wg *sync.WaitGroup) error {
	nowTime := time.Now().Unix()
	result, err := redisCoon.ZRangeByScoreWithScores(ctx, BaseDelayBucketKey, &redis2.ZRangeBy{
		Min: "-inf",
		Max: strconv.FormatInt(nowTime, 10),
	}).Result()
	if err == nil {
		for _, z := range result {
			//进入ready queue
			redisCoon.LPush(ctx, BaseReadyQueueKey, z.Member)
			//写入通道说明有数据了，可以进行消费
			err := p.Put(&Task{
				Member: z.Member.(string),
				Wg:     wg,
			})
			if err != nil {
				fmt.Println(err)
			}
			wg.Add(1)
		}
	}
	return err

}

// 消费队列
func ConsumeQueue(redisCoon *redis2.Client, ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	//先判断list中是否有数据，有数据才有必要执行，没有数据直接返回就好了
	lenQueue, err := redisCoon.LLen(ctx, BaseReadyQueueKey).Result()
	if err != nil {
		return err
	}
	if lenQueue == 0 {
		return nil
	}
	result, err := redisCoon.LPop(ctx, BaseReadyQueueKey).Result()
	if err != nil {
		return err
	}
	fmt.Println("我消费了一个数据：", result)
	//这里可以实现需要的操作，这里简单实现了删除操作
	redisCoon.Del(ctx, JobPoolKey+result)
	redisCoon.ZRem(ctx, BaseDelayBucketKey, result)
	return nil
}
