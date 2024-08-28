package test

import (
	"fmt"
	"sort"
	"testing"
)

func main() {
	var n int
	var s string
	fmt.Scan(&n)
	fmt.Scan(&s)

	// 统计每个长度的密码数量
	lengths := make(map[int]int)
	for i := 0; i < n; i++ {
		var password string
		fmt.Scan(&password)
		lengths[len(password)]++
	}

	// 获取正确密码的长度
	passwordLength := len(s)
	minTries, maxTries := 0, 0
	cumulativeCount := 0

	// 将所有的长度进行排序
	lenArr := make([]int, 0, len(lengths))
	for k := range lengths {
		lenArr = append(lenArr, k)
	}
	sort.Ints(lenArr)

	for _, l := range lenArr {
		if l < passwordLength {
			// 累积所有比正确密码短的密码的数量
			minTries += lengths[l]
			maxTries += lengths[l]
			cumulativeCount += lengths[l]
		} else if l == passwordLength {
			// 正确密码长度的情况
			minTries += 1
			maxTries += lengths[l]
			break
		}
	}

	// 输出最少和最多的尝试次数
	fmt.Println(minTries+cumulativeCount, maxTries+cumulativeCount)
}
func Test_c(t *testing.T) {
	main()
}
