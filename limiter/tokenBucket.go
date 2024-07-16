package limiter

import (
	"sync"
	"time"
)

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	capacity      int        // 容量
	currentTokens int        // 令牌数量
	rate          int        // 发放令牌速率/秒
	lastTime      time.Time  // 上次发放令牌时间
	mutex         sync.Mutex // 避免并发问题
}

// NewTokenBucketLimiter 创建一个新的令牌桶限流器实例。
func NewTokenBucketLimiter(capacity, rate int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity:      capacity,
		rate:          rate,
		lastTime:      time.Now(),
		currentTokens: 0, // 初始化时桶中没有令牌
	}
}

// TryAcquire 尝试从令牌桶中获取一个令牌。
func (l *TokenBucketLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	interval := now.Sub(l.lastTime) // 计算时间间隔

	// 如果距离上次发放令牌超过1秒，则发放新的令牌
	if interval >= time.Second {
		// 计算应该发放的令牌数量，但不超过桶的容量
		newTokens := int(interval/time.Second) * l.rate
		l.currentTokens = minInt(l.capacity, l.currentTokens+newTokens)

		// 更新上次发放令牌的时间
		l.lastTime = now
	}

	// 如果桶中没有令牌，则请求失败
	if l.currentTokens == 0 {
		return false
	}

	// 桶中有令牌，消费一个令牌
	l.currentTokens--

	return true
}

// minInt 返回两个整数中的较小值。
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
