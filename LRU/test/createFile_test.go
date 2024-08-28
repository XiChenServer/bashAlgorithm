package test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_CreateFile(t *testing.T) {
	file, err := os.OpenFile("/home/zwm/go_projects/bash_algorithm/LRU/test/createFile.txt",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		t.Fatalf("无法打开或创建文件: %v", err)
	}
	defer file.Close()

	for i := 0; i < 10000; i++ {
		str, err := generateRandomLetters()
		if err != nil {
			t.Errorf("生成随机字母失败: %v", err)
			continue // 可以选择跳过当前循环，或者根据需要处理错误
		}
		if _, err := file.WriteString(strconv.Itoa(i) + " " + str + "\n"); err != nil {
			t.Errorf("写入文件失败: %v", err)
			return // 可以选择返回，或者根据需要处理错误
		}
	}
}

// 生成随机的6位字母字符串
func generateRandomLetters() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 6)
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[n.Int64()]
	}
	return string(b), nil
}

func Test_1(t *testing.T) {
	mp := make(map[int]int)
	for i := 0; i < 15; i++ {
		mp[i] = i
	}
	for k, v := range mp {
		fmt.Println(k, v)
	}
	fmt.Println("111111111111111111")
	for k, v := range mp {
		fmt.Println(k, v)
	}
	fmt.Println("111111111111111111")
	for k, v := range mp {
		fmt.Println(k, v)
	}

}

func Test_Mid(t *testing.T) {
	nums := []int{}
	target := 0
	l, r := 0, len(nums)-1
	for r < l {
		mid := (l + r) / 2
		if target < nums[mid] {
			r = mid - 1
		}
		if target > nums[mid] {
			l = mid
		}
		if target == nums[mid] {
			break
		}

	}

}

var (
	a     = 0
	name  = "fjw"
	name1 = []int{1, 2, 3}
)

func Print(i int) {
	for true {
		if a%3 == i {
			time.Sleep(time.Second * 2)
			fmt.Println(byte(name[a%3] - 'a'))
			a++
		}
	}

}

func Test_dayin(t *testing.T) {

	go Print(0)
	go Print(1)
	go Print(2)
	fmt.Println(name)
	for true {

	}

}

func ChanPrint(ch <-chan int) {
	for {
		select {
		case a := <-ch:
			fmt.Println(a % 3)
		default:

		}
	}

}
func Test_chan(t *testing.T) {
	ch := make(chan int, 1)
	a := 0
	go ChanPrint(ch)
	for {
		time.Sleep(time.Second)
		ch <- a
		a++
	}
}

var (
	num int
	mu  sync.Mutex
	wg  sync.WaitGroup
)

func Print1(i int) {
	defer wg.Done()
	for {
		mu.Lock()
		if num%3 == i {
			fmt.Println(i + 1)
			num++
			mu.Unlock()
			time.Sleep(100 * time.Millisecond) // 加入短暂的休眠，以避免过多的 CPU 占用
		} else {
			mu.Unlock()
			time.Sleep(10 * time.Millisecond) // 加入短暂的休眠，以避免忙等
		}
	}
}

func Test_Mutest(t *testing.T) {
	wg.Add(3)
	go Print1(0)
	go Print1(1)
	go Print1(2)
	wg.Wait()
}

type mp struct {
	m map[int]int
	sync.RWMutex
}

func (m *mp) Get(key int) int {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	return m.m[key]
}
func (m *mp) Set(key int, value int) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.m[key] = value
}

func (m *mp) Del(key int) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	delete(m.m, key)
}

func suanfa(arr []int, x, y int) int {
	sort.Ints(arr)
	length := len(arr)
	for i := range arr {
		if i+1 >= x && i+1 <= y && length-i-1 >= x && length-i-1 <= y {
			return i + 1
		}
	}
	return 0
}

func Test_mei(t *testing.T) {

}
