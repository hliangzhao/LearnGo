package pool

import (
	`errors`
	`fmt`
	`io`
	`sync`
)

var (
	ErrPoolClosed = errors.New("pool has been closed")
)

// Pool 让协程安全地共享资源
type Pool struct {
	factory func() (io.Closer, error)    // 生成资源的工厂方法（我们要求这里的资源实现了close方法，即本资源至少是一个closer interface）
	resources chan io.Closer             // 存放资源的通道
	mtx sync.Mutex                       // 互斥锁
	closed bool                          // 资源池是否已关闭
}

func New(factory func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("invalid pool size")
	}
	return &Pool{
		factory:   factory,
		resources: make(chan io.Closer, size),
		closed:    false,
	}, nil
}

// AcquireResource 从pool（的通道）中取出资源来使用
func (p *Pool) AcquireResource() (io.Closer, error) {
	select {
	case resource, ok := <- p.resources:
		// 如果p.resources有资源，则进行本case，但是有资源也不一定能拿得到，需要pool没有被关闭
		if !ok {
			return nil, ErrPoolClosed
		}
		fmt.Println("get resource from the pool")
		return resource, nil
	default:
		// p.resources没有资源，调用工厂方法拿到本资源
		fmt.Println("acquire new resource")
		return p.factory()
	}
}

// ReleaseResource 资源使用完毕，放回pool（的通道）中
func (p *Pool) ReleaseResource(resource io.Closer) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.closed {
		resource.Close()
		return
	}

	select {
	case p.resources <- resource:
		fmt.Println("release resource back to pool")
	default:
		// 缓存区已经存满，不再接收新释放的资源
		fmt.Println("release resource closed")
		resource.Close()
	}
}

func (p *Pool) Close() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.closed {
		return
	}

	p.closed = true
	// 先关闭通道，不允许再向通道中添加资源，然后再以此关闭通道中的所有资源
	close(p.resources)
	for resource := range p.resources {
		resource.Close()
	}
}