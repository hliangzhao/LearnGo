package pool

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

var (
	ErrPoolClosed = errors.New("pool has been closed")
)

// Pool 让协程安全地共享资源
type Pool struct {
	// TODO：factory作为工厂，被协程们以并行方式调用的时候应该加锁。这才符合实际生活生产环境
	factory   func() (io.Closer, error) // 生成资源的工厂方法（我们要求这里的资源实现了close方法，即本资源至少是一个closer interface）
	resources chan io.Closer            // 存放资源的通道
	mtx       sync.Mutex                // 互斥锁
	closed    bool                      // 资源池是否已关闭
}

func New(factory func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("invalid pool size")
	}
	return &Pool{
		factory:   factory,
		resources: make(chan io.Closer, size),
		mtx:       sync.Mutex{},
		closed:    false,
	}, nil
}

// AcquireResource 从pool（的通道）中取出资源来使用
func (p *Pool) AcquireResource() (io.Closer, error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	select {
	case resource, ok := <-p.resources:
		// 如果p.resources有资源，则进行本case，但是有资源也不一定能拿得到，需要pool没有被关闭
		if !ok {
			fmt.Println("pool is closed")
			return nil, ErrPoolClosed
		}
		fmt.Println("successfully get resource from pool")
		return resource, nil
	default:
		// p.resources没有资源，调用工厂方法拿到本资源
		fmt.Println("call the factory to create new resource...")
		return p.factory()
	}
}

// ReleaseResource 资源使用完毕，放回pool（的通道）中
func (p *Pool) ReleaseResource(resource io.Closer) error {
	// TODO：操纵resources这个channel需要加锁
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.closed {
		// channel已经关闭，则直接关闭本资源，默默离开吧
		fmt.Println("pool is closed, do not accept released resource")
		return resource.Close()
	}

	select {
	case p.resources <- resource:
		fmt.Println("successfully release resource to pool")
		return nil
	default:
		// 缓存区已经存满，不再接收新释放的资源，直接关闭本资源吧
		fmt.Println("pool is full, do not accept released resource")
		return resource.Close()
	}
}

func (p *Pool) Close() error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	// 先关闭通道，不允许再向通道中添加资源，然后再依次关闭通道中的所有资源
	close(p.resources)
	for resource := range p.resources {
		if err := resource.Close(); err != nil {
			return err
		}
	}
	return nil
}
