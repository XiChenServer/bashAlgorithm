package limiter

import (
	"sync"
	"testing"
	"time"
)

func TestFixedWindowLimiter(t *testing.T) {
	limiter := NewFixedWindowLimiter(100, time.Second) // 假设我们设置每秒限流100次请求

	// 模拟并发请求的函数
	testRequest := func(wg *sync.WaitGroup, limiter *FixedWindowLimiter) {
		defer wg.Done()
		if !limiter.TryAcquire() {
			t.Errorf("Failed to acquire request limit")
		}
	}

	// 模拟1000个并发请求
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go testRequest(&wg, limiter)
	}

	// 等待所有并发请求完成
	wg.Wait()
}

func BenchmarkFixedWindowLimiter(b *testing.B) {
	limiter := NewFixedWindowLimiter(50, time.Second)

	// 基准测试的循环
	for i := 0; i < b.N; i++ {
		if !limiter.TryAcquire() {
			b.Errorf("Failed to acquire request limit")
		}
	}
}
