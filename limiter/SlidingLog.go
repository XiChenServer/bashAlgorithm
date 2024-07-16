package limiter

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// ViolationStrategyError 定义了违背限流策略时的错误结构。
type ViolationStrategyError struct {
	Limit  int           // 策略的请求上限
	Window time.Duration // 策略的窗口时间大小
}

// Error 实现了 error 接口，返回违背策略的错误信息。
func (e *ViolationStrategyError) Error() string {
	return fmt.Sprintf("violation strategy that limit = %d and window = %v", e.Limit, e.Window)
}

// SlidingLogLimiterStrategy 定义了滑动日志限流器的单个策略。
type SlidingLogLimiterStrategy struct {
	limit        int   // 该策略的窗口请求上限
	window       int64 // 该策略的窗口时间大小（纳秒）
	smallWindows int64 // 该策略窗口内的小窗口数量
}

// NewSlidingLogLimiterStrategy 创建并初始化一个新的滑动日志限流策略。
func NewSlidingLogLimiterStrategy(limit int, window time.Duration) *SlidingLogLimiterStrategy {
	return &SlidingLogLimiterStrategy{
		limit:  limit,
		window: int64(window),
	}
}

// SlidingLogLimiter 定义了滑动日志限流器，它包含多个策略。
type SlidingLogLimiter struct {
	strategies  []*SlidingLogLimiterStrategy // 滑动日志限流器的策略列表
	smallWindow int64                        // 小窗口时间大小（纳秒）
	counters    map[int64]int                // 每个小窗口的请求计数
	mutex       sync.Mutex                   // 互斥锁，避免并发问题
}

// NewSlidingLogLimiter 创建并初始化一个新的滑动日志限流器。
func NewSlidingLogLimiter(smallWindow time.Duration, strategies ...*SlidingLogLimiterStrategy) (*SlidingLogLimiter, error) {
	// 复制策略以避免外部修改
	strategiesCopy := make([]*SlidingLogLimiterStrategy, len(strategies))
	copy(strategiesCopy, strategies)

	// 检查策略列表是否为空
	if len(strategiesCopy) == 0 {
		return nil, errors.New("must be set strategies")
	}

	// 根据窗口大小对策略进行排序
	sort.Slice(strategiesCopy, func(i, j int) bool {
		if strategiesCopy[i].window == strategiesCopy[j].window {
			return strategiesCopy[i].limit > strategiesCopy[j].limit // 窗口相同，限制更大的排前面
		}
		return strategiesCopy[i].window > strategiesCopy[j].window // 窗口大的排前面
	})

	// 验证策略设置的合理性
	for i, strategy := range strategiesCopy {
		if i > 0 && strategy.limit >= strategiesCopy[i-1].limit {
			return nil, errors.New("the smaller window should have the smaller limit")
		}
		if strategy.window%int64(smallWindow) != 0 {
			return nil, errors.New("window cannot be split by integers")
		}
		strategy.smallWindows = strategy.window / int64(smallWindow)
	}

	return &SlidingLogLimiter{
		strategies:  strategiesCopy,
		smallWindow: int64(smallWindow),
		counters:    make(map[int64]int),
	}, nil
}

// TryAcquire 尝试根据设置的策略进行限流。
func (l *SlidingLogLimiter) TryAcquire() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now().UnixNano()
	currentSmallWindow := now / l.smallWindow * l.smallWindow // 当前小窗口的起始点

	// 计算每个策略的起始小窗口值
	startSmallWindows := make([]int64, len(l.strategies))
	for i, strategy := range l.strategies {
		startSmallWindows[i] = currentSmallWindow - l.smallWindow*(strategy.smallWindows-1)
	}

	// 清理过期的小窗口计数器并计算每个策略的当前请求总数
	counts := make([]int, len(l.strategies))
	for smallWindow, counter := range l.counters {
		if smallWindow < startSmallWindows[0] {
			delete(l.counters, smallWindow)
		}
		for i := range l.strategies {
			if smallWindow >= startSmallWindows[i] {
				counts[i] += counter
			}
		}
	}

	// 检查是否违背了策略
	for i, strategy := range l.strategies {
		if counts[i] > strategy.limit {
			return &ViolationStrategyError{
				Limit:  strategy.limit,
				Window: time.Duration(strategy.window),
			}
		}
	}

	// 如果没有违背策略，增加当前小窗口的计数
	l.counters[currentSmallWindow]++
	return nil
}
