package pool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type WorkerPoolTestSuite struct {
	suite.Suite
}

func (suite *WorkerPoolTestSuite) SetupTest() {}

func (suite *WorkerPoolTestSuite) TearDownTest() {}

func TestWorkerPoolTestSuite(t *testing.T) {
	suite.Run(t, new(WorkerPoolTestSuite))
}

func (suite *WorkerPoolTestSuite) TestErrCheckFunc() {
	suite.True(IsErrQueueIsFull(ErrQueueIsFull))
	suite.True(IsErrPoolNotStarted(ErrPoolNotStarted))
	suite.True(IsErrPoolTerminated(ErrPoolTerminated))

	suite.False(IsErrQueueIsFull(ErrPoolTerminated))
	suite.False(IsErrPoolNotStarted(ErrQueueIsFull))
	suite.False(IsErrPoolTerminated(ErrPoolNotStarted))
}

func (suite *WorkerPoolTestSuite) TestWorkerPoolSubmit() {
	opts := NewWorkerPoolOptions(nil)
	pool := NewWorkerPool(1, opts)

	// 向未启动的Pool提交工作
	suite.EqualError(ErrPoolNotStarted, pool.Submit(func() {}, true).Error())

	// 向已经启动的Pool提交工作
	suite.NoError(pool.Start())
	finished := false
	suite.NoError(pool.Submit(func() { finished = true }, true))
	suite.True(finished)

	// 向队列已满的Pool提交工作
	suite.NoError(pool.Submit(func() { time.Sleep(1 * time.Second) }, true))
	// 非阻塞提交，报错
	suite.EqualError(ErrQueueIsFull, pool.Submit(func() {}, false).Error())
	// 阻塞提交，成功
	suite.NoError(pool.Submit(func() {}, true))

	// 向已终止的Pool提交工作
	suite.NoError(pool.Terminate())
	suite.EqualError(ErrPoolTerminated, pool.Submit(func() {}, true).Error())
}

func (suite *WorkerPoolTestSuite) TestWorkerPool() {
	opts := NewWorkerPoolOptions(nil)
	pool := NewWorkerPool(3, opts)

	suite.NoError(pool.Start())

	wg := sync.WaitGroup{}
	expected := 2324
	executed := int32(0)

	for i := 0; i < expected; i++ {
		wg.Add(1)
		suite.NoError(pool.Submit(func() {
			atomic.AddInt32(&executed, 1)
			wg.Done()
		}, true))
	}

	wg.Wait()
	suite.Equal(int32(expected), atomic.LoadInt32(&executed))
	suite.NoError(pool.Terminate())
}

func (suite *WorkerPoolTestSuite) TestWorkerPoolGraceTermination() {
	// 无优雅退出时间，立即退出（GraceTimeout == 0）
	opts := NewWorkerPoolOptions(nil)
	pool := NewWorkerPool(2, opts)
	suite.NoError(pool.Start())
	jobFinished := false
	suite.NoError(pool.Submit(func() {
		time.Sleep(time.Second)
		jobFinished = true
	}, true))
	suite.NoError(pool.Terminate())
	suite.False(jobFinished)

	// 有优雅退出时间，任务应当能够完成
	opts = WorkerPoolOptions{GraceTimeout: time.Second}
	pool = NewWorkerPool(2, opts)
	suite.NoError(pool.Start())
	jobFinished = false
	suite.NoError(pool.Submit(func() {
		time.Sleep(500 * time.Millisecond)
		jobFinished = true
	}, true))
	st := time.Now()
	suite.NoError(pool.Terminate())
	end := time.Now().Sub(st)
	suite.True(jobFinished)
	fmt.Println(end.Seconds())
}

func (suite *WorkerPoolTestSuite) TestWorkerPoolPanicHandler() {
	opts := NewWorkerPoolOptions(nil)
	panicHandled := false
	opts.PanicHandler = func() {
		if err := recover(); err != nil {
			panicHandled = true
		}
	}
	pool := NewWorkerPool(1, opts)
	suite.NoError(pool.Start())
	suite.NoError(pool.Submit(func() { panic("intended panic") }, true))
	suite.True(panicHandled)
}
