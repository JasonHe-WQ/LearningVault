[一文读懂 K8s controller-runtime](https://mp.weixin.qq.com/s/rmcLiTF-vg0qIYVaVhPJUA)
================================================================================

k8s技术圈

作者：蔡奇，中国移动云能力中心软件研发工程师，专注于云原生、微服务、算力网络等领域。

前言
--

在K8s开发中，经常能听过controller的概念，那么这些概念在K8s底层是如何实现，本文将详细介绍。

Controller
----------

在K8s中，实现一个controller是通过controller-runtime(https://github.com/kubernetes-sigs/controller-runtime) 框架来实现的，包括Kubebuilder、operator-sdk等工具也只是在controller-runtime上做了封装，以便开发者快速生成项目的脚手架而已。

Controller定义在pkg/internal/controller/controller，一个controller主要包含Watch和Start两个方法，以及一个调协方法Reconcile。在controller的定义中，看上去没有资源对象的Informer或者Indexer数据，而在K8s中所有与kube-apiserver资源的交互是通过Informer实现的，实际上这里是通过下面的 startWatches 属性做了一层封装。

```go
type Controller struct {
   // Name is used to uniquely identify a Controller in tracing, logging and monitoring.  Name is required.
   Name string

   // MaxConcurrentReconciles is the maximum number of concurrent Reconciles which can be run. Defaults to 1.
   MaxConcurrentReconciles int

   // Reconciler is a function that can be called at any time with the Name / Namespace of an object and
   // ensures that the state of the system matches the state specified in the object.
   // Defaults to the DefaultReconcileFunc.
   Do reconcile.Reconciler

   // MakeQueue constructs the queue for this controller once the controller is ready to start.
   // This exists because the standard Kubernetes workqueues start themselves immediately, which
   // leads to goroutine leaks if something calls controller.New repeatedly.
   MakeQueue func() workqueue.RateLimitingInterface

   // Queue is an listeningQueue that listens for events from Informers and adds object keys to
   // the Queue for processing
   Queue workqueue.RateLimitingInterface

   // SetFields is used to inject dependencies into other objects such as Sources, EventHandlers and Predicates
   // Deprecated: the caller should handle injected fields itself.
   SetFields func(i interface{}) error

   // mu is used to synchronize Controller setup
   mu sync.Mutex

   // Started is true if the Controller has been Started
   Started bool

   // ctx is the context that was passed to Start() and used when starting watches.
   //
   // According to the docs, contexts should not be stored in a struct: https://golang.org/pkg/context,
   // while we usually always strive to follow best practices, we consider this a legacy case and it should
   // undergo a major refactoring and redesign to allow for context to not be stored in a struct.
   ctx context.Context

   // CacheSyncTimeout refers to the time limit set on waiting for cache to sync
   // Defaults to 2 minutes if not set.
   CacheSyncTimeout time.Duration

   // startWatches maintains a list of sources, handlers, and predicates to start when the controller is started.
   startWatches []watchDescription

   // LogConstructor is used to construct a logger to then log messages to users during reconciliation,
   // or for example when a watch is started.
   // Note: LogConstructor has to be able to handle nil requests as we are also using it
   // outside the context of a reconciliation.
   LogConstructor func(request *reconcile.Request) logr.Logger

   // RecoverPanic indicates whether the panic caused by reconcile should be recovered.
   RecoverPanic *bool
}
```

### Watch()

Watch方法首先会判断当前的controller是否已启动，如果未启动，会将watch的内容暂存到startWatches中等待controller启动。如果已启动，则会直接调用src.Start(c.ctx, evthdler, c.Queue, prct...), 其中Source可以为informer、kind、channel等。

```go
// Watch implements controller.Controller.
func (c *Controller) Watch(src source.Source, evthdler handler.EventHandler, prct ...predicate.Predicate) error {
    ...
    // Controller hasn't started yet, store the watches locally and return.
    //
    // These watches are going to be held on the controller struct until the manager or user calls Start(...).
    if !c.Started {
        c.startWatches = append(c.startWatches, watchDescription{src: src, handler: evthdler, predicates: prct})
        return nil
    }

    c.LogConstructor(nil).Info("Starting EventSource", "source", src)
    return src.Start(c.ctx, evthdler, c.Queue, prct...)
}
```

以informer为例，会通过以下方法添加对应的EventHandler：

```go
_, err := is.Informer.AddEventHandler(internal.EventHandler{Queue: queue, EventHandler: handler, Predicates: prct})
```

以kind为例，会通过以下方法添加对应的EventHandler：

```go
i, lastErr = ks.cache.GetInformer(ctx, ks.Type)
_, err := i.AddEventHandler(internal.EventHandler{Queue: queue, EventHandler: handler, Predicates: prct})
```

internal.EventHandler 实现了 OnAdd、OnUpdate、OnDelete 几个方法。也就是说src.Start方法作用是获取对应的informer，并注册对应的EventHandler。

### Start()

Start方法有两个主要功能，一是调用所有startWatches中Source的start方法，注册EventHandler。

```go
for _, watch := range c.startWatches {
   c.LogConstructor(nil).Info("Starting EventSource", "source", fmt.Sprintf("%s", watch.src))

   if err := watch.src.Start(ctx, watch.handler, c.Queue, watch.predicates...); err != nil {
      return err
   }
}
```

二是启动Work来处理资源对象。

```go
for i := 0; i < c.MaxConcurrentReconciles; i++ {
   go func() {
      defer wg.Done()
      // Run a worker thread that just dequeues items, processes them, and marks them done.
      // It enforces that the reconcileHandler is never invoked concurrently with the same object.
      for c.processNextWorkItem(ctx) {
      }
   }()
}
```

processNextWorkItem从Queue中获取资源对象，reconcileHandler 函数就是我们真正执行元素业务处理的地方，函数中包含了事件处理以及错误处理，真正的事件处理是通过c.Do.Reconcile(req) 暴露给开发者的，所以对于开发者来说，只需要在 Reconcile 函数中去处理业务逻辑就可以了。

```go
func (c *Controller) processNextWorkItem(ctx context.Context) bool {
    obj, shutdown := c.Queue.Get()
    if shutdown {
        // Stop working
        return false
    }

    // We call Done here so the workqueue knows we have finished
    // processing this item. We also must remember to call Forget if we
    // do not want this work item being re-queued. For example, we do
    // not call Forget if a transient error occurs, instead the item is
    // put back on the workqueue and attempted again after a back-off
    // period.
    defer c.Queue.Done(obj)

    ctrlmetrics.ActiveWorkers.WithLabelValues(c.Name).Add(1)
    defer ctrlmetrics.ActiveWorkers.WithLabelValues(c.Name).Add(-1)

    c.reconcileHandler(ctx, obj)
    return true
}

// Reconcile implements reconcile.Reconciler.
func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (_ reconcile.Result, err error) {
    defer func() {
        if r := recover(); r != nil {
            if c.RecoverPanic != nil && *c.RecoverPanic {
                for _, fn := range utilruntime.PanicHandlers {
                    fn(r)
                }
                err = fmt.Errorf("panic: %v [recovered]", r)
                return
            }

            log := logf.FromContext(ctx)
            log.Info(fmt.Sprintf("Observed a panic in reconciler: %v", r))
            panic(r)
        }
    }()
    return c.Do.Reconcile(ctx, req)
}
```

### Reconcile

Controller的调协逻辑在Reconcile中执行。

```go
type Reconciler interface {
   // Reconcile performs a full reconciliation for the object referred to by the Request.
   // The Controller will requeue the Request to be processed again if an error is non-nil or
   // Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
   Reconcile(context.Context, Request) (Result, error)
}
type Request struct {
   // NamespacedName is the name and namespace of the object to reconcile.
   types.NamespacedName
}
```

Reconcile方法的入参Request来自于controller.queue，并且会判断队列中的数据类型是否为Reconcile.Request，如果数据类型不一致，则不会执行Reconcile的逻辑。

```go
func (c *Controller) reconcileHandler(ctx context.Context, obj interface{}) {
   // Make sure that the object is a valid request.
   req, ok := obj.(reconcile.Request)
   if !ok {
      // As the item in the workqueue is actually invalid, we call
      // Forget here else we'd go into a loop of attempting to
      // process a work item that is invalid.
      c.Queue.Forget(obj)
      c.LogConstructor(nil).Error(nil, "Queue item was not a Request", "type", fmt.Sprintf("%T", obj), "value", obj)
      // Return true, don't take a break
      return
   }
}
```

那么数据是如何进入队列queue的呢，实际是通过Informer中的EventHandler入队的。回到src.Start(c.ctx, evthdler, c.Queue, prct...)方法中，该方法为informer注册了一个internal.EventHandler。internal.EventHandler实现了OnAdd、OnUpdate、OnDelete等方法，以OnAdd方法为例，该方法最后会调用EventHandler.Create 方法。

```go
type EventHandler struct {
    EventHandler handler.EventHandler
    Queue        workqueue.RateLimitingInterface
    Predicates   []predicate.Predicate
}

// OnAdd creates CreateEvent and calls Create on EventHandler.
func (e EventHandler) OnAdd(obj interface{}) {
    c := event.CreateEvent{}

    // Pull Object out of the object
    if o, ok := obj.(client.Object); ok {
        c.Object = o
    } else {
        log.Error(nil, "OnAdd missing Object",
            "object", obj, "type", fmt.Sprintf("%T", obj))
        return
    }

    for _, p := range e.Predicates {
        if !p.Create(c) {
            return
        }
    }

    // Invoke create handler
    e.EventHandler.Create(c, e.Queue)
}
```

EventHandler为一个接口，有EnqueueRequestForObject、Funcs、EnqueueRequestForOwner、enqueueRequestsFromMapFunc四个实现类。以EnqueueRequestForObject为例，其create方法为：

```go
// Create implements EventHandler.
func (e *EnqueueRequestForObject) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
    if evt.Object == nil {
        enqueueLog.Error(nil, "CreateEvent received with no metadata", "event", evt)
        return
    }
    q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
        Name:      evt.Object.GetName(),
        Namespace: evt.Object.GetNamespace(),
    }})
}
```

所以Reconcile协调执行的数据对象，实际是通过Informer中的EventHandler入队的。

kubebuilder等脚手架框架解析
-------------------

利用kubebuilder、operator-sdk等框架，可以快速生成相应资源对象的controller代码。接下来，以kubebuilder为例，对controller代码逻辑进行解析。

一个完整的controller启动逻辑包含以下步骤：

1） 在main.go启动函数中，会定义一个controllerManager对象。

```go
mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
    Scheme:                 scheme,
    MetricsBindAddress:     metricsAddr,
    Port:                   9443,
    HealthProbeBindAddress: probeAddr,
    LeaderElection:         enableLeaderElection,
    LeaderElectionID:       "9a82ee0d.my.domain",
    CertDir:                "dir",
    ...
})
```

2）通过SetUpWithManager()方法，注册每种资源对象的controller到controllerManager对象中。

```go
if err = (&controllers.AppServiceReconciler{
    Client: mgr.GetClient(),
    Scheme: mgr.GetScheme(),
}).SetupWithManager(mgr); err != nil {
    setupLog.Error(err, "unable to create controller", "controller", "AppService")
    os.Exit(1)
}
```

3）启动controllerManager，也即启动对应资源对象的controller。

```go
if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
    setupLog.Error(err, "problem running manager")
    os.Exit(1)
}
```

主要的代码逻辑在于SetUpWithManager()和mgr.Start()这两个方法中。

```go
// SetupWithManager sets up the controller with the Manager.
func (r *AppServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&appexamplecomv1.AppService{}).
        Complete(r)
}
```

### Builder

ctrl.NewControllerManagedBy(mgr)会返回一个builder对象。

```go
NewControllerManagedBy = builder.ControllerManagedBy

func ControllerManagedBy(m manager.Manager) *Builder {
    return &Builder{mgr: m}
}
```

builder为controller的构造器，其结构定义为:

```go
type Builder struct {
    forInput         ForInput
    ownsInput        []OwnsInput
    watchesInput     []WatchesInput
    mgr              manager.Manager
    globalPredicates []predicate.Predicate
    ctrl             controller.Controller
    ctrlOptions      controller.Options
    name             string
}
```

ctrlOptions指定构建controller的一些配置，主要是Reconciler。forInput指定被协调的对象本身，通过build.For()进行设置。ownsInput指定被协调监听的子对象资源，通过build.Owns()进行设置。watchesInput能够自定义EventHandler处理逻辑，通过build.Watches()进行设置。所以，kubebuilder生成的controller默认只会对协调的对象本身进行调协。

```go
type WatchesInput struct {
    src              source.Source
    eventhandler     handler.EventHandler
    predicates       []predicate.Predicate
    objectProjection objectProjection
}
```

builder.Complete()会调用Builder.Build()进行构造。Build()包含doController()和doWatch()这两个重要方法。

### DoController()

doController通过资源对象的 GVK 来获取 Controller 的名称，最后通过一个 newController 函数来实例化Controller。

```go
controllerName, err := blder.getControllerName(gvk, hasGVK)
blder.ctrl, err = newController(controllerName, blder.mgr, ctrlOptions)
```

newContrller为controller.New的别名，方法为：func New(name string, mgr manager.Manager, options Options) (Controller, error) { c, err := NewUnmanaged(name, mgr, options) if err != nil { return nil, err }

```go
// Add the controller as a Manager components
   return c, mgr.Add(c)
}
```

c, err := NewUnmanaged(name, mgr, options)初始化Controller实例，Controller 实例化完成后，又通过 mgr.Add(c) 函数将控制器添加到 Manager 中去进行管理。controllerManager 的 Add 函数传递的是一个 Runnable 参数，Runnable 是一个接口，用来表示可以启动的一个组件，而恰好 Controller 实际上就实现了这个接口的 Start 函数，所以可以通过 Add 函数来添加 Controller 实例。

### DoWatch()

DoWatch实现比较简单，就是调用controller.watch来注册EventHandler事件。DoWatch方法会调用controller.Watch()方法来注册EventHandler。可以看到对于forInput这类资源，默认的EventHandler为EnqueueRequestForObject，对于ownsInput这类资源，默认的EventHandler为EnqueueRequestForOwner，这两类handler已在上文提到过，均实现了Create()、Update()、Delete()等方法，能够将被调协的资源对象入队。

```go
func (blder *Builder) doWatch() error {
    // Reconcile type
    typeForSrc, err := blder.project(blder.forInput.object, blder.forInput.objectProjection)
    if err != nil {
        return err
    }
    src := &source.Kind{Type: typeForSrc}
    hdler := &handler.EnqueueRequestForObject{}
    allPredicates := append(blder.globalPredicates, blder.forInput.predicates...)
    if err := blder.ctrl.Watch(src, hdler, allPredicates...); err != nil {
        return err
    }

    // Watches the managed types
    for _, own := range blder.ownsInput {
        typeForSrc, err := blder.project(own.object, own.objectProjection)
        if err != nil {
            return err
        }
        src := &source.Kind{Type: typeForSrc}
        hdler := &handler.EnqueueRequestForOwner{
            OwnerType:    blder.forInput.object,
            IsController: true,
        }
        allPredicates := append([]predicate.Predicate(nil), blder.globalPredicates...)
        allPredicates = append(allPredicates, own.predicates...)
        if err := blder.ctrl.Watch(src, hdler, allPredicates...); err != nil {
            return err
        }
    }

    // Do the watch requests
    for _, w := range blder.watchesInput {
        allPredicates := append([]predicate.Predicate(nil), blder.globalPredicates...)
        allPredicates = append(allPredicates, w.predicates...)

        // If the source of this watch is of type *source.Kind, project it.
        if srckind, ok := w.src.(*source.Kind); ok {
            typeForSrc, err := blder.project(srckind.Type, w.objectProjection)
            if err != nil {
                return err
            }
            srckind.Type = typeForSrc
        }

        if err := blder.ctrl.Watch(w.src, w.eventhandler, allPredicates...); err != nil {
            return err
        }
    }
    return nil
}
```

watchesInput这类资源需要自己实现EventHandler，使用类似以下方式实现相应功能。根据之前的结论，controller中调协的资源对象来自于queue，而queue中的数据是通过EventHandler的Create、Update、Delete等处理逻辑进行入队的。因此这时controller的处理顺序为：EventHandler中定义的逻辑->入队->Reconcile。

```go
func (r *AppServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

    klog.Infof("开始Reconcile逻辑")
    ...
    return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        Named("appServiceController").
        Watches(
            &source.Kind{
                Type: &appexamplecomv1.AppService{},
            },
            handler.Funcs{
                CreateFunc: func(createEvent event.CreateEvent, limitingInterface workqueue.RateLimitingInterface) {
                    klog.Infof("createFunc")
                    limitingInterface.Add(reconcile.Request{NamespacedName: types.NamespacedName{
                       Name:      createEvent.Object.GetName(),
                       Namespace: createEvent.Object.GetNamespace(),
                    }})
                },
                UpdateFunc: func(updateEvent event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
                    klog.Infof("updateFunc")
                },
                DeleteFunc: func(deleteEvent event.DeleteEvent, limitingInterface workqueue.RateLimitingInterface) {
                    klog.Infof("deleteFunc")
                },
            }).
        Complete(r)
}
```

上述代码只有在Create时进行了入队处理，因此只有在创建资源时会进入Reconcile的逻辑。

### Manager.Start()

在注册controller到manager后，需要使用mgr.Start(ctrl.SetupSignalHandler())来启动manager。之前说过，注册Controller时调用DoController方法中的mgr.Add()将controller已runnable的形式添加到了Manager。Manager.start()正是调用了cm.runnables的start方法，也即controller.start()来启动controller。

```go
func (cm *controllerManager) Start(ctx context.Context) (err error) {
    ...
    if err := cm.runnables.Webhooks.Start(cm.internalCtx); err != nil {
        if !errors.Is(err, wait.ErrWaitTimeout) {
            return err
        }
    }

    // Start and wait for caches.
    if err := cm.runnables.Caches.Start(cm.internalCtx); err != nil {
        if !errors.Is(err, wait.ErrWaitTimeout) {
            return err
        }
    }

    // Start the non-leaderelection Runnables after the cache has synced.
    if err := cm.runnables.Others.Start(cm.internalCtx); err != nil {
        if !errors.Is(err, wait.ErrWaitTimeout) {
            return err
        }
    }
    ...
}
```
