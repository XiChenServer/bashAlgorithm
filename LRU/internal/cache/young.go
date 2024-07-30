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
		capacity: 1024,
		List:     list.New(),
		Items:    sync.Map{},
	}
}

// Add 向青年区缓存添加一个项
func (c *YoungCache) Add(key interface{}, value interface{}) {
	// 如果青年区已存在该项，则更新值并移动到列表前端
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查项是否已存在
	if element, ok := c.Items.Load(key); ok {
		// 如果项存在，更新值并移动到列表前端
		item := element.(*ItemCache) // 假设存储的是 *ItemCache 类型
		item.Value = value.(string)  // 假设值是 string 类型，根据实际情况调整
		c.List.MoveToFront(element.(*list.Element))
	} else {
		// 如果项不存在，创建新项并添加到列表和 Map 中
		newItem := &ItemCache{Key: key, Value: value}
		c.Items.Store(key, newItem)
		c.List.PushBack(newItem)
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
