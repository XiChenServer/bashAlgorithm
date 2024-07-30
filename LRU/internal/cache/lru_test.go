package cache

import (
	"fmt"
	"testing"
)

func Test_LRU(t *testing.T) {
	l := NewLRUCache()
	str := l.loadFromDisk("0")
	fmt.Println(str)
}
