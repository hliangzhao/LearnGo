package runner

import (
	`errors`
	`os`
	`os/signal`
	`time`
)

var (
	ErrTimeout = errors.New("cannot finish task within the timeout")
	ErrInterrupt = errors.New("received interrupt from OS")
)


// Runner 给定一些tasks，要求在规定的时间内完成，否则报错（timeout或interrupt）
type Runner struct {
	interrupt chan os.Signal        // 声明一个用来传递和接收OS信号的通道
	complete chan error             // 声明一个用来传递和接收任务执行时是否出错的通道
	timeout <- chan time.Time       // 声明一个单向的通道，存放存入数据时的时间
	tasks []func(id int)            // 声明一个任务列表（slice）
}

func New(t time.Duration) *Runner {
	// 创建一个Runner实例并返回其地址
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(t),           // 经过t秒之后，自动将time.Time实例传入timeout这个通道
		tasks:     make([]func(int), 0),
	}
}

// AddTasks 参数是一个变长列表，类似python中的*args
func (r *Runner) AddTasks(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

func (r *Runner) run() error {
	for id, task := range r.tasks {
		// select语句是处理channel数据的if-else逻辑
		select {
		case <- r.interrupt:
			// 如果r.interrupt这个通道里面有数据，则执行本case
			signal.Stop(r.interrupt)
			return ErrInterrupt
		default:
			task(id)
		}
	}
	// 所有tasks都正常执行完毕则返回nil
	return nil
}

func (r *Runner) Start() error {
	// 把os接收到的interrupt传递给Runner的interrupt通道中
	signal.Notify(r.interrupt, os.Interrupt)

	// 异步地执行run方法
	go func() {
		r.complete <- r.run()
	}()

	select {
	case err := <- r.complete:
		// r.complete不为nil当且仅当其执行出现了被中断的问题，即ErrInterrupt
		return err
	case <- r.timeout:
		// r.timeout这个通道里面有数据，则意味着已经到达了设定的时限，则返回超时错误
		return ErrTimeout
	}
}
