package Map

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_string(t *testing.T) {
	str := "天行健"
	for i, v := range str {
		fmt.Println(i, string(v), v)
	}
}

func Test_arr(t *testing.T) {
	var arr [10]int
	arr[4] = 1
	arr[5] = 24
	arr[4] = 3
	fmt.Println(arr)
}
func Test_arrAppend(t *testing.T) {
	arr := []int{}

	arr = append(arr, 10)
	for i := 0; i < 1000; i++ {
		arr = append(arr, []int{1, 2, 3, 4}...)
		if cap(arr) == len(arr) {
		}
	}
	fmt.Println(arr, cap(arr))
}

func Test_defer(t *testing.T) {
	for i := 0; i < 5; i++ {
		defer func(i int) {
			fmt.Println(i, &i)
		}(i)
	}
}

var wg sync.WaitGroup

func product(ch chan int) {
	defer wg.Done() // 确保在退出前调用 Done

	num := 0
	for {
		time.Sleep(time.Second)
		ch <- num
		num++
	}
}

func consumer(ch <-chan int) {
	defer wg.Done() // 确保在退出前调用 Done

	for num := range ch { // 用 range 简化接收数据
		fmt.Println(num)
	}
}

func Test_mq(t *testing.T) {
	ch := make(chan int, 10)

	wg.Add(2)
	go product(ch)
	go consumer(ch)

	time.Sleep(5 * time.Second) // 停顿一段时间，模拟生产和消费过程
	//close(ch)                   // 关闭通道，通知 consumer 退出
	wg.Wait()
}
