package cache

import (
	"fmt"
	"testing"
)

// 测试LRU加载哪一个数据的时候有没有问题
func Test_LRU(t *testing.T) {
	l := NewLRUCache()
	str := l.loadFromDisk("9")
	fmt.Println(str)
	// 遍历list并打印每个元素
	fmt.Println("old list")
	for e := l.Old.List.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	// 使用Range方法遍历map并打印键值对
	fmt.Println("old Map")
	l.Old.Items.Range(func(key, value interface{}) bool {
		fmt.Printf("key: %v, value: %v\n", key, value)
		return true
	})
	fmt.Println("y list")
	// 遍历list并打印每个元素
	for e := l.Young.List.Front(); e != nil; e = e.Next() {

		fmt.Println(e.Value)
	}
	// 使用Range方法遍历map并打印键值对
	fmt.Println("y Map")
	l.Young.Items.Range(func(key, value interface{}) bool {
		fmt.Printf("key: %v, value: %v\n", key, value)
		return true
	})
}
