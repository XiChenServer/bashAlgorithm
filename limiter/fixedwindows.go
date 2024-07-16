package limiter

import (
	"sync"
	"time"
)

// FixedWindowLimiter 结构代表一个固定窗口限流器。
// 它使用固定大小的时间窗口来限制请求的数量。
type FixedWindowLimiter struct {
	limit    int           // 请求上限，窗口内允许的最大请求数
	window   time.Duration // 窗口时间大小，即时间窗口的长度
	counter  int           // 计数器，记录当前窗口内的请求数
	lastTime time.Time     // 上一次请求的时间
	mutex    sync.Mutex    // 互斥锁，用于同步，避免并发访问导致的问题
}

// NewFixedWindowLimiter 构造函数创建并初始化一个新的 FixedWindowLimiter 实例。
func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		limit:    limit,
		window:   window,
		lastTime: time.Now(), // 初始化时设置当前时间为窗口开始时间
	}
}

// TryAcquire 尝试获取一个请求的机会。
// 如果当前窗口内请求数未达到上限，增加计数器并返回 true。
// 如果请求数已达到上限或窗口已过期，返回 false。
func (l *FixedWindowLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	// 检查当前时间与上次请求时间差是否超过窗口大小
	if now.Sub(l.lastTime) > l.window {
		l.counter = 0    // 如果窗口过期，重置计数器
		l.lastTime = now // 更新窗口开始时间为当前时间
	}
	// 如果当前请求数未达到上限，允许请求
	if l.counter < l.limit {
		l.counter++ // 请求成功，增加计数器
		return true // 返回 true 表示请求已成功获取
	}
	// 如果请求数已达到上限，请求失败
	return false
}
