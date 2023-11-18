[Golang程序性能分析（二）在Echo和Gin框架中使用pprof](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486654&idx=1&sn=ea7171f58254dfecebc61cfbec7b64e5#wechat_redirect)
================================================================================================================================================================

网管叨bi叨

前言
--

今天继续分享使用`Go`官方库`pprof`做性能分析相关的内容，上一篇文章：[Golang程序性能分析（一）pprof和go-torch](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486618&idx=1&sn=bb5e76e011ba99ebc2ffb8f9d3c00b89&chksm=fa80dd0dcdf7541b641be90bdf39001a8ac417c418bb19b9ee7d9e493a0627cfe7d144bb2163&token=228666042&lang=zh_CN&scene=21#wechat_redirect)中我花了很大的篇幅介绍了如何使用`pprof`采集`Go`应用程序的性能指标，如何找到运行缓慢的函数，以及函数中每一部的性能消耗细节。这一节的重点会放在如何在`Echo`和`Gin`这两个框架中增加对`pprof` HTTP请求的支持，因为`pprof`只是提供了对`net/http`包的`ServerMux`的路由支持，这些路由想放到`Echo`和`Gin`里使用时，还是需要有点额外的集成工作。

等集成到框架里，能通过HTTP访问`pprof`提供的几个路由后，`go tool pprof`工具还是通过访问这些URL把性能数拿到本地后来分析的，后续的性能数据采集和分析的操作就跟上篇文章里介绍的完全一样了，并没有因为使用的框架不一样而有什么差别。

在Echo中使用pprof
-------------

由于`Echo`框架使用的复用器`ServerMux是`自定义的，需要手动注册`pprof`提供的路由，网上有几个把他们封装成了包可以直接使用， 不过都不是官方提供的包。后来我看了一下`pprof`提供的路由`Handler`的源码，只需要把它转换成`Echo`框架的路由`Handler`后即可能正常处理那些`pprof`相关的请求，具体转换操作很简单我就直接放代码了。

```go
func RegisterRoutes(engine *echo.Echo) {
  router := engine.Group("")
  ......
  // 下面的路由根据要采集的数据需求注册，不用全都注册
 router.GET("/debug/pprof", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
 router.GET("/debug/pprof/allocs", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
 router.GET("/debug/pprof/block", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
 router.GET("/debug/pprof/goroutine", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
 router.GET("/debug/pprof/heap", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
 router.GET("/debug/pprof/mutex", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
 router.GET("/debug/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
 router.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
 router.GET("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
 router.GET("/debug/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
}
```

注册好路由后还需要对`Echo`框架的写响应超时`WriteTimeout`做一下配置，保证发生写超时的时间设置要大于`pprof`做数据采集的时间，这个配置对应的是`/debug/pprof`路由的`seconds`参数，默认采集时间是30秒，比如我通常要进行60秒的数据采集，那`WriteTimeout`配置的时间就要超过60秒，具体配置方式如下：

> 如果pprof做profiling的时间超过WriteTimeout会引发一个 "profile duration exceeds server's _WriteTimeout_"的错误。

```go
RegisterRoutes(engine)

err := engine.StartServer(&http.Server{
   Addr:              addr,
   ReadTimeout:       time.Second * 5,
   ReadHeaderTimeout: time.Second * 2,
   WriteTimeout:      time.Second * 90,
})
```

上面两步都设置完后就能够按照上面文件里介绍的pprof子命令进行性能分析了

```go
➜ go tool pprof http://{server_ip}:{port}/debug/pprof/profile
Fetching profile over HTTP from http://localhost/debug/pprof/profile
Saved profile in /Users/Kev/pprof/pprof.samples.cpu.005.pb.gz
Type: cpu
Time: Nov 15, 2020 at 3:32pm (CST)
Duration: 30.01s, Total samples = 0
No samples were found with the default sample value type.
Try "sample_index" command to analyze different sample values.
Entering interactive mode (type "help" for commands, "o" for options)
(pprof)
```

具体`pprof`常用子命令的使用方法，可以参考文章[Golang程序性能分析（一）pprof和go-torch](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486618&idx=1&sn=bb5e76e011ba99ebc2ffb8f9d3c00b89&chksm=fa80dd0dcdf7541b641be90bdf39001a8ac417c418bb19b9ee7d9e493a0627cfe7d144bb2163&token=228666042&lang=zh_CN&scene=21#wechat_redirect)里的内容。

在Gin中使用pprof
------------

在`Gin`框架可以通过安装`Gin`项目组提供的`gin-contrib/pprof`包，直接引入后使用就能提供`pprof`相关的路由访问。

```go
import "github.com/gin-contrib/pprof"

package main

import (
 "github.com/gin-contrib/pprof"
 "github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
  pprof.Register(router)
  router.Run(":8080")
}
```

这个包还支持把`pprof`路由划分到单独的路由组里，具体可以查阅gin-contrib/pprof的文档。

上面这些都配置完，启动服务后就能使用`go tool pprof`进行数据采集和分析：

*   内存使用信息采集
    

```go
go tool pprof http://localhost:8080/debug/pprof/heap
```

*   CPU使用情况信息采集
    

```go
go tool pprof http://localhost:8080/debug/pprof/profile
```

总结
--

用`go tool pprof`能对所有类型的`Go`应用程序的性能进行分析，这次的文章主要说的是怎么在`Echo`和`Gin`这两个框架里开启对`pprof`性能采集的支持，具体对程序性能分析的方法和步骤还是和第一篇[Golang程序性能分析（一）pprof和go-torch](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486618&idx=1&sn=bb5e76e011ba99ebc2ffb8f9d3c00b89&chksm=fa80dd0dcdf7541b641be90bdf39001a8ac417c418bb19b9ee7d9e493a0627cfe7d144bb2163&token=228666042&lang=zh_CN&scene=21#wechat_redirect)中重点讲的内容一样。

**近期文章推荐**

[Go语言init函数你必须记住的六个特征](http://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486639&idx=1&sn=cdf2497e73329fa2bc18652864aec7de&chksm=fa80dd38cdf7542eb38c3b4f904244020b590127b2542cc87fb8e904de13eebdca305d0d8ad6&scene=21#wechat_redirect)

[ConfigMap用管理对象的方式管理配置](http://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247486416&idx=1&sn=20d568f93d0f39e0f3c7ef3ce42ac1d8&chksm=fa80da47cdf75351af5919ece8169808807a3dab2ae2305dde19c268b4a60d335d37d50c0b05&scene=21#wechat_redirect)

> 感谢你能读到这里，后面还会有一篇说用`pprof`分析`gRPC`服务性能的文章，还请大家多多支持，持续关注公众号「网管叨bi叨」。最后铁针儿们求关注、求转发、求在看，不要再说下次一定啦~！
