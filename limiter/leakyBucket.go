package limiter

import (
	"errors"
	"math"
	"sync"
	"time"
)

// LeakyBucketLimiter 漏桶限流器
type LeakyBucketLimiter struct {
	peakLevel       int          // 最高水位
	currentLevel    int          // 当前水位
	currentVelocity int          // 水流速度/秒
	lastTime        time.Time    // 上次放水时间
	mutex           sync.RWMutex // 使用读写锁提高并发性能
}

// NewLeakyBucketLimiter 初始化漏桶限流器
func NewLeakyBucketLimiter(peakLevel, currentVelocity int) (*LeakyBucketLimiter, error) {
	if currentVelocity <= 0 {
		return nil, errors.New("currentVelocity must be greater than 0")
	}
	if peakLevel < currentVelocity {
		return nil, errors.New("peakLevel must be greater than or equal to currentVelocity")
	}
	return &LeakyBucketLimiter{
		peakLevel:       peakLevel,
		currentLevel:    0, // 初始化时水位为0
		currentVelocity: currentVelocity,
		lastTime:        time.Now(),
	}, nil
}

// TryAcquire 尝试获取处理请求的权限
func (l *LeakyBucketLimiter) TryAcquire() bool {
	l.mutex.RLock() // 读锁，允许多个并发读取
	defer l.mutex.RUnlock()

	// 如果上次放水时间距今不到1秒，不需要放水
	now := time.Now()
	interval := now.Sub(l.lastTime)

	l.mutex.Lock() // 写锁，更新水位
	defer l.mutex.Unlock()

	// 计算放水后的水位
	if interval >= time.Second {
		l.currentLevel = int(math.Max(0, float64(l.currentLevel)-(interval/time.Second).Seconds()*float64(l.currentVelocity)))
		l.lastTime = now
	}
	// 尝试增加水位
	if l.currentLevel < l.peakLevel {
		l.currentLevel++
		return true
	}
	return false
}
