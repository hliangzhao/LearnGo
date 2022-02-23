package main

import (
	"fmt"
	"github.com/hliangzhao/LearnGo/5-concurrent-model/pool"
	"github.com/hliangzhao/LearnGo/5-concurrent-model/runner"
	"io"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

/*
两种并发模型：runner和pool
*/

func main() {
	// 1 测试有缓冲区通道
	testBufferedChannel()

	// 2 测试runner
	r := runner.New(4 * time.Second)
	r.AddTasks(createTask(), createTask(), createTask())
	err := r.Start()
	switch err {
	case runner.ErrInterrupt:
		fmt.Printf("Task interrupted\n")
	case runner.ErrTimeout:
		fmt.Printf("Task timeout\n")
	default:
		fmt.Println("all finished")
	}

	// 3 测试pool
	p, err := pool.New(Factory, 5)
	if err != nil {
		log.Fatalln(err)
	}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		str := "select user from table " + strconv.Itoa(i+1)
		go Query(str, p)
	}
	wg.Wait()

	if err := p.Close(); err != nil {
		log.Fatalln(err)
	}
}

// testBufferedChannel 使用有缓冲区的通道
func testBufferedChannel() {
	// 创建一个有缓存的通道
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	fmt.Println("Here is reachable")

	// 对于有缓存空间的通道，往里面发送数据的操作只有在通道满的时候才会堵塞
	// ch <- 5
	// fmt.Println("Here is unreachable")

	//
	// TODO：没有缓存的通道，如果关闭了，则ok为false；如果通道为空，则dead lock；
	//  有缓存的通道，只要里面有东西，就可以按照存放顺序取出来，并且ok为true。缓冲区大小其实就是仓库中货架的个数
	//  channel里面放置的是"商品"，取出来1个就少1个
	value, ok := <-ch
	fmt.Println(value, ok)

	// TODO：有缓存的通道，只有close之后才可以"通过for range"往外拿东西
	close(ch)
	for val := range ch {
		fmt.Println(val)
	}
	// 此时通道里的数据已经全部被取出，所以本轮for range不会输出任何东西
	for val := range ch {
		fmt.Println(val)
	}
}

/* for runner */
func createTask() func(int) {
	return func(id int) {
		t := rand.Int()%10 + 1
		time.Sleep(time.Second + time.Duration(t))
		fmt.Printf("Task %d complete.\n", id)
	}
}

/* for pool */
var counter int32

// DBConnection 自定义一种模拟数据库连接的资源
type DBConnection struct {
	id int32
}

// Close 为这种资源定义Close方法使其成为一个Closer
func (conn DBConnection) Close() error {
	fmt.Printf("DB connection #%d is closed\n", conn.id)
	return nil
}

func Factory() (io.Closer, error) {
	atomic.AddInt32(&counter, 1)
	return DBConnection{id: counter}, nil
}

var wg sync.WaitGroup

func Query(sql string, pool *pool.Pool) {
	defer wg.Done()

	conn, err := pool.AcquireResource()
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err := pool.ReleaseResource(conn); err != nil {
			log.Fatalln(err)
		}
	}()

	// use conn to do some DB operation
	queryTime := time.Duration(rand.Intn(2)+1) * time.Second
	time.Sleep(time.Second)
	fmt.Printf("`%s` with connection %v in %v\n", sql, conn, queryTime)
}
