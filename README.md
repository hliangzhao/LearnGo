## LearnGo


#### [1-get-bilibili](https://github.com/hliangzhao/LearnGo/tree/main/1-get-bilibili)

通过一个简单的http get的例子熟悉Golang的语法。


#### [2-array-slice-map](https://github.com/hliangzhao/LearnGo/tree/main/2-array-slice-map)

介绍Golang中的数组、切片、映射的创建及使用（`make`、传值 or 传址）。


#### [3-type-system](https://github.com/hliangzhao/LearnGo/tree/main/3-type-system)

介绍内置类型、引用类型、结构体、接口以及基于此的、多态的实现。


#### [4-goroutine-channel](https://github.com/hliangzhao/LearnGo/tree/main/4-goroutine-channel)

介绍Golang最重要的特征：协程与通道。

协程是Go语言中实现并发的一种方式，协程之间传递信息和同步的渠道是channel。
此处的协程是轻量级线程，与Python中的协程完全不一样，会面临锁相关的问题。


#### [5-concurrent-model](https://github.com/hliangzhao/LearnGo/tree/main/5-concurrent-model)

介绍Golang中的两种并发模型：Runner和Pool。
* Runner结构体内包含**任务**列表和传递信息的通道；
* Pool内部存放的是共享**资源**。


#### [6-context](https://github.com/hliangzhao/LearnGo/tree/main/6-context)

context主要是用来在协程之间**传递上下文信息**，包括取消信号、超时时间、截止时间以及一些键值对等。
客户端发起取消等信号，服务端收到之后就会终止响应本次请求，从而节约服务端资源。


#### [7-testing-benchmarking](https://github.com/hliangzhao/LearnGo/tree/main/7-testing-benchmarking)

介绍Golang中的单元测试（测试函数是否有bug）和基准测试（测试函数的运行时间）的写法。


#### [8-reflection](https://github.com/hliangzhao/LearnGo/tree/main/8-reflection)

反射机制使得我们可以在程序动态运行时获得输入值的类型等信息。
最典型的案例是`fmt.Println()`，该函数动态获得传入的参数信息从而给出正确的格式输出。


#### [9-app-bili-stream](https://github.com/hliangzhao/LearnGo/tree/main/9-app-bili-stream)

通过`go get`使用外部包。


#### [10-gRPC](https://github.com/hliangzhao/LearnGo/tree/main/10-gRPC)

展示了如何通过`protoc`生成桩代码，并给出了gRPC server和client的一个简单案例。


#### [11-gin](https://github.com/hliangzhao/LearnGo/tree/main/11-gin)

展示了Golang的web框架`Gin`的使用。


#### [12-video-website](https://github.com/hliangzhao/LearnGo/tree/main/12-video-website)

综合运用`Gin`和`React`实现一个简易版B站。


#### [13-Docker](https://github.com/hliangzhao/LearnGo/tree/main/13-Docker)

给出了Dockerfile的一个构建示例。

给出了一个简易版的Docker实现。具体地，需要在代码中调用系统调用创建cgroup等内容对试图运行的程序进行进程隔离。
