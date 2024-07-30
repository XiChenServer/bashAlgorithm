package cache

import (
	"container/list"
	"sync"
)

// OldCache 用于存储旧数据的缓存结构
type OldCache struct {
	capacity int
	size     int
	mu       sync.RWMutex
	List     *list.List
	Items    sync.Map
}

// NewOldCache 初始化OldCache
func NewOldCache() *OldCache {
	return &OldCache{
		List:  list.New(),
		Items: sync.Map{},
	}
}

// Add 向老年区缓存添加一个项，如果老年区满了，则将最老的项移动到青年区
func (c *OldCache) Add(key interface{}, value interface{}, y *YoungCache) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查项是否已存在
	if _, exists := c.Items.Load(key); !exists {
		item := &ItemCache{Key: key, Value: value}
		c.Items.Store(key, item) // 存储项到sync.Map
		c.List.PushBack(item)    // 将新项添加到双向链表的末尾

		// 如果老年区已满，需要淘汰最老的项
		if c.List.Len() > c.capacity {
			c.evict(y) // 淘汰最老的项，可能移动到青年区
		}
	} else {
		// 如果项已存在，更新其在双向链表中的位置，表示最近被访问
		element, _ := c.Items.Load(key)
		c.List.MoveToFront(element.(*list.Element))
	}
}

// evict 从老年区缓存中淘汰最老的项
func (c *OldCache) evict(young *YoungCache) {
	back := c.List.Back()
	if back != nil {
		item := back.Value.(*ItemCache)
		c.Items.Delete(item.Key)
		c.List.Remove(back)
		// 将淘汰的项添加到青年区
		young.Add(item.Key, item.Value)
	}
}
