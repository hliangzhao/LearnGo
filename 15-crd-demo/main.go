package main

import (
	"flag"
	"github.com/golang/glog"
	clientset "github.com/hliangzhao/LearnGo/15-crd-demo/pkg/generated/clientset/versioned"
	informers "github.com/hliangzhao/LearnGo/15-crd-demo/pkg/generated/informers/externalversions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	flagSet              = flag.NewFlagSet("crddemo", flag.ExitOnError)
	master               = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	kubeconfig           = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM} // 来自OS的信号
)

func setupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler)

	stop := make(chan struct{})

	c := make(chan os.Signal, 2)
	// 将来自OS的信号relay to c
	signal.Notify(c, shutdownSignals...)

	// 启动一个协程不断监听是否有来自OS的interrupt和terminated的信号
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	return stop
}

func main() {
	flag.Parse()

	stopCh := setupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	mydemoClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	// 将生成的client注入InformerFactor，使得client可以和api-server通信，实现资源的list watch
	mydemoInformerFactory := informers.NewSharedInformerFactory(mydemoClient, time.Second*30)

	// 生成一个crddemo组的Mydemo对象传递给自定义控制器
	controller := NewController(kubeClient, mydemoClient, mydemoInformerFactory.Crddemo().V1().MyDemos())

	go mydemoInformerFactory.Start(stopCh)

	if err := controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Err running controller: %s", err.Error())
	}
}
