package main

import (
	"fmt"
	"github.com/golang/glog"
	samplecrdv1 "github.com/hliangzhao/LearnGo/15-crd-demo/pkg/apis/crddemo/v1"
	clientset "github.com/hliangzhao/LearnGo/15-crd-demo/pkg/generated/clientset/versioned"
	mydemoscheme "github.com/hliangzhao/LearnGo/15-crd-demo/pkg/generated/clientset/versioned/scheme"
	informers "github.com/hliangzhao/LearnGo/15-crd-demo/pkg/generated/informers/externalversions/crddemo/v1"
	listers "github.com/hliangzhao/LearnGo/15-crd-demo/pkg/generated/listers/crddemo/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"time"
)

const controllerAgentName = "mydemo-controller"

const (
	SuccessSynced         = "Synced"
	MessageResourceSynced = "Mydemo synced successfully"
)

type Controller struct {
	kubeClientSet   kubernetes.Interface
	mydemoClientSet clientset.Interface
	demoInformer    listers.MyDemoLister
	mydemoSynced    cache.InformerSynced

	// workQueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workQueue workqueue.RateLimitingInterface

	// recorder is an event recorder for recording Event resources to the Kubernetes API.
	recorder record.EventRecorder
}

func NewController(
	kubeClientSet kubernetes.Interface,
	mydemoClientSet clientset.Interface,
	mydemoInformer informers.MyDemoInformer) *Controller {

	// Create event broadcaster
	// Add mydemo-controller types to the default Kubernetes Scheme so Events can be
	// logged for mydemo-controller types.
	utilruntime.Must(mydemoscheme.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientSet.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeClientSet:   kubeClientSet,
		mydemoClientSet: mydemoClientSet,
		demoInformer:    mydemoInformer.Lister(),
		mydemoSynced:    mydemoInformer.Informer().HasSynced,
		workQueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "MyDemos"),
		recorder:        recorder,
	}

	glog.Info("Setting up mydemo event handlers")
	// mydemoInformer 注册了三个 Handler（AddFunc、UpdateFunc 和 DeleteFunc）
	// 分别对应 API 对象的“添加”“更新”和“删除”事件。而具体的处理操作，都是将该事件对应的 API 对象加入到工作队列中
	mydemoInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueMydemo,
		UpdateFunc: func(old, new interface{}) {
			oldMydemo := old.(*samplecrdv1.MyDemo)
			newMydemo := new.(*samplecrdv1.MyDemo)
			if oldMydemo.ResourceVersion == newMydemo.ResourceVersion {
				return
			}
			controller.enqueueMydemo(new)
		},
		DeleteFunc: controller.enqueueMydemoForDelete,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shut down the workQueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workQueue.ShutDown()

	glog.Info("Starting MyDemo control loop")

	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.mydemoSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	glog.Info("Started workers")

	// 只有当stopCh closed才可以从里面取数据，这样才可以执行到shutting down的代码
	<-stopCh
	glog.Info("Shutting down workers")
	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workQueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workQueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workQueue.Get()
	if shutdown {
		return false
	}

	// We wrap this block in a func, so we can defer c.workQueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workQueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workQueue and attempted again after a back-off
		// period.
		defer c.workQueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workQueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workQueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workQueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workQueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workQueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workQueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Mydemo resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item, so it does not
		// get queued again until another change happens.
		c.workQueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Mydemo resource
// with the current status of the resource.
// 尝试从 Informer 维护的缓存中拿到了它所对应的 MyDemo 对象
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// 代码中的 demoInformer，从namespace中通过key获取Mydemo对象这个操作，其实就是在访问本地缓存的索引，
	// 实际上，在 Kubernetes 的源码中，你会经常看到控制器从各种 Lister 里获取对象，
	// 比如：podLister、nodeLister 等等，它们使用的都是 Informer 和缓存机制。

	// 而如果控制循环从缓存中拿不到这个对象（demoInformer 返回了 IsNotFound 错误），那就意味着这个 Mydemo
	// 对象的 Key 是通过前面的“删除”事件添加进工作队列的。所以，尽管队列里有这个 Key，但是对应的 Mydemo
	// 对象已经被删除了。而如果能够获取到对应的 Mydemo 对象，就可以执行控制器模式里的对比“期望状态
	// （用户通过 YAML 文件提交到 APIServer 里的信息）”和“实际状态（我们的控制循环需要通过查询实际的Mydemo资源情况”
	// 的功能逻辑了

	// Get the Mydemo resource with this namespace/name
	mydemo, err := c.demoInformer.MyDemos(namespace).Get(name)
	// 从缓存中拿不到这个对象,那就意味着这个 MyDemo 对象的 Key 是通过前面的“删除”事件添加进工作队列的。
	if err != nil {
		// The Mydemo resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			glog.Warningf("DemoCRD: %s/%s does not exist in local cache, will delete it from Mydemo ...",
				namespace, name)

			glog.Infof("[DemoCRD] Deleting mydemo: %s/%s ...", namespace, name)

			return nil
		}

		utilruntime.HandleError(fmt.Errorf("failed to list mydemo by: %s/%s", namespace, name))

		return err
	}

	glog.Infof("[DemoCRD] Try to process mydemo: %#v ...", mydemo)

	c.recorder.Event(mydemo, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// enqueueMydemo takes a Mydemo resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Mydemo.
func (c *Controller) enqueueMydemo(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workQueue.AddRateLimited(key)
}

// enqueueMydemoForDelete takes a deleted Mydemo resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Mydemo.
func (c *Controller) enqueueMydemoForDelete(obj interface{}) {
	var key string
	var err error
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workQueue.AddRateLimited(key)
}
