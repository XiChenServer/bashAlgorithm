package test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"
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
