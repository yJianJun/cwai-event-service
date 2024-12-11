package pool

import (
	"fmt"
	"sync"
	"time"
)

// WorkerPool是一个轻量级的goroutine池，支持优雅退出
// 使用示例:
//
//	// 初始化配置
//	opts := NewWorkerPoolOptions(func() {
//		if err := recover(); err != nil {
//			glog.Error(err)
//		}
//	})
//	// 初始化WorkerPool，10个Worker
//	pool := NewWorkerPool(10, opts)
//	// 启动WorkerPool
//	if err := pool.Start(); err != nil {
//		glog.Error(err)
//	}
//	// 阻塞提交任务
//	if err := pool.Submit(func() { /*do something*/ }, true); err != nil {
//		glog.Error(err)
//	}
//	// 非阻塞提交任务
//	if err := pool.Submit(func() { /*do something*/ }, false); err != nil {
//		if IsErrQueueIsFull(err) {
//			// 队列已满处理逻辑
//		} else {
//			glog.Error(err)
//		}
//	}
//	// 终止WorkerPool
//	if err := pool.Terminate(); err != nil {
//		glog.Error(err)
//	}

const (
	NoGraceTimeout      time.Duration = 0 // 无优雅退出时间，立即退出
	UnbufferedQueueSize               = 0 // 工作队列无buffer，提交工作可能blocking
)

// WorkerPoolOptions WorkerPool配置结构体
type WorkerPoolOptions struct {
	GraceTimeout   time.Duration // 优雅退出时间
	MaxQueuingSize int           // 队列大小
	PanicHandler   func()        // panic处理函数
}

// NewWorkerPoolOptions 初始化一个默认WorkerPool配置结构体
func NewWorkerPoolOptions(panicHandler func()) WorkerPoolOptions {
	return WorkerPoolOptions{
		GraceTimeout:   NoGraceTimeout,
		MaxQueuingSize: UnbufferedQueueSize,
		PanicHandler:   panicHandler,
	}
}

var (
	ErrQueueIsFull    = fmt.Errorf("job queue is full")
	ErrPoolNotStarted = fmt.Errorf("pool not started")
	ErrPoolTerminated = fmt.Errorf("pool has already been terminated")
)

func IsErrQueueIsFull(err error) bool {
	return err == ErrQueueIsFull
}

func IsErrPoolNotStarted(err error) bool {
	return err == ErrPoolNotStarted
}

func IsErrPoolTerminated(err error) bool {
	return err == ErrPoolTerminated
}

// WorkerPool 基于goroutine的Worker池
type WorkerPool struct {
	size      int
	graceTime time.Duration
	maxQLen   int

	panicHandler func()

	funcChan chan func()
	termChan chan struct{}

	wg sync.WaitGroup
}

// NewWorkerPool 初始化一个WorkerPool，size为worker(goroutine)数量
func NewWorkerPool(size int, opts WorkerPoolOptions) *WorkerPool {
	if opts.MaxQueuingSize < 0 {
		panic("negative queue size")
	}
	return &WorkerPool{
		size:         size,
		graceTime:    opts.GraceTimeout,
		maxQLen:      opts.MaxQueuingSize,
		panicHandler: opts.PanicHandler,
		funcChan:     nil,
		termChan:     make(chan struct{}),
	}
}

// Start 启动WorkerPool
func (p *WorkerPool) Start() error {
	p.funcChan = make(chan func(), p.maxQLen)

	for i := 0; i < p.size; i++ {
		p.wg.Add(1)
		go p.run()
	}

	return nil
}

// Terminate 停止WorkerPool
func (p *WorkerPool) Terminate() error {
	grace := time.NewTimer(p.graceTime)
	defer grace.Stop()

	close(p.termChan)

	waitChan := make(chan struct{})

	go func() {
		p.wg.Wait()
		close(waitChan)
	}()

	for {
		select {
		case <-waitChan:
			return nil
		case <-grace.C:
			return nil
		}
	}
}

// Submit 向WorkerPool提交工作
func (p *WorkerPool) Submit(f func(), blocking bool) error {
	if p.funcChan == nil {
		return ErrPoolNotStarted
	}

	select {
	case _, opened := <-p.termChan:
		if !opened {
			return ErrPoolTerminated
		}
	default:
	}

	if blocking {
		p.funcChan <- f
		return nil
	}

	select {
	case p.funcChan <- f:
		return nil
	default:
		return ErrQueueIsFull
	}
}

// run 启动worker
func (p *WorkerPool) run() {
	defer p.wg.Done()

	if p.panicHandler != nil {
		defer p.panicHandler()
	}

	for {
		select {
		case f := <-p.funcChan:
			f()
		case <-p.termChan:
			return
		}
	}

}
