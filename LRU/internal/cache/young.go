package cache

import (
	"container/list"
	"sync"
)

type YoungCache struct {
	capacity int
	size     int
	mu       sync.RWMutex
	List     *list.List
	Items    sync.Map
}

// NewYoungCache 初始化OldCache
func NewYoungCache() *YoungCache {
	return &YoungCache{

		List:  list.New(),
		Items: sync.Map{},
	}
}

// Add 向青年区缓存添加一个项
func (c *YoungCache) Add(key interface{}, value interface{}) {
	// 如果青年区已存在该项，则更新值并移动到列表前端
	if _, ok := c.Items.Load(key); ok {
		c.Items.Store(key, &ItemCache{Key: key, Value: value})
		c.List.MoveToFront(c.List.Back()) // 假设是刚刚淘汰过来的，所以是列表最后一个
	} else {
		// 否则，添加新项
		c.Items.Store(key, &ItemCache{Key: key, Value: value})
		c.List.PushBack(&ItemCache{Key: key, Value: value})
	}

	// 如果超出容量，淘汰最老的项
	if c.List.Len() > c.capacity {
		c.Evict()
	}
}

// Evict 从青年区淘汰最老的数据项
func (c *YoungCache) Evict() {
	// 这里直接删除最老的项，没有移动到其他区域的逻辑
	back := c.List.Back()
	if back != nil {
		c.Items.Delete(back.Value.(*ItemCache).Key)
		c.List.Remove(back)
	}
}

// PromoteToOld 将青年区的项目晋升到老年区
func (y *YoungCache) PromoteToOld(key interface{}, o *OldCache) {
	y.mu.Lock()
	defer y.mu.Unlock()
	// 从青年区删除项目
	item, ok := y.Items.Load(key)
	if ok {
		y.List.Remove(item.(*list.Element))
		y.Items.Delete(key)

	}
	// 添加到老年区
	o.Add(key, item, y)
}
