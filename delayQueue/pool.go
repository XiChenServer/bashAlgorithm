package delayQueue

import (
	"bash_algorithm/delayQueue/redis"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// return if pool size
	ErrInvalidPoolCap = errors.New("invalid pool cap")
	// put task but pool already closed
	ErrPoolAlreadyClosed = errors.New("pool already closed")
)

// running status
const (
	RUNNING = 1
	STOPED  = 0
)

// Task task to-do
type Task struct {
	//Handler func(v ...interface{})
	//Params []interface{}
	Member string
	Wg     *sync.WaitGroup
}

// Pool task pool
type Pool struct {
	capacity       uint64
	runningWorkers uint64
	status         int64
	chTask         chan *Task
	PanicHandler   func(interface{})
	sync.Mutex
}

// NewPool init pool
func NewPool(capacity uint64) (*Pool, error) {
	if capacity <= 0 {
		return nil, ErrInvalidPoolCap
	}
	p := &Pool{
		capacity: capacity,
		status:   RUNNING,
		chTask:   make(chan *Task, capacity),
	}
	return p, nil
}

func (p *Pool) checkWorker() {
	p.Lock()
	defer p.Unlock()
	if p.runningWorkers == 0 && len(p.chTask) > 0 {
		p.run()
	}
}

// GetCap get capacity
func (p *Pool) GetCap() uint64 {
	return p.capacity
}

// GetRunningWorkers get running workers
func (p *Pool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&p.runningWorkers)
}

func (p *Pool) incRunning() {
	atomic.AddUint64(&p.runningWorkers, 1)
}

func (p *Pool) decRunning() {
	atomic.AddUint64(&p.runningWorkers, ^uint64(0))
}

// Put put a task to pool
func (p *Pool) Put(task *Task) error {
	p.Lock()
	defer p.Unlock()
	if p.status == STOPED {
		return ErrPoolAlreadyClosed
	}
	//run workers
	if p.GetRunningWorkers() < p.GetCap() {
		p.run()
	}
	//send task
	if p.status == RUNNING {
		p.chTask <- task
	}
	return nil
}

func (p *Pool) run() {
	p.incRunning()
	redisConn := redis.GetRedisDb()
	conn := context.Background()
	go func() {
		defer func() {
			p.decRunning()
			if r := recover(); r != nil {
				if p.PanicHandler != nil {
					p.PanicHandler(r)
				} else {
					log.Printf("Worker panic: %s\n", r)
				}
			}
			p.checkWorker() //check worker avoid no worker running
		}()

		for {
			select {
			case task, ok := <-p.chTask:
				if !ok {
					return
				}
				//task.Handler(task.Params...)
				fmt.Println(task.Member)
				err := ConsumeQueue(redisConn, conn, task.Wg)
				if err != nil {
					fmt.Println("消费队列发生错误：", err)
					return
				}
			}
		}
	}()
}

func (p *Pool) setStatus(status int64) bool {
	p.Lock()
	defer p.Unlock()
	if p.status == status {
		return false
	}
	p.status = status
	return true
}

// Close close pool graceful
func (p *Pool) Close() {
	if !p.setStatus(STOPED) {
		//stop put task
		return
	}
	for len(p.chTask) > 0 {
		//wait all task be consumed
		time.Sleep(1e6) //reduce cpu load
	}
	close(p.chTask)
}
