package bitMap

import (
	"fmt"
	"sync"
	"testing"
)

func TestAdd(t *testing.T) {
	bitmap := NewBitMap()

	// 测试添加不同的数值
	bitmap.Add(0) // 0 不应该被添加
	bitmap.Add(7)
	bitmap.Add(8)
	bitmap.Add(15)
	bitmap.Add(16)

	// 验证添加后的位图内容
	if !bitmap.Exist(7) {
		t.Errorf("Expected 7 to be in the bitmap")
	}
	if !bitmap.Exist(8) {
		t.Errorf("Expected 8 to be in the bitmap")
	}
	if !bitmap.Exist(15) {
		t.Errorf("Expected 15 to be in the bitmap")
	}
	if !bitmap.Exist(16) {
		t.Errorf("Expected 16 to be in the bitmap")
	}

	// 测试位图长度
	expectedLen := 2 // For indices 0-15 and 16
	if bitmap.Len() != expectedLen {
		t.Errorf("Expected length %d but got %d", expectedLen, bitmap.Len())
	}

	// 测试 ToString
	expectedStr := "0000001100000011"
	if bitmap.ToString() != expectedStr {
		t.Errorf("Expected binary string '%s' but got '%s'", expectedStr, bitmap.ToString())
	}
}

func TestDel(t *testing.T) {
	bitmap := NewBitMap()

	// 添加一些数值
	bitmap.Add(0)
	bitmap.Add(8)
	bitmap.Add(15)
	bitmap.Add(16)
	fmt.Println([]byte(bitmap.bitmap))
	// 删除数值
	bitmap.Del(8)
	bitmap.Del(15)

	// 验证删除后的位图内容
	if bitmap.Exist(8) {
		t.Errorf("Expected 8 to be removed from the bitmap")
	}

	if bitmap.Exist(15) {
		t.Errorf("Expected 15 to be removed from the bitmap")
	}
	if !bitmap.Exist(0) {
		t.Errorf("Expected 0 to still be in the bitmap")
	}
	if !bitmap.Exist(16) {
		t.Errorf("Expected 16 to still be in the bitmap")
	}
}

func TestExist(t *testing.T) {
	bitmap := NewBitMap()

	// 测试存在性检查
	bitmap.Add(1)
	bitmap.Add(9)
	bitmap.Add(16)

	if !bitmap.Exist(1) {
		t.Errorf("Expected 1 to be in the bitmap")
	}
	if !bitmap.Exist(9) {
		t.Errorf("Expected 9 to be in the bitmap")
	}
	if !bitmap.Exist(16) {
		t.Errorf("Expected 16 to be in the bitmap")
	}
	if bitmap.Exist(2) {
		t.Errorf("Expected 2 not to be in the bitmap")
	}
	if bitmap.Exist(10) {
		t.Errorf("Expected 10 not to be in the bitmap")
	}
}

func TestLen(t *testing.T) {
	bitmap := NewBitMap()

	// 测试位图长度
	bitmap.Add(0)
	bitmap.Add(7)
	bitmap.Add(8)
	bitmap.Add(15)
	bitmap.Add(16)
	bitmap.Add(23)

	expectedLen := 3 // 0-23, so 3 bytes should be needed
	if bitmap.Len() != expectedLen {
		t.Errorf("Expected length %d but got %d", expectedLen, bitmap.Len())
	}
}

func TestToString(t *testing.T) {
	bitmap := NewBitMap()

	// 添加一些数值
	bitmap.Add(0)
	bitmap.Add(1)
	bitmap.Add(7)
	bitmap.Add(8)
	bitmap.Add(15)

	expectedStr := "0000000000000000000000000000000111111110000000"
	if bitmap.ToString() != expectedStr {
		t.Errorf("Expected binary string '%s' but got '%s'", expectedStr, bitmap.ToString())
	}
}

func TestConcurrentAdd(t *testing.T) {
	bitmap := NewBitMap()
	var wg sync.WaitGroup
	const numGoroutines = 100
	const numAdds = 1000

	// 启动多个 goroutine 来并发地添加数据
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := start; j < start+numAdds; j++ {
				bitmap.Add(j)
			}
		}(i * numAdds)
	}

	wg.Wait()

	// 验证某些值是否在位图中
	for i := 0; i < numGoroutines*numAdds; i++ {
		if !bitmap.Exist(i) {
			t.Errorf("Expected %d to be in the bitmap", i)
		}
	}
}

func TestConcurrentDel(t *testing.T) {
	bitmap := NewBitMap()
	const numGoroutines = 100
	const numAdds = 1000
	const numDeletes = 500

	// 添加数据
	for i := 0; i < numGoroutines*numAdds; i++ {
		bitmap.Add(i)
	}

	var wg sync.WaitGroup

	// 启动多个 goroutine 来并发地删除数据
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := start; j < start+numDeletes; j++ {
				bitmap.Del(j)
			}
		}(i * numDeletes)
	}

	wg.Wait()

	// 验证删除的值是否不在位图中
	for i := 0; i < numGoroutines*numDeletes; i++ {
		if bitmap.Exist(i) {
			t.Errorf("Expected %d to be removed from the bitmap", i)
		}
	}
}

func TestConcurrentExist(t *testing.T) {
	bitmap := NewBitMap()
	const numGoroutines = 100
	const numAdds = 1000

	// 添加数据
	for i := 0; i < numGoroutines*numAdds; i++ {
		bitmap.Add(i)
	}

	var wg sync.WaitGroup

	// 启动多个 goroutine 来并发地检查存在性
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := start; j < start+numAdds; j++ {
				if !bitmap.Exist(j) {
					t.Errorf("Expected %d to be in the bitmap", j)
				}
			}
		}(i * numAdds)
	}

	wg.Wait()
}

func TestConcurrentLen(t *testing.T) {
	bitmap := NewBitMap()
	const numGoroutines = 100
	const numAdds = 1000

	// 启动多个 goroutine 来并发地添加数据
	for i := 0; i < numGoroutines; i++ {
		go func(start int) {
			for j := start; j < start+numAdds; j++ {
				bitmap.Add(j)
			}
		}(i * numAdds)
	}

	// 等待所有 goroutines 完成
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numAdds; j++ {
				_ = bitmap.Len() // 调用 Len() 方法来测试
			}
		}()
	}

	wg.Wait()

	// 验证长度是否在预期范围内
	expectedLen := numGoroutines * numAdds / 8
	if bitmap.Len() < expectedLen {
		t.Errorf("Expected length to be at least %d but got %d", expectedLen, bitmap.Len())
	}
}

func TestConcurrentToString(t *testing.T) {
	bitmap := NewBitMap()
	const numGoroutines = 100
	const numAdds = 1000

	// 启动多个 goroutine 来并发地添加数据
	for i := 0; i < numGoroutines; i++ {
		go func(start int) {
			for j := start; j < start+numAdds; j++ {
				bitmap.Add(j)
			}
		}(i * numAdds)
	}

	// 等待所有 goroutines 完成
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numAdds; j++ {
				_ = bitmap.ToString() // 调用 ToString() 方法来测试
			}
		}()
	}

	wg.Wait()
}

func Test_reverse(t *testing.T) {
	s := "123夏楠1231是个m"
	s1 := []rune(s)
	for i := range s {
		fmt.Println(s[i])
	}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s1[i], s1[j] = s1[j], s1[i]

	}
	fmt.Println(string(s1))
}

func Test_append(t *testing.T) {
	a := []int{}
	for i := 0; i < 10; i++ {
		a = append(a, i)
	}

	for i := 0; i < len(a); i++ {
		a = append(a, i)
		fmt.Println(i)
	}
}
