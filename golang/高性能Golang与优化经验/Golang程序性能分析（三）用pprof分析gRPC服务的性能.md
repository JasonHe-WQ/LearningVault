[Golang程序性能分析（三）用pprof分析gRPC服务的性能](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486669&idx=1&sn=2be1d2cb1b85c8d0f59a314cafb4637f#wechat_redirect)
==============================================================================================================================================================

网管叨bi叨

这是Golang程序性能分析系列文章的最后一篇，这次我们的主要内容是**如何使用`pprof`工具对[gRPC服务](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247483894&idx=1&sn=9bad832154f63015096debcb7fcdca63&chksm=fa80d061cdf75977db9171ee6bbf4fac3919743c8700ebd49671a170a8ceb5f314e74f15b2fb&scene=21&cur_album_id=1576438069854027776#wechat_redirect)的程序性能进行分析**。关于`gRPC`这个框架的文章之前已经写过不少文章了，如果你对它还不太熟悉，不知道它是用来干什么的，可以通过[**gRPC入门系列**](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247483889&idx=1&sn=141ab10eb58a3df89dba390fd5df0b1f&chksm=fa80d066cdf759705afbaf5b3a83df9783611665c3cb03d2a6605359e9ccdb70b836143f9fb7&scene=21&cur_album_id=1576438069854027776#wechat_redirect)的文章对它先做个了解。

怎么用pprof分析gRPC的性能
-----------------

`gRPC`底层基于`HTTP`协议的，一个典型的gRPC服务的启动程序可能像下面这样

```go
func main () {
  lis, err := net.Listen("tcp", 10000)
  grpcServer := grpc.NewServer()
  pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})
  grpcServer.Serve(lis)
}
```

它是一个`RPC`框架不是`Web`框架，不支持浏览器用`URL`访问，所以也就没法向上一节给[`Echo`和`Gin`框架单独注册`pprof`采集数据用的那些路由](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486654&idx=1&sn=ea7171f58254dfecebc61cfbec7b64e5&chksm=fa80dd29cdf7543f97f90723e0470fbcb437f9a614a44c51ccf64e98f56e54711d713e204ab0&token=1150263485&lang=zh_CN&scene=21#wechat_redirect)。但是我们可以换个角度来看这个问题，`pprof`做CPU分析原理是按照一定的频率采集程序CPU（包括寄存器）的使用情况，确定应用程序在主动消耗 CPU 周期时花费时间的位置。所以我们可以在`gRPC`服务启动时，异步启动一个监听其他端口的`HTTP`服务，通过这个`HTTP`服务间接获取`gRPC`服务的分析数据。

```go
go func() {
   http.ListenAndServe(":10001", nil)
}()
```

由于使用默认的`ServerMux`（服务复用器），所以只要匿名导入`net/http/pprof`包，这个`HTTP`的复用器默认就会注册`pprof`相关的路由。

此外建议在启动程序的最开端，调用`runtime.SetBlockProfileRate(1)`指示对阻塞超过1纳秒的`goroutine`进行数据采集。

```go
func main () {

  runtime.SetBlockProfileRate(1)
  go func() {
    http.ListenAndServe(":10001", nil)
  }()
  
  lis, err := net.Listen("tcp", 10000)
  grpcServer := grpc.NewServer()
  pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})
  grpcServer.Serve(lis)
}
```

服务启动后就能通过`{server_ip}:10001/debug/pprof/profile`采集CPU的使用情况了，具体`pprof`工具的使用方法的详细说明参考系列的[第一篇文章](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486618&idx=1&sn=bb5e76e011ba99ebc2ffb8f9d3c00b89&chksm=fa80dd0dcdf7541b641be90bdf39001a8ac417c418bb19b9ee7d9e493a0627cfe7d144bb2163&token=1150263485&lang=zh_CN&scene=21#wechat_redirect)。

下图是我用分析数据生成的函数调用图的一部分，可以看到虽然是用另外一个端口的`HTTP`服务拿到的分析数据，但依然能采集到监听另一个端口的`gRPC`服务程序的CPU使用情况。

![图片](https://mmbiz.qpic.cn/mmbiz_png/z4pQ0O5h0f4iaaTkrDicicucniapfr1IYvibPrmY9IQibkOqNdbDqRibpciaRjPNUwsSBzGG2kRNL1Z9oEhgTLSdXantVQ/640?wx_fmt=png)

gRPC服务的函数调用图

pprof的局限
--------

pprof这些功能虽然很有用，但是想分析出程序的性能问题还是挺费事儿的，从我使用下来的感觉主要有两点。

首先，因为调用图里把所有函数调用都显示出来了，有些耗时长的还是Go底层的`runtime`包内函数的直接，想要在这一堆里找到慢的业务函数还是得花不少力气。

再一个现在很多服务都是分布式的，如果服务A调用了服务B，服务B里的方法执行的比较耗时的话，在A的分析数据里只能知道`grpc.invoke`（客户端调用`gRPC`方法的请求都是由`invoke`发出的）耗时长，这时又得去服务B上采集数据，做不到全链路服务性能的采集，这块如果谁知道好的解决方案可以在留言里说一下。

> 这期的文章就到这里，欢迎在留言里多交流。下期会推送一篇关于用Kubernetes StatefulSet控制器编排有状态应用的超长文章，同时会对Headless Service做一个详细的分析，想入门K8s的铁汁儿们，微信还没关注公众号「网管叨bi叨」的，赶紧关注一波呀！！！

相关阅读
----

[Golang程序性能分析（一）pprof和go-torch](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486618&idx=1&sn=bb5e76e011ba99ebc2ffb8f9d3c00b89&chksm=fa80dd0dcdf7541b641be90bdf39001a8ac417c418bb19b9ee7d9e493a0627cfe7d144bb2163&token=1150263485&lang=zh_CN&scene=21#wechat_redirect)

[Golang程序性能分析（二）在Echo和Gin框架中使用pprof](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486654&idx=1&sn=ea7171f58254dfecebc61cfbec7b64e5&chksm=fa80dd29cdf7543f97f90723e0470fbcb437f9a614a44c51ccf64e98f56e54711d713e204ab0&token=1150263485&lang=zh_CN&scene=21#wechat_redirect)

- END -

关注公众号，每周教会你一个进阶知识

![图片](https://mmbiz.qpic.cn/mmbiz_png/z4pQ0O5h0f4icJbGAQ8RjXUUVdUZsGADuMBVWePgn7tfrWjjHyc6b8kXTQ7Sdkp0QQFFK4mel5tniczqooMna1CA/640?wx_fmt=png)

[grpc](javascript:void(0))[Go程序性能分析](javascript:void(0))
