package cache

import "sync"

// ItemCache 用于在列表和map中存储缓存项
type ItemCache struct {
	Key   interface{}
	Value interface{}
}

// LRUCache 包含老年区和青年区
type LRUCache struct {
	Old   *OldCache
	Young *YoungCache
	sync.RWMutex
}

// Access 访问缓存中的项，根据空间局部性决定放在老年区还是青年区
func (l *LRUCache) Access(key interface{}, value interface{}) {
	l.RLock()
	_, inOld := l.Old.Items[key]
	_, inYoung := l.Young.Items[key]
	l.RUnlock()

	if inOld {
		// 如果数据在老年区，直接返回
		return
	} else if inYoung {
		// 如果数据在青年区，可以考虑晋升到老年区
		l.Young.PromoteToOld(key)
	} else {
		// 数据不在缓存中，需要从磁盘加载
		l.loadToCache(key, value)
	}
}
