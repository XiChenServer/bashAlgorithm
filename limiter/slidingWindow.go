package limiter

import (
	"errors"
	"sync"
	"time"
)

// SlidingWindowLimiter 滑动窗口限流器，用于控制请求的速率。
type SlidingWindowLimiter struct {
	limit        int           // 窗口内允许的最大请求数
	window       int64         // 窗口时间大小（纳秒）
	smallWindow  int64         // 小窗口时间大小（纳秒）
	smallWindows int64         // 窗口内小窗口的数量
	counters     map[int64]int // 每个小窗口的请求计数
	mutex        sync.RWMutex  // 使用读写锁提高并发性能
}

// NewSlidingWindowLimiter 创建并初始化滑动窗口限流器。
func NewSlidingWindowLimiter(limit int, window, smallWindow time.Duration) (*SlidingWindowLimiter, error) {
	if int64(window%smallWindow) != 0 {
		return nil, errors.New("window size must be divisible by the small window size")
	}
	return &SlidingWindowLimiter{
		limit:        limit,
		window:       int64(window),
		smallWindow:  int64(smallWindow),
		smallWindows: int64(window / smallWindow),
		counters:     make(map[int64]int),
	}, nil
}

// TryAcquire 尝试在当前窗口内获取一个请求的机会。
func (l *SlidingWindowLimiter) TryAcquire() bool {
	l.mutex.RLock() // 读锁，允许多个并发读取
	defer l.mutex.RUnlock()

	now := time.Now().UnixNano()
	currentSmallWindow := now / l.smallWindow * l.smallWindow // 当前小窗口的起始点

	// 清理过期的小窗口计数器
	l.cleanExpiredWindows(now)

	// 检查并更新当前小窗口的计数
	l.mutex.Lock() // 写锁，更新计数器
	defer l.mutex.Unlock()

	count, exists := l.counters[currentSmallWindow]
	if !exists || count < l.limit {
		l.counters[currentSmallWindow] = count + 1
		return true
	}
	return false
}

// cleanExpiredWindows 清理已过期的小窗口计数器。
func (l *SlidingWindowLimiter) cleanExpiredWindows(now int64) {
	startSmallWindow := now/l.smallWindow*l.smallWindow - l.window
	for smallWindow := range l.counters {
		if smallWindow < startSmallWindow {
			delete(l.counters, smallWindow)
		}
	}
}

// 注意：cleanExpiredWindows 方法应该在持有读锁的情况下调用，以避免在遍历和修改计数器时产生竞态条件。
